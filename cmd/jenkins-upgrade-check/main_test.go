package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/keelinfra/jenkins-upgrade-check/internal/dryrun"
	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
	"github.com/keelinfra/jenkins-upgrade-check/internal/snapshot"
	"github.com/keelinfra/jenkins-upgrade-check/internal/verify"
)

// ucDir holds update-center.actual.json and plugin-versions.json for
// --update-center-dir.
const ucDir = "../../internal/plan/testdata"

// TestMain keeps the default cache directory out of the user's home.
func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "juc-cmd-test-")
	if err != nil {
		panic(err)
	}
	os.Setenv("XDG_CACHE_HOME", filepath.Join(tmp, "cache"))
	os.Setenv("HOME", filepath.Join(tmp, "home"))
	code := m.Run()
	os.RemoveAll(tmp)
	os.Exit(code)
}

type result struct {
	code           int
	stdout, stderr string
}

// cli runs the command line with the environment given.
func cli(env map[string]string, args ...string) result {
	var out, errOut bytes.Buffer
	code := run(args, &out, &errOut, func(k string) string { return env[k] })
	return result{code, out.String(), errOut.String()}
}

func (r result) wantCode(t *testing.T, code int) {
	t.Helper()
	if r.code != code {
		t.Fatalf("exit %d, want %d\nstdout:\n%s\nstderr:\n%s", r.code, code, r.stdout, r.stderr)
	}
}

// wantError checks a usage error: exit 2, and msg in the "error:" line
// that ends stderr.
func (r result) wantError(t *testing.T, msg string) {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(r.stderr, "\n"), "\n")
	last := lines[len(lines)-1]
	if r.code != 2 || !strings.HasPrefix(last, "error: ") || !strings.Contains(last, msg) {
		t.Errorf("exit %d, stderr %q; want 2 and an error line with %q", r.code, r.stderr, msg)
	}
}

func contains(t *testing.T, what, s string, wants ...string) {
	t.Helper()
	for _, w := range wants {
		if !strings.Contains(s, w) {
			t.Errorf("%s lacks %q:\n%s", what, w, s)
		}
	}
}

