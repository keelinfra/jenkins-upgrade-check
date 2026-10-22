package verify

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/keelinfra/jenkins-upgrade-check/internal/dryrun"
	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
	"github.com/keelinfra/jenkins-upgrade-check/internal/snapshot"
)

func before() *snapshot.Snapshot {
	return &snapshot.Snapshot{Core: "2.426.3",
		Plugins: []snapshot.Plugin{{Name: "git", Version: "5.0", Active: true}, {Name: "workflow-cps", Version: "3000", Active: true}, {Name: "old", Version: "1", Active: false}},
		Items:   []string{"app", "team", "team/deploy"},
		Agents:  []snapshot.Agent{{Name: "linux1"}, {Name: "aix", Offline: true}},
	}
}

func TestRollBack(t *testing.T) {
	after := &snapshot.Snapshot{Core: "2.580.1", QuietingDown: true,
		Plugins: []snapshot.Plugin{{Name: "git", Version: "5.8", Active: true}},
		Items:   []string{"app", "team"},
		Agents:  []snapshot.Agent{{Name: "linux1", Offline: true}, {Name: "aix", Offline: true}},
	}
	log := "2026-10-19 17:01:12.773+0000 [id=43]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed Loading plugin Pipeline: Groovy v3000 (workflow-cps)\n" +
		"java.io.IOException: Failed to load: Pipeline: Groovy (workflow-cps 3000)\n" +
		"[LF]>  - Update required: Script Security Plugin (script-security 1.75) to be updated to 1336.v1 or higher\n"
	r := Run(Input{Before: before(), After: after, Core: "2.580.1", Log: log,
		Plugins: []plan.Pick{{Name: "git", Version: "5.9"}, {Name: "workflow-cps", Version: "4000"}}})
	if r.Verdict != RollBack {
		t.Fatalf("verdict = %s", r.Verdict)
	}
	states := map[string]string{}
	for _, c := range r.Checks {
		states[c.Name] = c.State + " " + c.Detail + " " + strings.Join(c.Items, ",")
	}
	for name, want := range map[string]string{
		"Jenkins version":  "ok",
		"Plugins":          "fail 1 of 2 plugins that should run are not active workflow-cps",
		"Plugin versions":  "check",
		"Jobs and folders": "fail 1 of 3 did not load team/deploy",
		"Agents":           "check",
	} {
		if !strings.HasPrefix(states[name], want) {
			t.Errorf("%s = %q, want %q…", name, states[name], want)
		}
	}
	if len(r.Roots) != 1 || r.Roots[0].Diagnosis.Kind != "old-dependency" {
		t.Errorf("roots = %+v", r.Roots)
	}
	if !strings.Contains(states["Agents"], "linux1") || strings.Contains(states["Agents"], "aix") {
		t.Errorf("only agents online before count: %s", states["Agents"])
	}
}

func TestGo(t *testing.T) {
	after := &snapshot.Snapshot{Core: "2.580.1",
		Plugins: []snapshot.Plugin{{Name: "git", Version: "5.8", Active: true}, {Name: "workflow-cps", Version: "4000", Active: true}},
		Items:   []string{"app", "team", "team/deploy"},
		Agents:  []snapshot.Agent{{Name: "linux1"}},
	}
	r := Run(Input{Before: before(), After: after, Core: "2.580.1"})
	if r.Verdict != Go {
		t.Errorf("verdict = %s: %+v", r.Verdict, r.Checks)
	}
}

func TestWrongCore(t *testing.T) {
	after := *before()
	r := Run(Input{Before: before(), After: &after, Core: "2.580.1"})
	if r.Verdict != RollBack || r.Checks[0].State != Fail {
		t.Errorf("a controller still on the old version must not pass: %+v", r.Checks[0])
	}
}

