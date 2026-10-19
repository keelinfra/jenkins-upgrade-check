// Package runbook turns a plan into the commands an administrator runs to
// upgrade a controller without users noticing more than a short, announced
// window, and to roll back if the upgrade does not hold.
//
// The method is the same for every install type: rehearse on a copy first;
// in the window, let running builds finish, stop, take a backup, upgrade,
// start in quiet mode so no build runs on the new version, verify against a
// snapshot taken before, and only then let builds run. Rolling back before
// that point restores the backup and loses nothing.
package runbook

import (
	"fmt"
	"io"
	"strings"

	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
	"github.com/keelinfra/jenkins-upgrade-check/internal/version"
)

// Install types.
const (
	Deb    = "deb"
	RPM    = "rpm"
	War    = "war"
	Docker = "docker"
)

// Installs lists the install types a runbook can be written for.
var Installs = []string{Deb, RPM, War, Docker}

// Options describe the controller.
type Options struct {
	Install  string
	Home     string // JENKINS_HOME on the host
	URL      string
	PlanFile string // the plan JSON, as the commands refer to it
}

func (o Options) home() string {
	if o.Home != "" {
		return o.Home
	}
	switch o.Install {
	case Docker:
		return "/var/lib/docker/volumes/jenkins_home/_data"
	case War:
		return "/var/lib/jenkins"
	}
	return "/var/lib/jenkins"
}

func (o Options) owner() string {
	if o.Install == Docker {
		return "1000:1000"
	}
	return "jenkins:jenkins"
}

// Write writes the runbook for the plan as Markdown.
func Write(w io.Writer, p *plan.Plan, o Options) error {
	ok := false
	for _, i := range Installs {
		ok = ok || i == o.Install
	}
	if !ok {
		return fmt.Errorf("install type %q: use one of %s", o.Install, strings.Join(Installs, ", "))
	}
	if o.PlanFile == "" {
		o.PlanFile = "plan.json"
	}
	b := &strings.Builder{}
	r := &writer{b: b, p: p, o: o}
	r.intro()
	r.settings()
	r.before()
	n := 1
	if len(p.Steps) > 0 && len(p.Steps[0].Changes) > 0 {
		r.pluginWindow(p.Steps[0], n)
		n++
	}
	r.upgradeWindow(n)
	r.rollback()
	r.after()
	_, err := io.WriteString(w, b.String())
	return err
}

type writer struct {
	b *strings.Builder
	p *plan.Plan
	o Options
}

func (r *writer) f(format string, args ...any) { fmt.Fprintf(r.b, format, args...) }

func (r *writer) sh(lines ...string) {
	r.b.WriteString("```sh\n")
	for _, l := range lines {
		r.b.WriteString(l + "\n")
	}
	r.b.WriteString("```\n\n")
}

var installName = map[string]string{
	Deb:    "Debian/Ubuntu package (apt)",
	RPM:    "RHEL/Rocky/Alma/Fedora package (dnf)",
	War:    "jenkins.war run by a service",
	Docker: "Docker, official jenkins/jenkins image with JENKINS_HOME in a volume",
}

func (r *writer) hops() []plan.Step {
	var out []plan.Step
	for _, s := range r.p.Steps {
		if s.Upgrade {
			out = append(out, s)
		}
	}
	return out
}