func write(t *testing.T, path, body string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func writeJSONFile(t *testing.T, path string, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return write(t, path, string(b))
}

func TestUsage(t *testing.T) {
	r := cli(nil)
	r.wantCode(t, 2)
	contains(t, "stderr", r.stderr, "Usage:\n  jenkins-upgrade-check collect")

	for _, arg := range []string{"help", "-h", "--help"} {
		r = cli(nil, arg)
		r.wantCode(t, 0)
		contains(t, "stdout", r.stdout, "Run \"jenkins-upgrade-check COMMAND -h\" for the flags of one command.")
	}
	for _, arg := range []string{"version", "--version"} {
		if r = cli(nil, arg); r.code != 0 || r.stdout != "jenkins-upgrade-check dev\n" {
			t.Errorf("%s: %+v", arg, r)
		}
	}
	r = cli(nil, "upgrade")
	r.wantCode(t, 2)
	if !strings.HasPrefix(r.stderr, "unknown command \"upgrade\"\n\njenkins-upgrade-check plans, rehearses and verifies a Jenkins upgrade.\n") {
		t.Errorf("stderr = %q", r.stderr)
	}

	// -h on a command prints its flags, and is not an error.
	r = cli(nil, "plan", "-h")
	r.wantCode(t, 2)
	contains(t, "stderr", r.stderr, "Usage of plan:", "-runbook string")
	if strings.Contains(r.stderr, "error:") {
		t.Errorf("-h is reported as an error: %q", r.stderr)
	}
	for _, cmd := range []string{"collect", "plan", "dry-run", "download", "probe", "verify", "resolve"} {
		cli(nil, cmd, "--bogus").wantError(t, "flag provided but not defined: -bogus")
	}
}

// controller is a fake Jenkins: core, plugins, jobs and agents as given,
// for the user audit with the token tok.
type controller struct {
	core    string
	plugins string // JSON array
	agents  string // JSON array, after the built-in node
	jobs    string // JSON array
}

func (c controller) server(t *testing.T) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Jenkins", c.core)
		if u, p, _ := r.BasicAuth(); u != "audit" || p != "tok" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch {
		case r.URL.Path == "/api/json":
			w.Write([]byte(`{"jobs":` + c.jobs + `}`))
		case r.URL.Path == "/pluginManager/api/json":
			w.Write([]byte(`{"plugins":` + c.plugins + `}`))
		case r.URL.Path == "/computer/api/json" && c.agents == "forbidden":
			w.WriteHeader(http.StatusForbidden)
		case r.URL.Path == "/computer/api/json":
			w.Write([]byte(`{"computer":[{"_class":"hudson.model.Hudson$MasterComputer","displayName":"Built-In Node","monitorData":{"hudson.plugin.versioncolumn.JVMVersionMonitor":"21.0.5"}}` + c.agents + `]}`))
		case strings.HasSuffix(r.URL.Path, "/logText/progressiveText"):
			w.Write([]byte("SEVERE: java.lang.UnsupportedClassVersionError: hudson/remoting/Launcher (class file version 65.0), this version of the Java Runtime only recognizes class file versions up to 61.0\n"))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

var jenkinsEnv = map[string]string{"JENKINS_USER": "audit", "JENKINS_TOKEN": "tok"}

func withURL(url string) map[string]string {
	env := map[string]string{"JENKINS_URL": url}
	for k, v := range jenkinsEnv {
		env[k] = v
	}
	return env
}

func TestCollect(t *testing.T) {
	srv := controller{core: "2.462.3", jobs: `[{"fullName":"app"}]`, agents: "forbidden",
		plugins: `[{"shortName":"git","version":"5.2.1","enabled":true,"active":true}]`}.server(t)
	out := filepath.Join(t.TempDir(), "snapshot.json")

	r := cli(withURL(srv.URL), "collect", "-o", out)
	r.wantCode(t, 0)
	if want := "Jenkins 2.462.3: 1 plugins, 0 agents → " + out + " (3 GET requests, nothing written to the controller)\nnote: agents could not be read ("; !strings.HasPrefix(r.stderr, want) {
		t.Errorf("stderr = %q, want %q…", r.stderr, want)
	}
	s, err := snapshot.Read(out)
	if err != nil {
		t.Fatal(err)
	}
	if s.Core != "2.462.3" || s.Source != srv.URL || s.Tool != "jenkins-upgrade-check" || len(s.Plugins) != 1 || s.Items[0] != "app" {
		t.Errorf("snapshot = %+v", s)
	}

	// Flags instead of the environment; --token wins over JENKINS_TOKEN.
	r = cli(map[string]string{"JENKINS_TOKEN": "wrong"}, "collect", "--url", srv.URL, "--user", "audit", "--token", "tok", "-o", out)
	r.wantCode(t, 0)

	cli(nil, "collect").wantError(t, "--url is required")
	cli(map[string]string{"JENKINS_URL": srv.URL, "JENKINS_USER": "audit", "JENKINS_TOKEN": "bad"}, "collect", "-o", out).
		wantError(t, "the controller answered but refused the API token: refused: GET /api/json")
	cli(withURL(srv.URL), "collect", "-o", filepath.Join(t.TempDir(), "missing", "s.json")).wantError(t, "open ")
}

func TestPlan(t *testing.T) {
	dir := t.TempDir()
	plugins := write(t, filepath.Join(dir, "plugins.txt"), "workflow-cps:2.90\nscript-security:1.75\nold-thing:1.0\n")
	jsonOut, picksOut, rb, md := filepath.Join(dir, "plan.json"), filepath.Join(dir, "target.txt"), filepath.Join(dir, "runbook.md"), filepath.Join(dir, "plan.md")

	r := cli(nil, "plan", "--plugins-file", plugins, "--core", "2.346.3", "--to", "2.555.1", "--update-center-dir", ucDir,
		"--json", jsonOut, "--plugins-out", picksOut, "--runbook", rb, "--install", "deb", "-o", md)
	r.wantCode(t, 0)
	if r.stdout != "" || r.stderr != "reading update center data\n" {
		t.Errorf("stdout = %q, stderr = %q", r.stdout, r.stderr)
	}
	contains(t, "report", read(t, md), "# Jenkins upgrade plan: 2.346.3 → 2.555.1", "Source: "+plugins, "SECURITY-1", "| old-thing | 1.0 | 1.0 |")
	contains(t, "runbook", read(t, rb), "# Upgrade runbook: Jenkins 2.346.3 → 2.555.1", "Made from `"+jsonOut+"`", "JENKINS_URL=https://jenkins.example.com", "jenkins=2.555.1")
	if got := read(t, picksOut); got != "old-thing:1.0\nscript-security:1350.v1\nworkflow-api:1300.v1\nworkflow-cps:4000.v1\n" {
		t.Errorf("plugins-out = %q", got)
	}
	p, err := readPlan(jsonOut)
	if err != nil {
		t.Fatal(err)
	}
	if p.Tool != "jenkins-upgrade-check" || p.ToolVersion != "dev" || p.From != "2.346.3" || p.To != "2.555.1" || len(p.Steps) != 3 {
		t.Errorf("plan = %+v", p)
	}

	// From a snapshot, with the core given again; the runbook goes to
	// stdout and knows the controller's URL.
	snap := writeJSONFile(t, filepath.Join(dir, "snapshot.json"), snapshot.Snapshot{Source: "https://ci.example.com", Core: "2.346.1",
		Plugins: []snapshot.Plugin{{Name: "workflow-cps", Version: "2.90", Enabled: true, Active: true}}})
	r = cli(nil, "plan", "--snapshot", snap, "--core", "2.346.3", "--to", "2.555.1", "--update-center-dir", ucDir, "--runbook", "-", "--install", "docker", "-o", md)
	r.wantCode(t, 0)
	contains(t, "runbook", r.stdout, "JENKINS_URL=https://ci.example.com", "Made from `plan.json`", "docker pull jenkins/jenkins:2.555.1-jdk21")
	contains(t, "report", read(t, md), "# Jenkins upgrade plan: 2.346.3 → 2.555.1")

	// The report goes to stdout by default.
	r = cli(nil, "plan", "--snapshot", snap, "--to", "2.555.1", "--update-center-dir", ucDir)
	r.wantCode(t, 0)
	contains(t, "stdout", r.stdout, "# Jenkins upgrade plan: 2.346.1 → 2.555.1")
}

// lts and weekly come from the update center, as do the files, which are
// cached.
func TestPlanFromNetwork(t *testing.T) {
	files := map[string]string{
		"/stable/latestCore.txt":             "2.555.1\n",
		"/current/plugin-versions.json":      read(t, filepath.Join(ucDir, "plugin-versions.json")),
		"/current/update-center.actual.json": read(t, filepath.Join(ucDir, "update-center.actual.json")),
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ua := r.Header.Get("User-Agent"); ua != "jenkins-upgrade-check/dev" {
			t.Errorf("User-Agent = %q", ua)
		}
		b, ok := files[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(b))
	}))
	defer srv.Close()
	dir := t.TempDir()
	plugins := write(t, filepath.Join(dir, "plugins.txt"), "workflow-cps:2.90\n")
	cache := filepath.Join(dir, "cache")
	r := cli(nil, "plan", "--plugins-file", plugins, "--core", "2.462.3", "--update-center", srv.URL, "--cache-dir", cache)
	r.wantCode(t, 0)
	contains(t, "stdout", r.stdout, "# Jenkins upgrade plan: 2.462.3 → 2.555.1")
	if _, err := os.Stat(filepath.Join(cache, "current_plugin-versions.json")); err != nil {
		t.Errorf("not cached: %v", err)
	}
}

