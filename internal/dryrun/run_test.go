package dryrun

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
)

// fakeDockerScript stands in for the docker CLI. It records each call in
// "calls" and answers from files next to it, so a test sets up a scenario
// by writing files:
//
//	images          images "docker pull" finds, one per line
//	create.fail     "docker create" fails with this output
//	cp-in.fail      copying the configuration in fails with this output
//	start.fail      "docker start" fails with this output
//	logs.N          what the Nth "docker logs" prints (else jenkins.log)
//	state           what "docker inspect" prints (default "true 0")
//	probe.N.json    the probe record the Nth read finds (else probe.json)
//
// A read of the probe record also leaves a file named probe-read.
//
// The configuration copied into the container is kept as ref.tar.
const fakeDockerScript = `#!/bin/sh
d=$(dirname "$0")
echo "$*" >> "$d/calls"
n=$(($(cat "$d/$1.count" 2>/dev/null || echo 0) + 1))
echo "$n" > "$d/$1.count"
case "$1" in
pull)
	if grep -qxF "$3" "$d/images" 2>/dev/null; then echo "$3"; exit 0; fi
	echo "Error response from daemon: manifest for $3 not found: manifest unknown" >&2
	exit 1 ;;
create)
	if [ -f "$d/create.fail" ]; then cat "$d/create.fail" >&2; exit 1; fi
	echo 0123abcd ;;
cp)
	if [ "$2" = "-a" ]; then
		cat > "$d/ref.tar"
		if [ -f "$d/cp-in.fail" ]; then cat "$d/cp-in.fail" >&2; exit 1; fi
		exit 0
	fi
	: > "$d/probe-read"
	if [ -f "$d/probe.$n.json" ]; then cp "$d/probe.$n.json" "$3"; exit 0; fi
	if [ -f "$d/probe.json" ]; then cp "$d/probe.json" "$3"; exit 0; fi
	echo "Error response from daemon: Could not find the file /var/jenkins_home/jenkins-upgrade-check-probe.json in container" >&2
	exit 1 ;;
start)
	if [ -f "$d/start.fail" ]; then cat "$d/start.fail" >&2; exit 1; fi ;;
logs)
	if [ -f "$d/logs.$n" ]; then cat "$d/logs.$n"; else cat "$d/jenkins.log" 2>/dev/null; fi ;;
inspect)
	if [ -f "$d/state" ]; then cat "$d/state"; else echo "true 0"; fi ;;
esac
exit 0
`

type fakeDocker struct {
	t   *testing.T
	dir string
}

func newFakeDocker(t *testing.T) *fakeDocker {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "docker"), []byte(fakeDockerScript), 0o755); err != nil {
		t.Fatal(err)
	}
	return &fakeDocker{t, dir}
}

func (f *fakeDocker) path() string { return filepath.Join(f.dir, "docker") }

func (f *fakeDocker) set(name, body string) {
	f.t.Helper()
	if err := os.WriteFile(filepath.Join(f.dir, name), []byte(body), 0o644); err != nil {
		f.t.Fatal(err)
	}
}

func (f *fakeDocker) calls() []string {
	b, _ := os.ReadFile(filepath.Join(f.dir, "calls"))
	return strings.Split(strings.TrimSpace(string(b)), "\n")
}

func (f *fakeDocker) called(prefix string) bool {
	for _, c := range f.calls() {
		if strings.HasPrefix(c, prefix) {
			return true
		}
	}
	return false
}

// refTar reads the configuration Run copied into the container, by path.
func (f *fakeDocker) refTar() map[string]*tar.Header {
	f.t.Helper()
	b, err := os.ReadFile(filepath.Join(f.dir, "ref.tar"))
	if err != nil {
		f.t.Fatal(err)
	}
	out := map[string]*tar.Header{}
	tr := tar.NewReader(bytes.NewReader(b))
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return out
		}
		if err != nil {
			f.t.Fatal(err)
		}
		out[h.Name] = h
	}
}