func TestDown(t *testing.T) {
	log := "2026-10-19 17:01:12.773+0000 [id=43]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed Loading plugin Pipeline: Groovy v3000 (workflow-cps)\n" +
		"java.io.IOException: Failed to load: Pipeline: Groovy (workflow-cps 3000)\n" +
		"[LF]>  - Update required: Script Security Plugin (script-security 1.75) to be updated to 1336.v1 or higher\n" +
		"2026-10-19 17:01:13.001+0000 [id=30]\tSEVERE\thudson.WebAppMain$3#run: Failed to initialize Jenkins\n" +
		"java.lang.NoClassDefFoundError: org/jenkinsci/main/modules/instance_identity/InstanceIdentity\n"
	r := Run(Input{Before: before(), Core: "2.580.1", Log: log, Down: "could not reach the controller: connection refused"})
	if r.Verdict != RollBack || !r.Down || r.Core != "2.580.1" {
		t.Fatalf("verdict = %s, down = %v, core = %s", r.Verdict, r.Down, r.Core)
	}
	if len(r.Checks) != 1 || r.Checks[0].State != Fail {
		t.Fatalf("checks = %+v", r.Checks)
	}
	if len(r.Roots) != 1 || r.Roots[0].Name != "workflow-cps" {
		t.Errorf("roots = %+v", r.Roots)
	}
	if len(r.Errors) != 1 || !strings.Contains(r.Errors[0], "InstanceIdentity") {
		t.Errorf("errors = %q", r.Errors)
	}
}

// An agent on Java 17 cannot connect to a core that needs 21: a reason to
// roll back, not to wait. Its Java is known from before, or from its log.
func TestAgentTooOldJava(t *testing.T) {
	b := before()
	b.Agents = []snapshot.Agent{{Name: "linux1", Java: "17.0.12"}, {Name: "linux2"}, {Name: "linux3"}, {Name: "linux4"}}
	after := &snapshot.Snapshot{Core: "2.580.1",
		Plugins: []snapshot.Plugin{{Name: "git", Version: "5.8", Active: true}, {Name: "workflow-cps", Version: "4000", Active: true}},
		Items:   []string{"app", "team", "team/deploy"},
		Agents:  []snapshot.Agent{{Name: "linux1", Offline: true}, {Name: "linux2", Offline: true}, {Name: "linux3", Offline: true}, {Name: "linux4", Offline: true}},
	}
	logs := map[string]string{
		"linux2": "Caused by: java.lang.UnsupportedClassVersionError: hudson/slaves/SlaveComputer$SlaveVersion has been compiled by a more recent version of the Java Runtime (class file version 65.0), this version of the Java Runtime only recognizes class file versions up to 61.0\n",
		"linux3": "INFO: Connected\njava.io.EOFException: unexpected stream termination\nINFO: retrying\n",
	}
	r := Run(Input{Before: b, After: after, Core: "2.580.1", AgentLogs: logs})
	if r.Verdict != RollBack {
		t.Fatalf("verdict = %s", r.Verdict)
	}
	var fail, warn Check
	for _, c := range r.Checks {
		if c.Name == "Agents" && c.State == Fail {
			fail = c
		}
		if c.Name == "Agents" && c.State == Warn {
			warn = c
		}
	}
	got := strings.Join(fail.Items, "\n")
	for _, want := range []string{"linux1: runs Java 17; Jenkins 2.580.1 needs Java 21 or newer on agents", "linux2: runs Java 17; Jenkins 2.580.1 needs Java 21 or newer on agents (from its connection log)"} {
		if !strings.Contains(got, want) {
			t.Errorf("fail items lack %q:\n%s", want, got)
		}
	}
	if w := strings.Join(warn.Items, "\n"); w != "linux3 (its log ends with: java.io.EOFException: unexpected stream termination)\nlinux4" {
		t.Errorf("warn items = %q", w)
	}
}

func checks(r *Result) map[string]Check {
	m := map[string]Check{}
	for _, c := range r.Checks {
		m[c.Name+" "+c.State] = c
	}
	return m
}