func TestPlanErrors(t *testing.T) {
	dir := t.TempDir()
	plugins := write(t, filepath.Join(dir, "plugins.txt"), "workflow-cps:2.90\n")
	uc := []string{"--update-center-dir", ucDir}
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--runbook", "rb.md"}, "--runbook needs --install, one of deb, rpm, war, docker"},
		{[]string{"--snapshot", "s.json", "--plugins-file", plugins}, "give --snapshot or --plugins-file, not both"},
		{[]string{"--plugins-file", plugins}, "--plugins-file needs --core, the Jenkins version those plugins run on"},
		{[]string{}, "give --snapshot (from collect) or --plugins-file with --core"},
		{[]string{"--snapshot", filepath.Join(dir, "missing.json")}, "missing.json: no such file or directory"},
		{append([]string{"--plugins-file", plugins, "--core", "2.346.3", "--to", "lts"}, uc...), "--to lts needs the network; give a version number with --update-center-dir"},
		{[]string{"--plugins-file", plugins, "--core", "2.346.3", "--to", "2.555.1", "--update-center-dir", dir}, "plugin-versions.json: no such file or directory (with --update-center-dir, download it from"},
		{append([]string{"--plugins-file", plugins, "--core", "2.555.1", "--to", "2.541.3"}, uc...), "the controller already runs 2.555.1, which is not older than 2.541.3"},
		{append([]string{"--plugins-file", plugins, "--core", "2.346.3", "--to", "2.555.1", "--json", filepath.Join(dir, "no", "p.json")}, uc...), "no such file or directory"},
		{append([]string{"--plugins-file", plugins, "--core", "2.346.3", "--to", "2.555.1", "--plugins-out", filepath.Join(dir, "no", "p.txt")}, uc...), "no such file or directory"},
		{append([]string{"--plugins-file", plugins, "--core", "2.346.3", "--to", "2.555.1", "--runbook", filepath.Join(dir, "no", "rb.md"), "--install", "deb"}, uc...), "no such file or directory"},
		{append([]string{"--plugins-file", plugins, "--core", "2.346.3", "--to", "2.555.1", "--runbook", filepath.Join(dir, "rb.md"), "--install", "msi"}, uc...), `install type "msi": use one of deb, rpm, war, docker`},
	} {
		cli(nil, append([]string{"plan"}, c.args...)...).wantError(t, c.want)
	}
	// An update center with update-center.actual.json missing.
	only := t.TempDir()
	write(t, filepath.Join(only, "plugin-versions.json"), read(t, filepath.Join(ucDir, "plugin-versions.json")))
	cli(nil, "plan", "--plugins-file", plugins, "--core", "2.346.3", "--to", "2.555.1", "--update-center-dir", only).
		wantError(t, "update-center.actual.json: no such file or directory")
}

