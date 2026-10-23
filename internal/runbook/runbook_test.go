package runbook

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
)

func testPlan() *plan.Plan {
	return &plan.Plan{
		From: "2.346.3", To: "2.580.1", JavaFrom: []int{8, 11, 17}, JavaTo: []int{21, 25},
		Steps: []plan.Step{
			{Core: "2.346.3", Changes: []plan.Change{{Name: "git", From: "4.0", To: "4.11.5"}}},
			{Core: "2.541.3", Java: 17, Upgrade: true, Changes: []plan.Change{{Name: "git", From: "4.11.5", To: "5.2.2"}}},
			{Core: "2.580.1", Java: 21, Upgrade: true, Changes: []plan.Change{{Name: "git", From: "5.2.2", To: "5.8.0"}}},
		},
		Agents: []plan.AgentNote{{Name: "solaris10", Problem: "no Java 21"}},
	}
}

func TestDeb(t *testing.T) {
	var b strings.Builder
	if err := Write(&b, testPlan(), Options{Install: Deb, URL: "https://ci.example.com"}); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{
		"JENKINS_URL=https://ci.example.com",
		"Window 1: update plugins on 2.346.3",
		"Window 2: upgrade Jenkins to 2.580.1",
		"--step 0 --dir \"$WORK/step0-plugins\"",
		"--step 2 --dir \"$WORK/step2-plugins\"",
		"jenkins=2.541.3",
		"jenkins=2.580.1",
		"Install Java 17 and 21",
		`JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64`,
		"api -X POST \"$JENKINS_URL/quietDown\"",
		"jenkins-upgrade-check probe --quiet",
		"--core 2.541.3",
		"--core 2.580.1",
		"## Roll back",
		"jenkins=2.346.3",
		`tar -tzf "$BACKUP" > /dev/null &&`,
		`mkdir -p "$WORK" && chmod 700 "$WORK"`,
		`[ "$(stat -c %u "$WORK")" = 0 ] || echo "STOP:`,
		`(umask 077 && rm -f "$BACKUP" && tar -C "$JENKINS_HOME"`,
		`df -h "$(dirname "$BACKUP")" "$JENKINS_HOME"`,
		"grep -qE 'Failed to initialize Jenkins",
		"solaris10",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("runbook lacks %q", want)
		}
	}
}

func TestDocker(t *testing.T) {
	var b strings.Builder
	if err := Write(&b, testPlan(), Options{Install: Docker}); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{"docker pull jenkins/jenkins:2.541.3-jdk17", "docker pull jenkins/jenkins:2.580.1-jdk21", "chown -R 1000:1000", `docker logs --since "$START" "$CONTAINER" > "$WORK/start.log" 2>&1`} {
		if !strings.Contains(out, want) {
			t.Errorf("runbook lacks %q", want)
		}
	}
	if strings.Contains(out, "Install Java") {
		t.Error("the image carries its Java; the Docker runbook must not ask to install one")
	}
}

// TestShellSyntax runs every command block through bash -n: a quoting
// slip in a long line would otherwise surface in the middle of a window.
func TestShellSyntax(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("no bash")
	}
	for _, install := range []string{Deb, RPM, War, Docker} {
		var b strings.Builder
		if err := Write(&b, testPlan(), Options{Install: install}); err != nil {
			t.Fatal(err)
		}
		var script strings.Builder
		in := false
		for _, line := range strings.Split(b.String(), "\n") {
			switch {
			case line == "```sh":
				in = true
			case line == "```":
				in = false
			case in:
				script.WriteString(line + "\n")
			}
		}
		// Jenkins logs to stderr: "2>&1 > file" would send it to the
		// terminal and leave verify a log without the reason.
		if strings.Contains(script.String(), "2>&1 >") {
			t.Errorf("%s runbook redirects stderr before stdout", install)
		}
		cmd := exec.Command(bash, "-n")
		cmd.Stdin = strings.NewReader(script.String())
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("%s runbook: %v\n%s", install, err, out)
		}
	}
}

func TestUnknownInstall(t *testing.T) {
	if err := Write(&strings.Builder{}, testPlan(), Options{Install: "msi"}); err == nil {
		t.Error("an unknown install type must be refused")
	}
}