func (f *fakeDocker) refFile(name string) string {
	f.t.Helper()
	b, err := os.ReadFile(filepath.Join(f.dir, "ref.tar"))
	if err != nil {
		f.t.Fatal(err)
	}
	tr := tar.NewReader(bytes.NewReader(b))
	for {
		h, err := tr.Next()
		if err != nil {
			f.t.Fatalf("%s not in the tar: %v", name, err)
		}
		if h.Name == name {
			c, _ := io.ReadAll(tr)
			return string(c)
		}
	}
}

// fastWaits makes Run and the downloads poll and retry at once.
func fastWaits(t *testing.T) {
	old := [3]time.Duration{startPoll, probePoll, retryWait}
	startPoll, probePoll, retryWait = time.Millisecond, time.Millisecond, time.Millisecond
	t.Cleanup(func() { startPoll, probePoll, retryWait = old[0], old[1], old[2] })
}

// pluginFile builds a plugin archive: a manifest, classes, and jars in
// WEB-INF/lib holding more classes.
func pluginFile(t *testing.T, manifest string, classes []string, libs map[string][]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	z := zip.NewWriter(&buf)
	if manifest != "" {
		w, _ := z.Create("META-INF/MANIFEST.MF")
		io.WriteString(w, manifest)
	}
	for _, c := range classes {
		if _, err := z.Create("WEB-INF/classes/" + c); err != nil {
			t.Fatal(err)
		}
	}
	for name, cls := range libs {
		var jar bytes.Buffer
		jz := zip.NewWriter(&jar)
		for _, c := range cls {
			jz.Create(c)
		}
		jz.Close()
		w, _ := z.Create("WEB-INF/lib/" + name)
		w.Write(jar.Bytes())
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func sum(b []byte) string {
	s := sha256.Sum256(b)
	return base64.StdEncoding.EncodeToString(s[:])
}

// pluginServer serves files by path and counts the requests for each.
type pluginServer struct {
	*httptest.Server
	mu   sync.Mutex
	hits map[string]int
}

func newPluginServer(t *testing.T, files map[string][]byte) *pluginServer {
	ps := &pluginServer{hits: map[string]int{}}
	ps.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ps.mu.Lock()
		ps.hits[r.URL.Path]++
		ps.mu.Unlock()
		b, ok := files[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Write(b)
	}))
	t.Cleanup(ps.Close)
	return ps
}