func TestResolve(t *testing.T) {
	dir := t.TempDir()
	uc := []string{"--update-center-dir", ucDir}

	r := cli(nil, append([]string{"resolve", "--core", "2.479.3", "--names", "workflow-cps, ,gone"}, uc...)...)
	r.wantCode(t, 0)
	if r.stdout != "script-security:1350.v1\nworkflow-api:1300.v1\nworkflow-cps:4000.v1\n" {
		t.Errorf("stdout = %q", r.stdout)
	}
	if r.stderr != "warning: gone: the update center no longer distributes it and no version was given; left out\n3 plugins for Jenkins 2.479.3\n" {
		t.Errorf("stderr = %q", r.stderr)
	}

	plugins := write(t, filepath.Join(dir, "plugins.txt"), "workflow-cps:2.90\n")
	out := filepath.Join(dir, "out.txt")
	r = cli(nil, append([]string{"resolve", "--core", "2.462.3", "--plugins-file", plugins, "-o", out}, uc...)...)
	r.wantCode(t, 0)
	if got := read(t, out); got != "script-security:1300.v1\nworkflow-api:1200.v1\nworkflow-cps:3800.v1\n" {
		t.Errorf("out = %q", got)
	}

	// The most installed plugins that have a release for the core: the
	// header, a row without a count and a plugin too new for the core are
	// left out.
	csv := write(t, filepath.Join(dir, "top.csv"), "\"name\",\"installations\"\n\"needs-future\",\"900\"\n\"ionicons-api\",\"100\"\n\"lonely\",\"n/a\"\n\"workflow-api\",\"500\"\n\"script-security\",\"700\"\n")
	r = cli(nil, append([]string{"resolve", "--core", "2.479.3", "--top-plugins", csv, "--top", "2"}, uc...)...)
	r.wantCode(t, 0)
	if r.stdout != "script-security:1350.v1\nworkflow-api:1300.v1\n" {
		t.Errorf("top 2 = %q", r.stdout)
	}
	// Rows of one field are no plugin counts.
	r = cli(nil, append([]string{"resolve", "--core", "2.479.3", "--top-plugins", write(t, filepath.Join(dir, "names.csv"), "workflow-api\nscript-security\n")}, uc...)...)
	r.wantCode(t, 0)
	if r.stdout != "" || r.stderr != "0 plugins for Jenkins 2.479.3\n" {
		t.Errorf("one-field rows: stdout %q, stderr %q", r.stdout, r.stderr)
	}

	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{}, "--core is required"},
		{[]string{"--core", "weekly", "--names", "x", "--update-center-dir", ucDir}, "--to weekly needs the network"},
		{[]string{"--core", "2.479.3", "--names", "x", "--update-center-dir", dir}, "plugin-versions.json: no such file or directory"},
		{[]string{"--core", "2.479.3", "--update-center-dir", ucDir}, "give --plugins-file, --names or --top-plugins"},
		{[]string{"--core", "2.479.3", "--plugins-file", filepath.Join(dir, "none.txt"), "--update-center-dir", ucDir}, "none.txt: no such file or directory"},
		{[]string{"--core", "2.479.3", "--top-plugins", filepath.Join(dir, "none.csv"), "--update-center-dir", ucDir}, "none.csv: no such file or directory"},
		{[]string{"--core", "2.479.3", "--top-plugins", write(t, filepath.Join(dir, "bad.csv"), "\"a\"b,1\n"), "--update-center-dir", ucDir}, "bad.csv: parse error on line 1"},
		{[]string{"--core", "2.479.3", "--names", "x", "-o", filepath.Join(dir, "no", "out.txt"), "--update-center-dir", ucDir}, "no such file or directory"},
	} {
		cli(nil, append([]string{"resolve"}, c.args...)...).wantError(t, c.want)
	}
}