func (r *writer) intro() {
	p := r.p
	r.f("# Upgrade runbook: Jenkins %s → %s\n\n", p.From, p.To)
	r.f("For: %s. Made from `%s`.\n\n", installName[r.o.Install], r.o.PlanFile)
	r.b.WriteString("How this keeps users safe and the way back open:\n\n")
	r.b.WriteString("- **Rehearse first.** The upgrade runs on a copy of your configuration in a container with no network, until it comes up clean. Nothing reaches users or other systems.\n")
	r.b.WriteString("- **One announced window, no killed builds.** Jenkins stops accepting new builds and waits for running ones to finish before it stops.\n")
	r.b.WriteString("- **Nothing runs on the new version until it is verified.** It starts in quiet mode; `verify` compares it with a snapshot from before and says GO or ROLL BACK, and why.\n")
	r.b.WriteString("- **Rollback is a restore, not a repair.** A backup is taken right before the upgrade. Because no build ran on the new version, restoring it loses nothing. The commands are below, with the exact versions.\n\n")
	hops := r.hops()
	if len(hops) > 1 {
		r.f("No Java version runs both %s and %s, so the upgrade stops at", p.From, p.To)
		for i, h := range hops[:len(hops)-1] {
			if i > 0 {
				r.b.WriteString(",")
			}
			r.f(" %s", h.Core)
		}
		r.b.WriteString(" on the way. Jenkins has to start once on each stop, so its data migrations run.\n\n")
	}
}

func (r *writer) settings() {
	o := r.o
	r.b.WriteString("## Settings\n\nSet these once in the shell you run the commands in. The API token is for a user with Overall/Administer.\n\n")
	url := o.URL
	if url == "" {
		url = "https://jenkins.example.com"
	}
	lines := []string{
		"JENKINS_URL=" + url,
		"JENKINS_USER=admin",
		"read -rs JENKINS_TOKEN && export JENKINS_TOKEN   # paste the API token",
		"JENKINS_HOME=" + o.home(),
		"WORK=/var/tmp/jenkins-upgrade              # downloads, snapshots, reports",
		"BACKUP=/var/backups/jenkins-before-upgrade.tgz",
	}
	switch o.Install {
	case Docker:
		lines = append(lines, "CONTAINER=jenkins                           # docker ps shows the name",
			"# JENKINS_HOME is the host side of the volume: docker volume inspect jenkins_home -f '{{.Mountpoint}}'")
	case War:
		lines = append(lines, "WAR_DIR=/opt/jenkins                        # where jenkins.war is", "SERVICE=jenkins                             # the service that runs it")
	}
	lines = append(lines, `mkdir -p "$WORK"`, `api() { curl -fsS -u "$JENKINS_USER:$JENKINS_TOKEN" "$@"; }`)
	r.sh(lines...)
}

func (r *writer) before() {
	p := r.p
	pf := r.o.PlanFile
	r.b.WriteString("## Before the window\n\n")
	r.b.WriteString("1. **Rehearse until it is clean.** On any machine with Docker and a copy of JENKINS_HOME (only the configuration is read, not build records):\n\n")
	r.sh(`sudo tar -C "$JENKINS_HOME" --exclude=builds --exclude=workspace --exclude=./plugins --exclude=./war --exclude=./caches --exclude=./logs -czf "$WORK/jenkins-config.tgz" .   # holds secrets: keep it private`,
		`mkdir -p "$WORK/config" && tar -C "$WORK/config" -xzf "$WORK/jenkins-config.tgz"`,
		fmt.Sprintf(`jenkins-upgrade-check dry-run --plan %s --home "$WORK/config" -o "$WORK/rehearsal.md"`, pf))
	r.b.WriteString("   Fix every row under \"What failed, and what to do\", make a new plan, and rehearse again, until the report says every plugin loaded and all jobs came back. Delete the configuration copy when done.\n\n")
	r.b.WriteString("2. **Download the plugins now,** so the window does not depend on the network:\n\n")
	var dl []string
	for i, s := range p.Steps {
		if len(s.Changes) == 0 {
			continue
		}
		dl = append(dl, fmt.Sprintf(`jenkins-upgrade-check download --plan %s --step %d --dir "$WORK/step%d-plugins"   # for %s`, pf, i, i, s.Core))
	}
	if len(dl) == 0 {
		dl = append(dl, "# no plugin changes")
	}
	r.sh(dl...)
	if pkg := r.packagePrefetch(); len(pkg) > 0 {
		r.b.WriteString("3. **Fetch the new Jenkins** without installing it:\n\n")
		r.sh(pkg...)
	}
	if javas := r.javaChanges(); len(javas) > 0 {
		r.f("4. **Install Java %s next to the current one** (do not switch yet; skip a version the controller already runs). %s\n\n", joinInts(javas, " and "), r.javaInstallHint(javas))
	}
	if len(p.Agents) > 0 {
		r.f("5. **Agents that cannot follow** (from the plan): %s. Sort them out first, or they stay offline after the upgrade.\n\n", agentNames(p))
	}
	r.b.WriteString("6. **Time the backup**, so you know how long the window is: the backup below copies everything but workspaces.\n\n")
	r.sh(`sudo du -sh --exclude=workspace "$JENKINS_HOME"`)
	r.b.WriteString("   Hundreds of GB of build history? Take a disk or filesystem snapshot (LVM, ZFS, a cloud volume snapshot) instead of the tar backup, and restore from it in the rollback.\n\n")
	r.b.WriteString("7. **Announce the window** to the people who use Jenkins.\n\n")
}

