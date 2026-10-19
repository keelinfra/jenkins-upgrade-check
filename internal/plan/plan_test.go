package plan

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/keelinfra/jenkins-upgrade-check/internal/snapshot"
	"github.com/keelinfra/jenkins-upgrade-check/internal/uc"
)

func fixtures(t *testing.T) (*uc.Releases, *uc.UpdateCenter) {
	t.Helper()
	f, err := os.Open("testdata/plugin-versions.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rel, err := uc.ReadReleases(f)
	if err != nil {
		t.Fatal(err)
	}
	g, err := os.Open("testdata/update-center.actual.json")
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	center, err := uc.ReadUpdateCenter(g)
	if err != nil {
		t.Fatal(err)
	}
	return rel, center
}

func TestResolveRaisesAndAddsDependencies(t *testing.T) {
	rel, _ := fixtures(t)
	r := Resolve(rel, []Pick{{Name: "workflow-cps", Version: "2.90"}, {Name: "script-security", Version: "1.75"}}, "2.479.3")
	got := r.ByName()
	if got["workflow-cps"].Version != "4000.v1" {
		t.Errorf("workflow-cps = %s, want 4000.v1 (the beta is skipped)", got["workflow-cps"].Version)
	}
	if got["script-security"].Version != "1350.v1" {
		t.Errorf("script-security = %s, want 1350.v1", got["script-security"].Version)
	}
	api := got["workflow-api"]
	if !api.Added || api.Version != "1300.v1" || !reflect.DeepEqual(api.RequiredBy, []string{"workflow-cps"}) {
		t.Errorf("workflow-api should be added for workflow-cps, got %+v", api)
	}
	if _, ok := got["ionicons-api"]; ok {
		t.Error("an optional dependency that is not installed must not be added")
	}
	if len(r.Problems) != 0 {
		t.Errorf("unexpected problems: %v", r.Problems)
	}
}

func TestResolveStaysWithinCore(t *testing.T) {
	rel, _ := fixtures(t)
	r := Resolve(rel, []Pick{{Name: "workflow-cps", Version: "2.90"}}, "2.462.3")
	got := r.ByName()
	if got["workflow-cps"].Version != "3800.v1" || got["script-security"].Version != "1300.v1" {
		t.Errorf("on 2.462.3 want workflow-cps 3800.v1 + script-security 1300.v1, got %+v", r.Plugins)
	}
}

func TestResolveUnknownAndImpossible(t *testing.T) {
	rel, _ := fixtures(t)
	r := Resolve(rel, []Pick{{Name: "private-plugin", Version: "1.0"}, {Name: "needs-future", Version: "0"}}, "2.479.3")
	got := r.ByName()
	if !got["private-plugin"].Unknown || got["private-plugin"].Version != "1.0" {
		t.Errorf("a plugin the update center never published is kept as it is, got %+v", got["private-plugin"])
	}
	if len(r.Problems) != 1 || !strings.Contains(r.Problems[0], "needs-future") {
		t.Errorf("want one problem about needs-future, got %v", r.Problems)
	}
}

func TestBuild(t *testing.T) {
	rel, center := fixtures(t)
	s := &snapshot.Snapshot{
		Source: "test",
		Core:   "2.346.3",
		Plugins: []snapshot.Plugin{
			{Name: "workflow-cps", Version: "2.90", Enabled: true},
			{Name: "script-security", Version: "1.75", Enabled: true},
			{Name: "old-thing", Version: "1.0", Enabled: false},
			{Name: "private-plugin", Version: "1.0", Enabled: true},
		},
		Agents: []snapshot.Agent{
			{Name: "sol", OS: "SunOS (sparc)"},
			{Name: "lin", OS: "Linux (amd64)", Java: "17.0.9"},
			{Name: "lin21", OS: "Linux (amd64)", Java: "21.0.5"},
			{Name: "aix73", OS: "AIX (ppc64)", Java: "21.0.4"},
		},
	}
	p, err := Build(s, "2.555.1", rel, center)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Steps) != 3 || p.Steps[0].Upgrade || p.Steps[1].Core != "2.541.3" || p.Steps[1].Java != 17 || p.Steps[2].Core != "2.555.1" || p.Steps[2].Java != 21 {
		t.Fatalf("steps = %+v", p.Steps)
	}
	by := map[string]PluginPlan{}
	for _, pp := range p.Plugins {
		by[pp.Name] = pp
	}
	ss := by["script-security"]
	if len(ss.VulnerableNow) != 1 || len(ss.VulnerableAfter) != 0 {
		t.Errorf("script-security 1.75 is vulnerable and the target is not: %+v", ss)
	}
	old := by["old-thing"]
	if old.Deprecated == "" || !old.NotDistributed || len(old.VulnerableAfter) != 1 || !old.Disabled {
		t.Errorf("old-thing should be deprecated, not distributed, still vulnerable, disabled: %+v", old)
	}
	if !by["private-plugin"].Unknown {
		t.Error("private-plugin should be unknown")
	}
	if !by["workflow-api"].Added {
		t.Error("workflow-api should be a new dependency")
	}
	if len(p.CoreWarnings) != 1 {
		t.Errorf("2.346.3 is affected by the core warning: %+v", p.CoreWarnings)
	}
	agents := map[string]bool{}
	for _, a := range p.Agents {
		agents[a.Name] = true
	}
	if !agents["sol"] || !agents["lin"] || agents["lin21"] || agents["aix73"] {
		t.Errorf("agents at risk = %v, want sol and lin", agents)
	}
	if p.Summary.CoreUpgrades != 2 || p.Summary.Added != 1 || p.Summary.Unknown != 1 {
		t.Errorf("summary = %+v", p.Summary)
	}
}

func TestBuildRefusesDowngrade(t *testing.T) {
	rel, center := fixtures(t)
	_, err := Build(&snapshot.Snapshot{Core: "2.555.1"}, "2.541.3", rel, center)
	if err == nil {
		t.Fatal("want an error for a target that is not newer")
	}
}

func TestJavaMajor(t *testing.T) {
	for in, want := range map[string]int{"1.8.0_392": 8, "17.0.9": 17, "21": 21, "11.0.2+9": 11, "": 0, "n/a": 0} {
		if got := JavaMajor(in); got != want {
			t.Errorf("JavaMajor(%q) = %d, want %d", in, got, want)
		}
	}
}
