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