func (r *writer) packagePrefetch() []string {
	hops := r.hops()
	var out []string
	for _, h := range hops {
		switch r.o.Install {
		case Deb:
			out = append(out, fmt.Sprintf("sudo apt-get install --download-only -y --allow-change-held-packages jenkins=%s", h.Core))
		case RPM:
			out = append(out, fmt.Sprintf("sudo dnf install --downloadonly -y jenkins-%s", h.Core))
		case War:
			out = append(out, fmt.Sprintf(`curl -fsSLo "$WORK/jenkins-%[1]s.war" https://get.jenkins.io/war-stable/%[1]s/jenkins.war`, h.Core),
				fmt.Sprintf(`echo "$(curl -fsSL https://get.jenkins.io/war-stable/%[1]s/jenkins.war.sha256 | cut -d' ' -f1)  $WORK/jenkins-%[1]s.war" | sha256sum -c -`, h.Core))
		case Docker:
			out = append(out, "docker pull "+r.image(h.Core, h.Java))
		}
	}
	if r.o.Install == Deb && len(out) > 0 {
		out = append([]string{"sudo apt-get update   # if apt cannot find a version: the repository must be pkg.jenkins.io/debian-stable, with its current signing key"}, out...)
	}
	return out
}

// image is the official image for a core on the Java the plan picked.
func (r *writer) image(core string, java int) string {
	if java == 0 {
		if js := version.SupportedJava(core); len(js) > 0 {
			java = js[len(js)-1]
		}
	}
	if java == 0 {
		return "jenkins/jenkins:" + core
	}
	return fmt.Sprintf("jenkins/jenkins:%s-jdk%d", core, java)
}

// javaChanges lists the Java versions the controller has to switch to on
// the way. When the plan does not know the Java it runs, every hop's Java
// is listed; the commands say to skip a switch to the Java already running.
func (r *writer) javaChanges() []int {
	if r.o.Install == Docker {
		return nil // the image carries its Java
	}
	cur := plan.JavaMajor(r.p.ControllerJava)
	var out []int
	for _, h := range r.hops() {
		if h.Java != 0 && h.Java != cur {
			out = append(out, h.Java)
			cur = h.Java
		}
	}
	return out
}

func (r *writer) javaInstallHint(javas []int) string {
	j := javas[len(javas)-1]
	switch r.o.Install {
	case Deb:
		return fmt.Sprintf("For example `sudo apt-get install -y openjdk-%d-jre-headless`, or Eclipse Temurin from packages.adoptium.net if your release does not have it.", j)
	case RPM:
		return fmt.Sprintf("For example `sudo dnf install -y java-%d-openjdk-headless`, or Eclipse Temurin from packages.adoptium.net.", j)
	}
	return "Eclipse Temurin from adoptium.net, for example."
}

func (r *writer) javaHome(java int) string {
	switch r.o.Install {
	case Deb:
		return fmt.Sprintf("/usr/lib/jvm/java-%d-openjdk-amd64", java)
	case RPM:
		return fmt.Sprintf("/usr/lib/jvm/jre-%d", java)
	}
	return fmt.Sprintf("/usr/lib/jvm/java-%d", java)
}