func TestProbe(t *testing.T) {
	r := cli(nil, "probe")
	r.wantCode(t, 0)
	if r.stdout != string(dryrun.ProbeScript(false)) {
		t.Error("probe prints the probe script")
	}
	r = cli(nil, "probe", "--quiet")
	r.wantCode(t, 0)
	if !strings.HasPrefix(r.stdout, "// Start quiet:") {
		t.Errorf("probe --quiet starts %q", r.stdout[:40])
	}
	var errOut bytes.Buffer
	if code := run([]string{"probe"}, failingWriter{}, &errOut, os.Getenv); code != 2 || errOut.String() != "error: disk full\n" {
		t.Errorf("exit %d, stderr %q", code, errOut.String())
	}
	errOut.Reset()
	if code := run([]string{"resolve", "--core", "2.479.3", "--names", "workflow-api", "--update-center-dir", ucDir}, failingWriter{}, &errOut, os.Getenv); code != 2 || !strings.HasSuffix(errOut.String(), "\nerror: disk full\n") {
		t.Errorf("resolve to a full disk: exit %d, stderr %q", code, errOut.String())
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("disk full") }

func sum(b []byte) string {
	s := sha256.Sum256(b)
	return base64.StdEncoding.EncodeToString(s[:])
}

// pluginServer serves plugin files by name, as /NAME.hpi.
func pluginServer(t *testing.T, files map[string][]byte) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, ok := files[strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/"), ".hpi")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Write(b)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// testPlan has a plugin set for each step, downloadable from srv.
func testPlan(t *testing.T, srv string, files map[string][]byte) string {
	pick := func(name string) plan.Pick {
		return plan.Pick{Name: name, Version: "1.0", URL: srv + "/" + name + ".hpi", SHA256: sum(files[name])}
	}
	p := plan.Plan{From: "2.462.3", To: "2.555.1",
		Steps: []plan.Step{
			{Core: "2.462.3", Plugins: []plan.Pick{pick("git")}},
			{Core: "2.541.3", Java: 17, Upgrade: true, Plugins: []plan.Pick{pick("git"), pick("mailer")}},
			{Core: "2.555.1", Java: 21, Upgrade: true},
		},
		Target: plan.Resolution{Core: "2.555.1", Plugins: []plan.Pick{pick("git"), pick("mailer"), {Name: "private", Version: "0.1"}}},
	}
	return writeJSONFile(t, filepath.Join(t.TempDir(), "plan.json"), p)
}

func TestDownload(t *testing.T) {
	files := map[string][]byte{"git": []byte("git plugin"), "mailer": []byte("mailer plugin")}
	srv := pluginServer(t, files)
	pf := testPlan(t, srv.URL, files)
	dir := filepath.Join(t.TempDir(), "step1")
	cache := t.TempDir()

	r := cli(nil, "download", "--plan", pf, "--step", "1", "--dir", dir, "--cache-dir", cache)
	r.wantCode(t, 0)
	if r.stderr != "downloading 2 plugins for Jenkins 2.541.3 into "+dir+"\n2 of 2 plugins downloaded\n" {
		t.Errorf("stderr = %q", r.stderr)
	}
	if read(t, filepath.Join(dir, "git.jpi")) != "git plugin" || read(t, filepath.Join(dir, "mailer.jpi")) != "mailer plugin" {
		t.Error("plugins not downloaded")
	}
	if _, err := os.Stat(filepath.Join(cache, "plugins", "git", "1.0", "git.hpi")); err != nil {
		t.Errorf("not cached under the cache directory: %v", err)
	}

	// The target has a plugin with nowhere to download it from.
	r = cli(nil, "download", "--plan", pf, "--dir", t.TempDir(), "--cache-dir", cache)
	r.wantError(t, "1 plugin(s) could not be downloaded; the plugins directory would be incomplete")
	contains(t, "stderr", r.stderr, "not downloaded: private 0.1: no download known\n2 of 3 plugins downloaded\n")

	file := write(t, filepath.Join(t.TempDir(), "file"), "")
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--plan", pf}, "--plan and --dir are required"},
		{[]string{"--plan", filepath.Join(t.TempDir(), "none.json"), "--dir", dir}, "none.json: no such file or directory"},
		{[]string{"--plan", write(t, filepath.Join(t.TempDir(), "bad.json"), "{"), "--dir", dir}, "bad.json: unexpected end of JSON input"},
		{[]string{"--plan", pf, "--step", "2", "--dir", dir}, "step 2 has no plugin set; make the plan again with this version"},
		{[]string{"--plan", pf, "--step", "1", "--dir", filepath.Join(file, "x")}, "not a directory"},
	} {
		cli(nil, append([]string{"download"}, c.args...)...).wantError(t, c.want)
	}
}

// fakeDocker puts a docker on the PATH that starts "Jenkins" at once,
// printing jenkins.log and leaving probe.json in its directory as the
// probe's record.
func fakeDocker(t *testing.T, log, probe string) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "jenkins.log"), log)
	write(t, filepath.Join(dir, "probe.json"), probe)
	script := `#!/bin/sh
d=$(dirname "$0")
case "$1" in
create) echo 0123abcd ;;
cp)
	if [ "$2" = "-a" ]; then cat > /dev/null; exit 0; fi
	cp "$d/probe.json" "$3" ;;
logs) cat "$d/jenkins.log" ;;
inspect) echo "true 0" ;;
esac
exit 0
`
	if err := os.WriteFile(filepath.Join(dir, "docker"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

const upLine = "2026-10-21 13:36:10.000+0000 [id=25]\tINFO\thudson.lifecycle.Lifecycle#onReady: Jenkins is fully up and running\n"

func TestDryRunPlan(t *testing.T) {
	files := map[string][]byte{"git": []byte("git plugin"), "mailer": []byte("mailer plugin")}
	srv := pluginServer(t, files)
	pf := testPlan(t, srv.URL, files)
	fakeDocker(t, upLine, `{"completed":true,"plugins":[{"name":"git","version":"1.0","active":true,"enabled":true},{"name":"mailer","version":"1.0","active":true,"enabled":true}]}`)
	dir := t.TempDir()
	jsonOut := filepath.Join(dir, "result.json")

	// Step 1 is rehearsed as an upgrade from step 0's core, on the plan's Java.
	r := cli(nil, "dry-run", "--plan", pf, "--step", "1", "--cache-dir", t.TempDir(), "--json", jsonOut)
	r.wantCode(t, 0)
	contains(t, "report", r.stdout, "# Dry run: Jenkins 2.541.3 with 2 plugins", "Image `jenkins/jenkins:2.541.3-jdk17`",
		"Jenkins started as an upgrade from 2.462.3.", "and every plugin loaded.**")
	contains(t, "progress", r.stderr, "downloading 2 plugins", "pulling jenkins/jenkins:2.541.3-jdk17", "starting Jenkins 2.541.3 without network")
	var res dryrun.Result
	if err := json.Unmarshal([]byte(read(t, jsonOut)), &res); err != nil {
		t.Fatal(err)
	}
	if !res.OK() || res.Core != "2.541.3" || res.Upgrading != "2.462.3" || res.Java != 17 {
		t.Errorf("result = %+v", res)
	}

	// The target, from a given version; a plugin that cannot be downloaded
	// is reported, and is no failure of the rehearsal.
	out := filepath.Join(dir, "report.md")
	r = cli(nil, "dryrun", "--plan", pf, "--from", "2.346.3", "--cache-dir", t.TempDir(), "-o", out)
	r.wantCode(t, 0)
	contains(t, "report", read(t, out), "# Dry run: Jenkins 2.555.1 with 2 plugins", "Jenkins started as an upgrade from 2.346.3.", "Not rehearsed (could not be downloaded): private 0.1: no download known.")

	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--plan", filepath.Join(dir, "none.json")}, "none.json: no such file or directory"},
		{[]string{"--plan", write(t, filepath.Join(dir, "bad.json"), "[]")}, "bad.json: json: cannot unmarshal array"},
		{[]string{"--plan", pf, "--step", "7"}, "the plan has steps 0 to 2"},
		{[]string{"--plan", pf, "--cache-dir", t.TempDir(), "--json", filepath.Join(dir, "no", "r.json")}, "no such file or directory"},
		{[]string{"--plan", pf, "--cache-dir", t.TempDir(), "-o", filepath.Join(dir, "no", "r.md")}, "no such file or directory"},
	} {
		cli(nil, append([]string{"dry-run"}, c.args...)...).wantError(t, c.want)
	}

	t.Setenv("PATH", t.TempDir())
	cli(nil, "dry-run", "--plan", pf).wantError(t, "dry-run needs Docker: ")
}

// ucFiles writes an update center whose plugins download from srv.
func ucFiles(t *testing.T, srv string, files map[string][]byte) string {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "update-center.actual.json"), `{"plugins":{}}`)
	write(t, filepath.Join(dir, "plugin-versions.json"), `{"plugins":{
		"git":{"1.0":{"requiredCore":"2.400","url":"`+srv+`/git.hpi","sha256":"`+sum(files["git"])+`","dependencies":[{"name":"mailer","version":"1.0","optional":false}]}},
		"mailer":{"1.0":{"requiredCore":"2.400","url":"`+srv+`/mailer.hpi","sha256":"`+sum(files["mailer"])+`"}},
		"future":{"9.0":{"requiredCore":"9.0","url":"`+srv+`/future.hpi","sha256":"`+sum(files["future"])+`"}}
	}}`)
	return dir
}