func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, body := range files {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

const upLine = "2026-10-21 13:36:10.000+0000 [id=25]\tINFO\thudson.lifecycle.Lifecycle#onReady: Jenkins is fully up and running\n"

// A rehearsal of an upgrade with a real configuration: Jenkins starts
// after a while, the probe reports a moment later, and the report puts
// every problem on the plugin behind it.
func TestRunWithHome(t *testing.T) {
	fastWaits(t)
	d := newFakeDocker(t)
	d.set("images", "jenkins/jenkins:2.580.1-jdk21\n")
	d.set("logs.1", "Running from: /usr/share/jenkins/jenkins.war\n")
	d.set("jenkins.log", "Running from: /usr/share/jenkins/jenkins.war\n"+
		"2026-10-21 13:36:01.000+0000 [id=30]\tINFO\thudson.PluginManager#loadDetachedPlugins: Upgraded Jenkins from version 2.462.3 to version 2.580.1. Loaded detached plugins (and dependencies): [sshd.hpi, instance-identity.jpi]\n"+
		"2026-10-21 13:36:05.000+0000 [id=40]\tSEVERE\thudson.ExtensionFinder$GuiceFinder$SezpozModule#configure: Failed to load org.example.Foo\n"+
		"java.lang.NoClassDefFoundError: javax/servlet/ServletException\n"+
		"2026-10-21 13:36:06.000+0000 [id=46]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed Loading item gone-job\n"+
		"com.thoughtworks.xstream.mapper.CannotResolveClassException: flow-definition\n"+
		upLine)
	d.set("probe.1.json", `{"completed":false}`)
	d.set("probe.json", `{"completed":true,"version":"2.580.1",
		"plugins":[{"name":"good","version":"1.0","active":true,"enabled":true},{"name":"extra","version":"1.0","active":false,"enabled":true}],
		"failed":[{"name":"extra","causes":["java.io.IOException: Failed to load: Extra (extra 1.0)\n - Plugin is missing: gone (2.0)"]}],
		"items":["jobs/app"],
		"unreadable":[{"item":"app","path":"jobs/app","problem":"com.thoughtworks.xstream.mapper.CannotResolveClassException: hudson.plugins.good.Thing"}],
		"monitors":[{"id":"hudson.diagnosis.OldDataMonitor","name":"Manage Old Data"},{"id":"jenkins.security.QueueItemAuthenticatorMonitor","name":"Builds run as SYSTEM"}],
		"probeErrors":["items: boom"]}`)

	good := pluginFile(t, "Manifest-Version: 1.0\nMinimum-Java-Version: 21\n", nil, map[string][]string{"good-core.jar": {"org/example/Foo.class"}})
	srv := newPluginServer(t, map[string][]byte{"/good.hpi": good})

	home := t.TempDir()
	writeTree(t, home, map[string]string{
		"config.xml":                 "<hudson><numExecutors>2</numExecutors></hudson>",
		"jobs/app/config.xml":        `<project><builders><hudson.plugins.good.Thing plugin="good@1.0"/></builders></project>`,
		"jobs/gone-job/config.xml":   `<flow-definition plugin="workflow-job@1385.v1"/>`,
		"jobs/silent/config.xml":     `<project/>`,
		"jobs/app/builds/1/log":      "not configuration",
		"jobs/app/nextBuildNumber":   "2",
		"plugins/good.jpi":           "stale",
		"secrets/master.key":         "m",
		"jenkins.model.Jenkins.yaml": "x",
	})
	casc := filepath.Join(t.TempDir(), "jenkins.yaml")
	os.WriteFile(casc, []byte("jenkins:\n  systemMessage: hi\n"), 0o644)

	out := filepath.Join(t.TempDir(), "out")
	var progress bytes.Buffer
	res, err := Run(context.Background(), Options{
		Core: "2.580.1", Upgrading: "2.462.3", Home: home, CasC: casc,
		Plugins: []plan.Pick{
			{Name: "good", Version: "1.0", URL: srv.URL + "/good.hpi", SHA256: sum(good)},
			{Name: "nourl", Version: "0.1"},
			{Name: "gone", Version: "2.0", URL: srv.URL + "/gone.hpi"},
		},
		OutDir: out, CacheDir: t.TempDir(), Docker: d.path(), Client: srv.Client(), Progress: &progress,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !res.Started || res.Image != "jenkins/jenkins:2.580.1-jdk21" || res.Java != 21 || res.JavaReason != "" || res.Plugins != 1 {
		t.Errorf("result = %+v", res)
	}
	// Plugins with nowhere to download them from come first.
	wantSkipped := []string{"nourl 0.1: no download known", "gone 2.0: GET " + srv.URL + "/gone.hpi: 404 Not Found"}
	if !reflect.DeepEqual(res.Skipped, wantSkipped) {
		t.Errorf("skipped = %q, want %q", res.Skipped, wantSkipped)
	}
	if res.Upgrading != "2.462.3" || !reflect.DeepEqual(res.Detached, []string{"sshd", "instance-identity"}) {
		t.Errorf("upgrading = %q, detached = %q", res.Upgrading, res.Detached)
	}
	if len(res.Failed) != 1 || res.Failed[0].Name != "extra" || !reflect.DeepEqual(res.Failed[0].Reasons, []string{"Plugin is missing: gone (2.0)"}) {
		t.Errorf("failed = %+v", res.Failed)
	}
	if len(res.Roots) != 1 || res.Roots[0].Diagnosis.Kind != CauseMissingDependency {
		t.Errorf("roots = %+v", res.Roots)
	}
	if len(res.Broken) != 1 || res.Broken[0].Name != "good" || res.Broken[0].Version != "1.0" {
		t.Errorf("broken = %+v", res.Broken)
	}
	if !reflect.DeepEqual(res.Errors, []string{"probe: items: boom"}) {
		t.Errorf("errors = %q", res.Errors)
	}
	h := res.Home
	if h == nil || h.Source != home || h.Files != 6 || h.Items != 3 || h.Loaded != 1 {
		t.Fatalf("home = %+v", h)
	}
	wantMissing := []ConfigProblem{
		{Item: "gone-job", Path: "jobs/gone-job", Problem: "CannotResolveClassException: flow-definition", Plugin: "workflow-job 1385.v1", PluginState: PluginNotInstalled},
		// Gone without a word in the log, and its configuration names no
		// plugin.
		{Item: "silent", Path: "jobs/silent", Problem: "not loaded; the log has no reason"},
	}
	if !reflect.DeepEqual(h.Missing, wantMissing) {
		t.Errorf("missing = %+v", h.Missing)
	}
	wantUnreadable := []ConfigProblem{{Item: "app", Path: "jobs/app", Problem: "CannotResolveClassException: hudson.plugins.good.Thing", Plugin: "good 1.0", PluginState: PluginLoaded}}
	if !reflect.DeepEqual(h.Unreadable, wantUnreadable) {
		t.Errorf("unreadable = %+v", h.Unreadable)
	}
	if !reflect.DeepEqual(h.Monitors, []string{"Builds run as SYSTEM"}) {
		t.Errorf("monitors = %q", h.Monitors)
	}
	if res.LogFile != filepath.Join(out, "dry-run-2.580.1.log") {
		t.Errorf("log file = %q", res.LogFile)
	} else if b, _ := os.ReadFile(res.LogFile); !strings.Contains(string(b), "Jenkins is fully up and running") {
		t.Errorf("the log file holds %q", b)
	}

	// What went into the container.
	ref := d.refTar()
	for _, name := range []string{"ref/", "ref/plugins/good.jpi", "ref/init.groovy.d/zz-jenkins-upgrade-check-probe.groovy", "ref/config.xml",
		"ref/jobs/app/config.xml", "ref/secrets/master.key", "ref/jenkins.model.Jenkins.yaml", "ref/casc_configs/jenkins.yaml", "ref/" + lastExecVersion} {
		h, ok := ref[name]
		if !ok {
			t.Errorf("%s not copied in; got %v", name, reflect.ValueOf(ref).MapKeys())
			continue
		}
		if h.Uid != 1000 || h.Gid != 1000 {
			t.Errorf("%s belongs to %d:%d, want the jenkins user", name, h.Uid, h.Gid)
		}
	}
	for _, name := range []string{"ref/jobs/app/builds/1/log", "ref/jobs/app/nextBuildNumber"} {
		if _, ok := ref[name]; ok {
			t.Errorf("%s copied in; it is not configuration", name)
		}
	}
	if got := d.refFile("ref/config.xml"); got != "<hudson><numExecutors>0</numExecutors></hudson>" {
		t.Errorf("config.xml = %q", got)
	}
	if got := d.refFile("ref/" + lastExecVersion); got != "2.462.3" {
		t.Errorf("last version = %q", got)
	}
	if got := d.refFile("ref/plugins/good.jpi"); got != string(good) {
		t.Error("the downloaded plugin, not the stale copy in the home, goes in")
	}

	calls := d.calls()
	create := ""
	for _, c := range calls {
		if strings.HasPrefix(c, "create ") {
			create = c
		}
	}
	for _, want := range []string{"--name juc-2-580-1-", "--network none", "-Djenkins.install.runSetupWizard=false", "-e CASC_JENKINS_CONFIG=/var/jenkins_home/casc_configs", "jenkins/jenkins:2.580.1-jdk21"} {
		if !strings.Contains(create, want) {
			t.Errorf("docker create lacks %q: %s", want, create)
		}
	}
	if !d.called("cp -a - juc-2-580-1-") || !d.called("start juc-") || !d.called("rm -f juc-2-580-1-") {
		t.Errorf("calls = %q", calls)
	}
	for _, want := range []string{"copying the configuration of " + home + " (no build records)", "downloading 2 plugins", "pulling jenkins/jenkins:2.580.1-jdk21", "starting Jenkins 2.580.1 without network (container juc-2-580-1-"} {
		if !strings.Contains(progress.String(), want) {
			t.Errorf("progress lacks %q:\n%s", want, progress.String())
		}
	}
}

// Jenkins exits as it starts. With --keep the container and the work
// directory stay, and the probe's partial record still counts.
func TestRunContainerExits(t *testing.T) {
	fastWaits(t)
	d := newFakeDocker(t)
	d.set("images", "jenkins/jenkins:2.541.3-jdk21\n")
	d.set("state", "false 1\n")
	d.set("jenkins.log", "2026-10-21 13:19:20.000+0000 [id=26]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed to load the plugins\n"+
		"java.lang.LinkageError: boom\n")
	d.set("probe.json", `{"completed":false,"plugins":[{"name":"modern","version":"2.0","active":false,"enabled":true}]}`)
	modern := pluginFile(t, "Minimum-Java-Version: 21\n", nil, nil)
	srv := newPluginServer(t, map[string][]byte{"/modern.hpi": modern})

	res, err := Run(context.Background(), Options{
		Core: "2.541.3", PreferJava: 17, Keep: true, Docker: d.path(), CacheDir: t.TempDir(),
		Plugins: []plan.Pick{{Name: "modern", Version: "2.0", URL: srv.URL + "/modern.hpi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.LogFile != "" {
		t.Cleanup(func() { os.RemoveAll(filepath.Dir(res.LogFile)) })
	}
	if res.Started || res.Reason != "the container exited (state: false 1); the log ends: jenkins.InitReactorRunner$1#onTaskFailed: Failed to load the plugins — java.lang.LinkageError: boom" {
		t.Errorf("started = %v, reason = %q", res.Started, res.Reason)
	}
	// The plan said Java 17; the plugin needs 21.
	if res.Java != 21 || res.JavaReason != "modern needs Java 21 or newer" {
		t.Errorf("java = %d (%q)", res.Java, res.JavaReason)
	}
	want := []Failure{{Name: "modern", Version: "2.0", Reasons: []string{reasonInactive}}}
	if !reflect.DeepEqual(res.Failed, want) {
		t.Errorf("failed = %+v", res.Failed)
	}
	// A rehearsal that did not start does not complain that the probe's
	// list is missing; the log's errors are all there is.
	if !reflect.DeepEqual(res.Errors, []string{"jenkins.InitReactorRunner$1#onTaskFailed: Failed to load the plugins — java.lang.LinkageError: boom"}) {
		t.Errorf("errors = %q", res.Errors)
	}
	if d.called("rm ") {
		t.Error("--keep keeps the container")
	}
	if res.LogFile == "" || filepath.Base(res.LogFile) != "dry-run-2.541.3.log" {
		t.Fatalf("with --keep the log stays in the work directory: %q", res.LogFile)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(res.LogFile), "ref", "plugins", "modern.jpi")); err != nil {
		t.Errorf("--keep keeps the work directory: %v", err)
	}
}

// Jenkins starts but the probe never writes its record: the report says
// only the log was read. Errors past the first forty are counted.
func TestRunWithoutProbe(t *testing.T) {
	fastWaits(t)
	d := newFakeDocker(t)
	d.set("images", "jenkins/jenkins:2.541.3-jdk21\n")
	var log strings.Builder
	for i := 0; i < 45; i++ {
		log.WriteString("2026-10-21 13:19:20.000+0000 [id=26]\tSEVERE\tsome.Thing#run: failure " + string(rune('A'+i)) + "\njava.lang.IllegalStateException: boom\n")
	}
	log.WriteString(upLine)
	d.set("jenkins.log", log.String())
	modern := pluginFile(t, "Minimum-Java-Version: 21\n", nil, nil)
	srv := newPluginServer(t, map[string][]byte{"/modern.hpi": modern})

	res, err := Run(context.Background(), Options{
		Core: "2.541.3", Docker: d.path(), CacheDir: t.TempDir(),
		Plugins: []plan.Pick{{Name: "modern", Version: "2.0", URL: srv.URL + "/modern.hpi"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Started || res.LogFile != "" {
		t.Errorf("started = %v, log file = %q (none without --out-dir or --keep)", res.Started, res.LogFile)
	}
	// No plan: 2.541.3 runs on Java 17, but the plugin needs 21.
	if res.Java != 21 || res.JavaReason != "modern needs Java 21 or newer" {
		t.Errorf("java = %d (%q)", res.Java, res.JavaReason)
	}
	if len(res.Errors) != maxErrors+1 || res.Errors[maxErrors] != "and 6 more errors in the log" {
		t.Errorf("%d errors, last %q", len(res.Errors), res.Errors[len(res.Errors)-1])
	}
	// 45 from the log and the probe's absence make 46.
	if !strings.HasPrefix(res.Errors[0], "some.Thing#run: failure A — java.lang.IllegalStateException: boom") {
		t.Errorf("first error = %q", res.Errors[0])
	}
	n := 0
	for _, c := range d.calls() {
		if strings.HasPrefix(c, "cp juc-") {
			n++
		}
	}
	if n != 31 {
		t.Errorf("read the probe %d times, want 30 tries and one last look", n)
	}
}

func TestRunNotUpInTime(t *testing.T) {
	fastWaits(t)
	d := newFakeDocker(t)
	d.set("images", "jenkins/jenkins:2.580.1-jdk21\n")
	d.set("jenkins.log", "Running from: /usr/share/jenkins/jenkins.war\n")
	res, err := Run(context.Background(), Options{Core: "2.580.1", Timeout: time.Nanosecond, Docker: d.path()})
	if err != nil {
		t.Fatal(err)
	}
	if res.Started || res.Reason != "not up after 1ns" {
		t.Errorf("started = %v, reason = %q", res.Started, res.Reason)
	}
}

// cancelOn cancels a context when a progress line contains text.
type cancelOn struct {
	text   string
	cancel context.CancelFunc
}

func (c cancelOn) Write(b []byte) (int, error) {
	if strings.Contains(string(b), c.text) {
		c.cancel()
	}
	return len(b), nil
}

// Each step before Jenkins runs says what failed.
func TestRunErrors(t *testing.T) {
	fastWaits(t)
	ctx := context.Background()
	if _, err := Run(ctx, Options{}); err == nil || err.Error() != "no Jenkins version to rehearse" {
		t.Errorf("no core: %v", err)
	}
	if _, err := Run(ctx, Options{Core: "2.580.1", Docker: filepath.Join(t.TempDir(), "docker")}); err == nil || !strings.HasPrefix(err.Error(), "dry-run needs Docker: ") {
		t.Errorf("no docker: %v", err)
	}

	d := newFakeDocker(t)
	t.Run("temp dir", func(t *testing.T) {
		t.Setenv("TMPDIR", filepath.Join(t.TempDir(), "missing"))
		if _, err := Run(ctx, Options{Core: "2.580.1", Docker: d.path()}); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("err = %v", err)
		}
	})
	notHome := t.TempDir()
	if _, err := Run(ctx, Options{Core: "2.580.1", Docker: d.path(), Home: notHome}); err == nil || err.Error() != notHome+" has no config.xml; is it a JENKINS_HOME?" {
		t.Errorf("not a home: %v", err)
	}
	if _, err := Run(ctx, Options{Core: "2.580.1", Docker: d.path(), CasC: filepath.Join(notHome, "casc.yaml")}); err == nil || !strings.HasPrefix(err.Error(), "--casc: stat ") {
		t.Errorf("missing CasC: %v", err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := Run(cancelled, Options{Core: "2.580.1", Docker: d.path()}); !errors.Is(err, context.Canceled) {
		t.Errorf("cancelled: %v", err)
	}
	if _, err := os.Stat(filepath.Join(d.dir, "calls")); err == nil {
		t.Errorf("a rehearsal that failed before the container ran docker: %q", d.calls())
	}

	if _, err := Run(ctx, Options{Core: "2.580.1", Docker: d.path()}); err == nil || err.Error() != "no Jenkins image found for 2.580.1 (tried jenkins/jenkins:2.580.1-jdk21, jenkins/jenkins:2.580.1-jdk25, jenkins/jenkins:2.580.1): Error response from daemon: manifest for jenkins/jenkins:2.580.1 not found: manifest unknown" {
		t.Errorf("no image: %v", err)
	}
	d.set("images", "jenkins/jenkins:2.580.1-jdk21\n")

	d.set("create.fail", "Error response from daemon: Conflict.")
	if _, err := Run(ctx, Options{Core: "2.580.1", Docker: d.path()}); err == nil || err.Error() != "docker create: exit status 1\nError response from daemon: Conflict." {
		t.Errorf("create: %v", err)
	}
	os.Remove(filepath.Join(d.dir, "create.fail"))

	d.set("cp-in.fail", "no space left on device")
	if _, err := Run(ctx, Options{Core: "2.580.1", Docker: d.path()}); err == nil || err.Error() != "docker cp: exit status 1\nno space left on device" {
		t.Errorf("copy in: %v", err)
	}
	os.Remove(filepath.Join(d.dir, "cp-in.fail"))

	d.set("start.fail", "Error response from daemon: OCI runtime create failed")
	if _, err := Run(ctx, Options{Core: "2.580.1", Docker: d.path()}); err == nil || err.Error() != "docker start: exit status 1\nError response from daemon: OCI runtime create failed" {
		t.Errorf("start: %v", err)
	}
	os.Remove(filepath.Join(d.dir, "start.fail"))
	if !d.called("rm -f juc-") {
		t.Error("a container that did not start is removed")
	}
}

// Ctrl-C while Jenkins starts, or while waiting for the probe, stops the
// rehearsal.
func TestRunCancelled(t *testing.T) {
	old := [2]time.Duration{startPoll, probePoll}
	startPoll, probePoll = time.Hour, time.Hour
	t.Cleanup(func() { startPoll, probePoll = old[0], old[1] })

	d := newFakeDocker(t)
	d.set("images", "jenkins/jenkins:2.580.1-jdk21\n")
	d.set("jenkins.log", "Running from: /usr/share/jenkins/jenkins.war\n")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := Run(ctx, Options{Core: "2.580.1", Docker: d.path(), Progress: cancelOn{"starting Jenkins", cancel}})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("while starting: %v", err)
	}

	// Up at once, but the probe has not written its record: cancel once
	// the first read of it has begun.
	d = newFakeDocker(t)
	d.set("images", "jenkins/jenkins:2.580.1-jdk21\n")
	d.set("jenkins.log", upLine)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			if _, err := os.Stat(filepath.Join(d.dir, "probe-read")); err == nil {
				cancel()
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond):
			}
		}
	}()
	_, err = Run(ctx, Options{Core: "2.580.1", Docker: d.path()})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("while waiting for the probe: %v", err)
	}
}