// drain is the start of every window: no new builds, wait for running ones,
// a snapshot to compare with, stop, back up.
func (r *writer) drain() {
	r.b.WriteString("**Let running builds finish, then stop.** New builds wait in the queue; nothing running is killed.\n\n")
	r.sh(`jenkins-upgrade-check collect --url "$JENKINS_URL" --user "$JENKINS_USER" -o "$WORK/before.json"`,
		`api -X POST "$JENKINS_URL/quietDown"`,
		`until api "$JENKINS_URL/computer/api/json?tree=busyExecutors" | grep -q '"busyExecutors":0'; do sleep 30; done`,
		r.stop())
	r.b.WriteString("**Back up.** Everything but workspaces, which builds recreate:\n\n")
	r.sh(`sudo tar -C "$JENKINS_HOME" --exclude=./workspace --exclude=./caches -czf "$BACKUP" .`)
}

func (r *writer) stop() string {
	switch r.o.Install {
	case Docker:
		return `docker stop "$CONTAINER"`
	case War:
		return `sudo systemctl stop "$SERVICE"`
	}
	return "sudo systemctl stop jenkins"
}

func (r *writer) start() []string {
	switch r.o.Install {
	case Docker:
		return []string{`START=$(date -u +%Y-%m-%dT%H:%M:%SZ)`, "docker compose up -d   # or your docker run command, with the image above"}
	case War:
		return []string{`START=$(date '+%Y-%m-%d %H:%M:%S')`, `sudo systemctl start "$SERVICE"`}
	}
	return []string{`START=$(date '+%Y-%m-%d %H:%M:%S')`, "sudo systemctl start jenkins"}
}

func (r *writer) logCmd() string {
	switch r.o.Install {
	case Docker:
		return `docker logs --since "$START" "$CONTAINER" > "$WORK/start.log" 2>&1`
	case War:
		return `sudo journalctl -u "$SERVICE" --since "$START" --no-pager > "$WORK/start.log"`
	}
	return `sudo journalctl -u jenkins --since "$START" --no-pager > "$WORK/start.log"`
}

// quietStart installs the probe, starts Jenkins and verifies it.
func (r *writer) quietStart(core string) {
	r.b.WriteString("**Start in quiet mode.** The script below puts Jenkins into quiet mode as it starts and records what loaded and what did not; `jenkins-upgrade-check probe --quiet` prints it, so you can read it first.\n\n")
	lines := []string{
		`sudo mkdir -p "$JENKINS_HOME/init.groovy.d"`,
		`jenkins-upgrade-check probe --quiet | sudo tee "$JENKINS_HOME/init.groovy.d/zz-jenkins-upgrade-check.groovy" > /dev/null`,
		fmt.Sprintf(`sudo chown -R %s "$JENKINS_HOME/init.groovy.d"`, r.o.owner()),
	}
	lines = append(lines, r.start()...)
	lines = append(lines, `until curl -fsS -o /dev/null "$JENKINS_URL/login"; do sleep 5; done`, "sleep 60   # let agents reconnect")
	r.sh(lines...)
	r.b.WriteString("**Verify.** GO exits 0, ROLL BACK exits 1; the report says what failed and why.\n\n")
	r.sh(r.logCmd(),
		fmt.Sprintf(`jenkins-upgrade-check verify --before "$WORK/before.json" --url "$JENKINS_URL" --user "$JENKINS_USER" --plan %s --core %s --home "$JENKINS_HOME" --log "$WORK/start.log" -o "$WORK/verify-%s.md"`, r.o.PlanFile, core, core))
	r.b.WriteString("Then look at two or three jobs that matter most (their configuration page opens, their last build page renders). Agents still offline? Run verify again after a few minutes.\n\n")
}

func (r *writer) open() {
	r.b.WriteString("**GO: let builds run.**\n\n")
	r.sh(`sudo rm "$JENKINS_HOME/init.groovy.d/zz-jenkins-upgrade-check.groovy"`,
		`api -X POST "$JENKINS_URL/cancelQuietDown"`)
	r.b.WriteString("From here, builds run on the new version and the backup no longer holds their records. Keep the backup for a few days anyway.\n\n")
}