func TestDryRunPluginsFile(t *testing.T) {
	files := map[string][]byte{"git": []byte("git plugin"), "mailer": []byte("mailer plugin"), "future": []byte("future plugin")}
	srv := pluginServer(t, files)
	uc := ucFiles(t, srv.URL, files)
	plugins := write(t, filepath.Join(t.TempDir(), "plugins.txt"), "git:1.0\n")
	log := "2026-10-21 12:41:52.000+0000 [id=41]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed Loading plugin Git v1.0 (git)\n" +
		"java.io.IOException: Failed to load: Git (git 1.0)\n" +
		" - Plugin is missing: mailer (1.0)\n" + upLine
	fakeDocker(t, log, `{"completed":true,"plugins":[{"name":"mailer","version":"1.0","active":true,"enabled":true}]}`)

	// As given: git's dependency is missing, so git fails and the exit
	// code says so.
	r := cli(nil, "dry-run", "--plugins-file", plugins, "--core", "2.462.3", "--update-center-dir", uc, "--cache-dir", t.TempDir(), "--timeout", "1m")
	r.wantCode(t, 1)
	contains(t, "report", r.stdout, "# Dry run: Jenkins 2.462.3 with 1 plugins", "but 1 plugin(s) did not load",
		"| git 1.0 | needs the mailer plugin (1.0 or newer), which is not installed | install mailer 1.0 or newer |")

	// --resolve adds the dependency, and says what it could not settle.
	fakeDocker(t, upLine, `{"completed":true,"plugins":[{"name":"git","version":"1.0","active":true,"enabled":true},{"name":"mailer","version":"1.0","active":true,"enabled":true},{"name":"future","version":"9.0","active":true,"enabled":true}]}`)
	plugins = write(t, filepath.Join(t.TempDir(), "plugins.txt"), "git:1.0\nfuture:9.0\n")
	r = cli(nil, "dry-run", "--plugins-file", plugins, "--core", "2.462.3", "--resolve", "--update-center-dir", uc, "--cache-dir", t.TempDir())
	r.wantCode(t, 0)
	contains(t, "report", r.stdout, "# Dry run: Jenkins 2.462.3 with 3 plugins", "every plugin loaded")
	contains(t, "stderr", r.stderr, "warning: future: no release runs on Jenkins 2.462.3; kept at 9.0\n")

	cli(nil, "dry-run", "--plugins-file", plugins).wantError(t, "--plugins-file needs --core")
	cli(nil, "dry-run", "--plugins-file", plugins, "--core", "2.462.3", "--update-center-dir", t.TempDir()).wantError(t, "plugin-versions.json: no such file or directory")
}

