package verify

import (
	"strings"
	"testing"

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