func (r *writer) copyPlugins(step int) []string {
	return []string{
		fmt.Sprintf(`sudo cp "$WORK/step%d-plugins/"*.jpi "$JENKINS_HOME/plugins/"`, step),
		fmt.Sprintf(`for f in "$WORK/step%d-plugins/"*.jpi; do sudo rm -f "$JENKINS_HOME/plugins/$(basename "$f" .jpi).hpi"; done   # an old .hpi would shadow the new .jpi`, step),
		fmt.Sprintf(`sudo chown -R %s "$JENKINS_HOME/plugins"`, r.o.owner()),
	}
}

func (r *writer) pluginWindow(s plan.Step, n int) {
	r.f("## Window %d: update plugins on %s\n\n", n, s.Core)
	r.f("The Jenkins upgrade guides ask for plugins to be current before the core moves; this window changes %d plugin(s) and nothing else. Do it a few days before the core upgrade, so a problem is easy to place.\n\n", len(s.Changes))
	r.drain()
	r.b.WriteString("**Update plugins.**\n\n")
	r.sh(r.copyPlugins(0)...)
	r.quietStart(s.Core)
	r.open()
}

func (r *writer) upgradeWindow(n int) {
	hops := r.hops()
	r.f("## Window %d: upgrade Jenkins to %s\n\n", n, r.p.To)
	r.drain()
	for i, s := range r.p.Steps {
		if !s.Upgrade {
			continue
		}
		if len(hops) > 1 {
			r.f("### Stop %s\n\n", s.Core)
		}
		r.f("**Upgrade to %s.**\n\n", s.Core)
		var lines []string
		if i > 0 && len(hops) > 1 && s.Core != hops[0].Core {
			lines = append(lines, r.stop())
		}
		lines = append(lines, r.upgradeCmds(s)...)
		if len(s.Changes) > 0 {
			lines = append(lines, r.copyPlugins(i)...)
		}
		r.sh(lines...)
		for _, pr := range s.Problems {
			r.f("- ⚠ %s\n", pr)
		}
		if len(s.Problems) > 0 {
			r.b.WriteString("\n")
		}
		r.quietStart(s.Core)
		if s.Core != r.p.To {
			r.b.WriteString("GO? Leave it in quiet mode and go on to the next stop. ROLL BACK? See below.\n\n")
		}
	}
	r.open()
}

func (r *writer) upgradeCmds(s plan.Step) []string {
	var out []string
	switch r.o.Install {
	case Deb:
		if s.Java > 0 && r.needsJavaSwitch(s) {
			out = append(out, r.systemdJava(s.Java, "jenkins")...)
		}
		out = append(out, fmt.Sprintf("sudo apt-get install -y --allow-downgrades --allow-change-held-packages jenkins=%s && sudo apt-mark hold jenkins", s.Core))
	case RPM:
		if s.Java > 0 && r.needsJavaSwitch(s) {
			out = append(out, r.systemdJava(s.Java, "jenkins")...)
		}
		out = append(out, fmt.Sprintf("sudo dnf install -y jenkins-%s", s.Core))
	case War:
		if s.Java > 0 && r.needsJavaSwitch(s) {
			out = append(out, r.systemdJava(s.Java, "$SERVICE")...)
		}
		out = append(out, fmt.Sprintf(`sudo cp "$WORK/jenkins-%s.war" "$WAR_DIR/" && sudo ln -sfn "jenkins-%s.war" "$WAR_DIR/jenkins.war"   # the service runs $WAR_DIR/jenkins.war`, s.Core, s.Core))
	case Docker:
		out = append(out, fmt.Sprintf("# in your compose file or docker run command, change the image to %s", r.image(s.Core, s.Java)),
			"# (built your own image FROM jenkins/jenkins with plugins.txt? rebuild it from this image and the step's plugins: plan --plugins-out)")
	}
	return out
}