// before is the controller as collected before the upgrade.
func before(t *testing.T) string {
	return writeJSONFile(t, filepath.Join(t.TempDir(), "before.json"), snapshot.Snapshot{Core: "2.462.3",
		Plugins: []snapshot.Plugin{{Name: "git", Version: "5.2.1", Enabled: true, Active: true}},
		Items:   []string{"app"},
		Agents:  []snapshot.Agent{{Name: "linux1"}}})
}

func TestVerifyGo(t *testing.T) {
	srv := controller{core: "2.555.1", jobs: `[{"fullName":"app"}]`,
		plugins: `[{"shortName":"git","version":"5.8","enabled":true,"active":true}]`,
		agents:  `,{"_class":"hudson.slaves.SlaveComputer","displayName":"linux1","offline":false}`}.server(t)
	pf := writeJSONFile(t, filepath.Join(t.TempDir(), "plan.json"), plan.Plan{From: "2.462.3", To: "2.555.1",
		Steps:  []plan.Step{{Core: "2.462.3"}, {Core: "2.555.1", Upgrade: true, Plugins: []plan.Pick{{Name: "git", Version: "5.8"}}}},
		Target: plan.Resolution{Plugins: []plan.Pick{{Name: "git", Version: "9.9"}}}})
	home := t.TempDir()
	write(t, filepath.Join(home, dryrun.ProbeFile), `{"completed":true,"version":"2.555.1","monitors":[{"id":"x","name":"Builds run as SYSTEM"}]}`)
	logFile := write(t, filepath.Join(t.TempDir(), "start.log"), upLine)
	dir := t.TempDir()
	jsonOut, md := filepath.Join(dir, "verify.json"), filepath.Join(dir, "verify.md")

	r := cli(withURL(srv.URL), "verify", "--before", before(t), "--plan", pf, "--home", home, "--log", logFile, "--json", jsonOut, "-o", md)
	r.wantCode(t, 0)
	if r.stderr != "GO\n" {
		t.Errorf("stderr = %q", r.stderr)
	}
	contains(t, "report", read(t, md), "# Upgrade check: Jenkins 2.555.1", "## GO", "- Builds run as SYSTEM")
	var res verify.Result
	if err := json.Unmarshal([]byte(read(t, jsonOut)), &res); err != nil {
		t.Fatal(err)
	}
	// The step's plugin set is what the upgrade meant to install, not the
	// target's: no version warning.
	if res.Verdict != verify.Go || len(res.Checks) != 5 {
		t.Errorf("result = %+v", res)
	}

	// A probe record from another version is ignored, and a missing one
	// says how to install it.
	write(t, filepath.Join(home, dryrun.ProbeFile), `{"completed":true,"version":"2.541.3"}`)
	r = cli(withURL(srv.URL), "verify", "--before", before(t), "--core", "2.555.1", "--home", home)
	r.wantCode(t, 0)
	contains(t, "stderr", r.stderr, "note: the probe record is from Jenkins 2.541.3, not 2.555.1; ignoring it\n")
	empty := t.TempDir()
	r = cli(withURL(srv.URL), "verify", "--before", before(t), "--home", empty)
	r.wantCode(t, 0)
	contains(t, "stderr", r.stderr, "note: no probe record in "+empty+" (", "); install it with: jenkins-upgrade-check probe --quiet > "+filepath.Join(empty, "init.groovy.d", "zz-jenkins-upgrade-check.groovy"))
}