// A war controller on Java 11 that has three stops to make: the runbook
// lists every stop, switches Java only where the controller does not run it
// already, and skips the plugin window when the current core needs no
// plugin change.
func TestWarThreeStops(t *testing.T) {
	var agents []plan.AgentNote
	for i := 1; i <= 12; i++ {
		agents = append(agents, plan.AgentNote{Name: fmt.Sprintf("hpux%d", i), Problem: "no Java 21"})
	}
	p := &plan.Plan{
		From: "2.249.3", To: "2.580.1", ControllerJava: "11.0.20",
		Steps: []plan.Step{
			{Core: "2.249.3"},
			{Core: "2.462.3", Java: 11, Upgrade: true},
			{Core: "2.541.3", Java: 17, Upgrade: true, Changes: []plan.Change{{Name: "git", From: "4.0", To: "5.2.2"}},
				Problems: []string{"gitlab-plugin 1.5 needs jersey3-api 2.0 or newer, and no such release runs on Jenkins 2.541.3"}},
			{Core: "2.580.1", Java: 21, Upgrade: true, Changes: []plan.Change{{Name: "git", From: "5.2.2", To: "5.8.0"}}},
		},
		Agents: agents,
	}
	var b strings.Builder
	if err := Write(&b, p, Options{Install: War, Home: "/srv/jenkins", PlanFile: "out/plan.json"}); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{
		"For: jenkins.war run by a service. Made from `out/plan.json`.",
		"so the upgrade stops at 2.462.3, 2.541.3 on the way",
		"JENKINS_HOME=/srv/jenkins\n",
		"WAR_DIR=/opt/jenkins",
		"## Window 1: upgrade Jenkins to 2.580.1",
		"### Stop 2.462.3",
		"### Stop 2.541.3",
		"--step 2 --dir \"$WORK/step2-plugins\"   # for 2.541.3",
		"--step 3 --dir \"$WORK/step3-plugins\"   # for 2.580.1",
		"Install Java 17 and 21 next to the current one",
		"Eclipse Temurin from adoptium.net, for example.",
		`printf '[Service]\nEnvironment="JAVA_HOME=/usr/lib/jvm/java-17"\n' | tee "/etc/systemd/system/$SERVICE.service.d/java.conf"`,
		"- ⚠ gitlab-plugin 1.5 needs jersey3-api",
		"GO? Leave it in quiet mode and go on to the next stop.",
		"hpux1, hpux2, hpux3, hpux4, hpux5, hpux6, hpux7, hpux8, hpux9, hpux10 and 2 more",
		`ln -sfn "jenkins-2.249.3.war" "$WAR_DIR/jenkins.war"`,
		`rm -f "/etc/systemd/system/$SERVICE.service.d/java.conf"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("runbook lacks %q", want)
		}
	}
	for _, unwanted := range []string{"update plugins on 2.249.3", "--step 0", "--step 1", "JAVA_HOME=/usr/lib/jvm/java-11"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("runbook has %q", unwanted)
		}
	}
}

// With no plugin to change, the download step says so instead of listing
// nothing.
func TestNoPluginChanges(t *testing.T) {
	p := &plan.Plan{From: "2.541.3", To: "2.555.1", ControllerJava: "21",
		Steps: []plan.Step{{Core: "2.541.3"}, {Core: "2.555.1", Java: 21, Upgrade: true}}}
	var b strings.Builder
	if err := Write(&b, p, Options{Install: RPM}); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{"# no plugin changes", "dnf install -y jenkins-2.555.1", "dnf downgrade -y jenkins-2.541.3", "JENKINS_HOME=/var/lib/jenkins\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("runbook lacks %q", want)
		}
	}
	if strings.Contains(out, "Install Java") || strings.Contains(out, "java.conf") {
		t.Error("the controller already runs Java 21: no Java switch")
	}
}

// When a stop has no Java (a core older than the Java support policy), the
// Docker runbook takes the newest Java the core runs on, or the plain tag
// when no Java is known.
func TestDockerImageWithoutJava(t *testing.T) {
	p := &plan.Plan{From: "1.651.3", To: "2.580.1",
		Steps: []plan.Step{{Core: "1.651.3"}, {Core: "1.651.3", Upgrade: true}, {Core: "2.580.1", Upgrade: true}}}
	var b strings.Builder
	if err := Write(&b, p, Options{Install: Docker, Home: "/data/jenkins"}); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{"docker pull jenkins/jenkins:1.651.3\n", "docker pull jenkins/jenkins:2.580.1-jdk25\n", "JENKINS_HOME=/data/jenkins\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("runbook lacks %q", want)
		}
	}
}