func (r *writer) needsJavaSwitch(s plan.Step) bool {
	for _, j := range r.javaChanges() {
		if j == s.Java {
			return true
		}
	}
	return false
}

func (r *writer) systemdJava(java int, unit string) []string {
	return []string{
		fmt.Sprintf(`sudo mkdir -p "/etc/systemd/system/%s.service.d"`, unit),
		fmt.Sprintf(`printf '[Service]\nEnvironment="JAVA_HOME=%s"\n' | sudo tee "/etc/systemd/system/%s.service.d/java.conf"   # check the path: ls /usr/lib/jvm`, r.javaHome(java), unit),
		"sudo systemctl daemon-reload",
	}
}

func (r *writer) rollback() {
	p := r.p
	r.b.WriteString("## Roll back\n\n")
	r.b.WriteString("When verify says ROLL BACK, or anything else looks wrong before you let builds run. Keep `$WORK/verify-*.md` and `$WORK/start.log`: they say what to fix before the next attempt.\n\n")
	lines := []string{r.stop()}
	switch r.o.Install {
	case Deb:
		lines = append(lines, fmt.Sprintf("sudo apt-get install -y --allow-downgrades --allow-change-held-packages jenkins=%s && sudo apt-mark hold jenkins", p.From))
	case RPM:
		lines = append(lines, fmt.Sprintf("sudo dnf downgrade -y jenkins-%s", p.From))
	case War:
		lines = append(lines, fmt.Sprintf(`sudo ln -sfn "jenkins-%s.war" "$WAR_DIR/jenkins.war"   # or copy back the war you ran before`, p.From))
	case Docker:
		lines = append(lines, "# put the image back to the tag you ran before, for example jenkins/jenkins:"+p.From)
	}
	if len(r.javaChanges()) > 0 {
		unit := "jenkins"
		if r.o.Install == War {
			unit = "$SERVICE"
		}
		lines = append(lines, fmt.Sprintf(`sudo rm -f "/etc/systemd/system/%s.service.d/java.conf" && sudo systemctl daemon-reload   # back to the Java it ran on`, unit))
	}
	lines = append(lines,
		`sudo mkdir -p "$WORK/restore" && sudo tar -C "$WORK/restore" -xzf "$BACKUP"`,
		`sudo rsync -a --delete --exclude=/workspace --exclude=/caches "$WORK/restore/" "$JENKINS_HOME/"   # also removes the quiet-start script`,
		`sudo rm -rf "$WORK/restore"`)
	lines = append(lines, r.start()...)
	lines = append(lines, `until curl -fsS -o /dev/null "$JENKINS_URL/login"; do sleep 5; done`,
		fmt.Sprintf(`jenkins-upgrade-check verify --before "$WORK/before.json" --url "$JENKINS_URL" --user "$JENKINS_USER" --core %s -o "$WORK/verify-rollback.md"`, p.From))
	r.sh(lines...)
	r.b.WriteString("The rollback is done when this last verify says GO: the controller is what it was before the window. Then tell users it is open again.\n\n")
}

func (r *writer) after() {
	r.b.WriteString("## After\n\n")
	r.b.WriteString("- Watch the first builds of the jobs that matter most; a plugin that loads can still fail inside a build step.\n")
	r.b.WriteString("- Keep the backup until you are sure, then delete it, `$WORK` and the configuration copy (they hold secrets).\n")
	r.b.WriteString("- Run `jenkins-upgrade-check collect` and `plan` again in a few months: staying close to the current LTS keeps the next upgrade one step.\n")
}

func agentNames(p *plan.Plan) string {
	var out []string
	for _, a := range p.Agents {
		out = append(out, a.Name)
	}
	if len(out) > 10 {
		return strings.Join(out[:10], ", ") + fmt.Sprintf(" and %d more", len(out)-10)
	}
	return strings.Join(out, ", ")
}

func joinInts(xs []int, sep string) string {
	s := make([]string, len(xs))
	for i, x := range xs {
		s[i] = fmt.Sprint(x)
	}
	return strings.Join(s, sep)
}