// An agent that was online is not, and its log says why: roll back.
func TestVerifyRollBack(t *testing.T) {
	srv := controller{core: "2.555.1", jobs: `[{"fullName":"app"}]`,
		plugins: `[{"shortName":"git","version":"5.8","enabled":true,"active":true}]`,
		agents:  `,{"_class":"hudson.slaves.SlaveComputer","displayName":"linux1","offline":true}`}.server(t)
	r := cli(withURL(srv.URL), "verify", "--before", before(t))
	r.wantCode(t, 1)
	if r.stderr != "ROLL BACK\n" {
		t.Errorf("stderr = %q", r.stderr)
	}
	contains(t, "report", r.stdout, "## ROLL BACK", "linux1: runs Java 17; Jenkins 2.555.1 needs Java 21 or newer on agents (from its connection log)")
}

// A controller that does not answer did not start: a verdict, with what
// the log says, not a usage error.
func TestVerifyDown(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close()
	logFile := write(t, filepath.Join(t.TempDir(), "start.log"),
		"2026-10-19 17:01:13.001+0000 [id=30]\tSEVERE\thudson.WebAppMain$3#run: Failed to initialize Jenkins\njava.lang.IllegalStateException: boom\n")
	r := cli(withURL(url), "verify", "--before", before(t), "--core", "2.555.1", "--log", logFile)
	r.wantCode(t, 1)
	contains(t, "report", r.stdout, "# Upgrade check: Jenkins 2.555.1", "Jenkins did not come up", "could not reach the controller", "- hudson.WebAppMain$3#run: Failed to initialize Jenkins — java.lang.IllegalStateException: boom")
}

func TestVerifyErrors(t *testing.T) {
	srv := controller{core: "2.555.1", jobs: `[]`, plugins: `[]`}.server(t)
	b := before(t)
	dir := t.TempDir()
	env := withURL(srv.URL)
	for _, c := range []struct {
		env  map[string]string
		args []string
		want string
	}{
		{env, nil, "--before and --url are required"},
		{jenkinsEnv, []string{"--before", b}, "--before and --url are required"},
		{env, []string{"--before", filepath.Join(dir, "none.json")}, "none.json: no such file or directory"},
		{env, []string{"--before", b, "--plan", filepath.Join(dir, "none-plan.json")}, "none-plan.json: no such file or directory"},
		{env, []string{"--before", b, "--log", filepath.Join(dir, "none.log")}, "none.log: no such file or directory"},
		{map[string]string{"JENKINS_URL": srv.URL, "JENKINS_USER": "audit", "JENKINS_TOKEN": "bad"}, []string{"--before", b}, "the controller answered but refused the API token"},
		{env, []string{"--before", b, "--json", filepath.Join(dir, "no", "v.json")}, "no such file or directory"},
		{env, []string{"--before", b, "-o", filepath.Join(dir, "no", "v.md")}, "no such file or directory"},
	} {
		cli(c.env, append([]string{"verify"}, c.args...)...).wantError(t, c.want)
	}
}