// writeFile writes body to home/rel.
func writeFile(t *testing.T, home, rel, body string) {
	t.Helper()
	p := filepath.Join(home, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func probe(t *testing.T, js string) *dryrun.Probe {
	t.Helper()
	var p dryrun.Probe
	if err := json.Unmarshal([]byte(js), &p); err != nil {
		t.Fatal(err)
	}
	return &p
}

// With the probe's record, settings Jenkins could not read are blamed on
// their plugin: one that ran before and does not now loses them for good.
func TestProbeSettingsLost(t *testing.T) {
	home := t.TempDir()
	writeFile(t, home, "jobs/app/config.xml", `<project><publishers><hudson.plugins.foo.Notifier plugin="foo@1.0"/></publishers></project>`)
	b := &snapshot.Snapshot{Core: "2.462.3", Plugins: []snapshot.Plugin{
		{Name: "git", Version: "5.0", Active: true}, {Name: "foo", Version: "1.0", Active: true}}}
	after := &snapshot.Snapshot{Core: "2.580.1", Plugins: []snapshot.Plugin{
		{Name: "git", Version: "5.8", Active: true}, {Name: "foo", Version: "1.0", Active: false}}}
	pr := probe(t, `{"completed":true,"version":"2.580.1",
		"plugins":[{"name":"git","version":"5.8","active":true,"enabled":true},{"name":"foo","version":"1.0","active":false,"enabled":true}],
		"failed":[{"name":"foo","causes":["java.io.IOException: Failed to load: Foo (foo 1.0)\n - Jenkins (2.999) or higher required"]}],
		"unreadable":[{"item":"app","path":"jobs/app","problem":"com.thoughtworks.xstream.mapper.CannotResolveClassException: hudson.plugins.foo.Notifier"}],
		"monitors":[{"id":"hudson.diagnosis.OldDataMonitor","name":"Manage Old Data"},{"id":"jenkins.security.QueueItemAuthenticatorMonitor","name":"Builds run as SYSTEM"}]}`)
	r := Run(Input{Before: b, After: after, Probe: pr, Home: home})
	if r.Verdict != RollBack {
		t.Fatalf("verdict = %s", r.Verdict)
	}
	c := checks(r)
	if got := c["Jenkins version ok"].Detail; got != "2.580.1 (was 2.462.3)" {
		t.Errorf("version detail = %q", got)
	}
	if got := c["Settings fail"].Detail; got != "1 setting(s) come from plugins that ran before and do not now; Jenkins drops them when the item is next saved" {
		t.Errorf("settings = %+v", r.Checks)
	}
	if got := c["Jobs and folders check"].Detail; !strings.HasPrefix(got, "not compared: the snapshot from before the upgrade has no job list") {
		t.Errorf("jobs = %q", got)
	}
	want := dryrun.ConfigProblem{Item: "app", Path: "jobs/app", Problem: "com.thoughtworks.xstream.mapper.CannotResolveClassException: hudson.plugins.foo.Notifier", Plugin: "foo 1.0", PluginState: dryrun.PluginFailed}
	if len(r.Unreadable) != 1 || !reflect.DeepEqual(r.Unreadable[0], want) {
		t.Errorf("unreadable = %+v", r.Unreadable)
	}
	if !reflect.DeepEqual(r.Monitors, []string{"Builds run as SYSTEM"}) {
		t.Errorf("monitors = %q", r.Monitors)
	}
	// The probe says why foo did not load; the version comes from the
	// controller.
	if len(r.Roots) != 1 || r.Roots[0].Name != "foo" || r.Roots[0].Version != "1.0" || r.Roots[0].Diagnosis.Kind != dryrun.CauseCoreTooOld {
		t.Errorf("roots = %+v", r.Roots)
	}
}

// Settings from a plugin that was not running before either are a warning,
// and a probe with nothing unreadable is a pass.
func TestProbeSettingsWarnAndOK(t *testing.T) {
	b := &snapshot.Snapshot{Core: "2.462.3", Items: []string{"app"}, Plugins: []snapshot.Plugin{{Name: "git", Version: "5.0", Active: true}, {Name: "bar", Version: "1", Active: false}}}
	after := &snapshot.Snapshot{Core: "2.580.1", Items: []string{"app"}, Plugins: []snapshot.Plugin{{Name: "git", Version: "5.8", Active: true}}}
	home := t.TempDir()
	writeFile(t, home, "jobs/app/config.xml", `<project><bar.Thing plugin="bar@1"/></project>`)
	pr := probe(t, `{"completed":true,"unreadable":[{"item":"app","path":"jobs/app","problem":"CannotResolveClassException: bar.Thing"},{"item":"(global configuration)","path":"","problem":"old field"}]}`)
	r := Run(Input{Before: b, After: after, Probe: pr, Home: home, Core: "2.580.1"})
	c := checks(r)
	if got := c["Settings check"].Detail; !strings.HasPrefix(got, "2 setting(s) Jenkins cannot read, from plugins that were not running before the upgrade either") {
		t.Errorf("checks = %+v", r.Checks)
	}
	if r.Verdict != Go {
		t.Errorf("verdict = %s: a warning is no reason to roll back", r.Verdict)
	}
	if r.Unreadable[0].PluginState != dryrun.PluginNotInstalled {
		t.Errorf("bar: %+v", r.Unreadable[0])
	}

	r = Run(Input{Before: b, After: after, Probe: probe(t, `{"completed":true}`), Core: "2.580.1"})
	if c := checks(r); c["Settings ok"].Detail != "Jenkins read all configuration" || c["Jobs and folders ok"].Detail != "all 1 are back" {
		t.Errorf("checks = %+v", r.Checks)
	}
}

// Without a log or a probe, a plugin that is down still gets a line that
// says what to pass to see why.
func TestPluginDownWithoutReasons(t *testing.T) {
	b := &snapshot.Snapshot{Core: "2.462.3", Plugins: []snapshot.Plugin{
		{Name: "a", Version: "1", Active: true}, {Name: "b", Version: "2", Active: true}, {Name: "c", Version: "3", Active: true}, {Name: "d", Version: "4", Active: true}}}
	after := &snapshot.Snapshot{Core: "2.580.1", Plugins: []snapshot.Plugin{{Name: "a", Version: "1.1"}, {Name: "c", Version: "3"}, {Name: "d", Version: "4"}}}
	// The log names c with no reason and d with one; the probe has c's.
	log := "2026-10-19 17:01:12.773+0000 [id=43]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed Loading plugin C v3 (c)\n" +
		"2026-10-19 17:01:12.774+0000 [id=43]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed Loading plugin D v4 (d)\n" +
		"java.io.IOException: Failed to load: D (d 4)\n" +
		" - Plugin is missing: e (1.0)\n"
	pr := probe(t, `{"completed":true,"failed":[{"name":"c","causes":["java.lang.UnsupportedClassVersionError: c/C (class file version 65.0)"]},{"name":"d","causes":["java.lang.LinkageError: other"]}]}`)
	r := Run(Input{Before: b, After: after, Log: log, Probe: pr})
	why := map[string]dryrun.Root{}
	for _, rt := range r.Roots {
		why[rt.Name] = rt
	}
	if got := why["a"]; got.Version != "1.1" || got.Diagnosis.Why != "installed but not active; pass the Jenkins log with --log, or --home after a start with the probe script, to see why" {
		t.Errorf("a = %+v", got)
	}
	if got := why["b"]; got.Version != "" || !strings.HasPrefix(got.Diagnosis.Why, "not installed any more, or failed to load") {
		t.Errorf("b = %+v", got)
	}
	if got := why["c"]; got.Diagnosis.Kind != dryrun.CauseNewerJava || got.Version != "3" {
		t.Errorf("c: the probe's reason fills in for a log without one: %+v", got)
	}
	if got := why["d"]; got.Diagnosis.Kind != dryrun.CauseMissingDependency {
		t.Errorf("d: the log's reason stays: %+v", got)
	}
}

// A job that did not load is blamed on the plugin its config.xml names,
// with the reason from the log.
func TestMissingItemBlamed(t *testing.T) {
	home := t.TempDir()
	writeFile(t, home, "jobs/team/jobs/deploy/config.xml", `<?xml version='1.1' encoding='UTF-8'?>
<flow-definition plugin="workflow-job@1385.v1"><definition class="org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition" plugin="workflow-cps@3000"/></flow-definition>`)
	after := &snapshot.Snapshot{Core: "2.580.1",
		Plugins: []snapshot.Plugin{{Name: "git", Version: "5.8", Active: true}},
		Items:   []string{"app", "team"},
		Agents:  []snapshot.Agent{{Name: "linux1"}}}
	log := "2026-10-19 17:01:14.257+0000 [id=46]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed Loading item deploy\n" +
		"com.thoughtworks.xstream.mapper.CannotResolveClassException: org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition\n"
	r := Run(Input{Before: before(), After: after, Core: "2.580.1", Log: log, Home: home})
	want := dryrun.ConfigProblem{Item: "team/deploy", Path: "jobs/team/jobs/deploy", Problem: "CannotResolveClassException: org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition", Plugin: "workflow-cps 3000", PluginState: dryrun.PluginFailed}
	if len(r.Missing) != 1 || !reflect.DeepEqual(r.Missing[0], want) {
		t.Errorf("missing = %+v\nwant %+v", r.Missing, want)
	}
	if c := checks(r); c["Agents ok"].Detail != "all 1 agents that were online are online" {
		t.Errorf("checks = %+v", r.Checks)
	}
}

// A controller that does not answer: without a log, the report asks for
// one; with JENKINS_HOME, an error that names a class is put on the plugin
// that ships it; and only the first ten other errors are kept.
func TestDownWithHome(t *testing.T) {
	r := Run(Input{Before: before(), Core: "2.580.1", Down: "connection refused"})
	if len(r.Checks) != 1 || !strings.HasSuffix(r.Checks[0].Detail, "; pass the Jenkins log with --log to see why") || r.Roots != nil {
		t.Errorf("no log: %+v", r)
	}

	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, "plugins"), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(filepath.Join(home, "plugins", "example.jpi"))
	if err != nil {
		t.Fatal(err)
	}
	z := zip.NewWriter(f)
	if _, err := z.Create("WEB-INF/classes/org/example/Foo.class"); err != nil {
		t.Fatal(err)
	}
	z.Close()
	f.Close()

	// example ships org.example.Foo; workflow-cps failed on its own, so
	// an error in one of its classes would say nothing new.
	log := "2024-10-30 12:00:01.000+0000 [id=35]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed Loading plugin Pipeline: Groovy v2.90 (workflow-cps)\n" +
		"java.io.IOException: Failed to load: Pipeline: Groovy (workflow-cps 2.90)\n" +
		" - Jenkins (2.999) or higher required\n" +
		"2024-10-30 12:00:05.000+0000 [id=40]\tSEVERE\thudson.ExtensionFinder$GuiceFinder$SezpozModule#configure: Failed to load org.example.Foo\n" +
		"java.lang.NoClassDefFoundError: javax/servlet/ServletException\n"
	for i := 0; i < 12; i++ {
		log += fmt.Sprintf("2024-10-30 12:00:06.%03d+0000 [id=40]\tSEVERE\tsome.Thing#run: failure %d\njava.lang.IllegalStateException: boom %d\n", i, i, i)
	}
	r = Run(Input{Before: before(), Core: "2.580.1", Down: "connection refused", Log: log, Home: home})
	if len(r.Roots) != 2 || r.Roots[0].Name != "workflow-cps" || r.Roots[1].Name != "example" || r.Roots[1].Diagnosis.Kind != dryrun.CauseJakartaEE9 {
		t.Errorf("roots = %+v", r.Roots)
	}
	if len(r.Errors) != 10 || r.Errors[0] != "some.Thing#run: failure 0 — java.lang.IllegalStateException: boom 0" {
		t.Errorf("errors = %q", r.Errors)
	}
}

// A long line at the end of an agent's log is cut.
func TestAgentCauseLongLine(t *testing.T) {
	cause, final := agentCause("", "SEVERE: "+strings.Repeat("x", 300), "2.580.1")
	if final || cause != "its log ends with: SEVERE: "+strings.Repeat("x", 192)+"…" {
		t.Errorf("cause = %q, %v", cause, final)
	}
	if cause, final := agentCause("", "INFO: all good\n", "2.580.1"); cause != "" || final {
		t.Errorf("a log without an error: %q, %v", cause, final)
	}
}
