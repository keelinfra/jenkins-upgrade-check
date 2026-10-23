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

func TestResolveRemovedDependency(t *testing.T) {
	rel, _ := fixtures(t)
	// By name only, as resolve --names does: the dependent says which
	// version, and the download archive has it.
	r := Resolve(rel, []Pick{{Name: "momentjs", Version: "0"}, {Name: "stage-view", Version: "0"}, {Name: "gone", Version: "0"}}, "2.346.3")
	got := r.ByName()
	m := got["momentjs"]
	if m.Version != "1.1.1" || !m.Unknown || !strings.HasSuffix(m.URL, "/momentjs/1.1.1/momentjs.hpi") {
		t.Errorf("momentjs = %+v, want 1.1.1 from the archive", m)
	}
	if _, ok := got["gone"]; ok {
		t.Error("a removed plugin with no version and no dependent must be left out")
	}
	if len(r.Problems) != 2 {
		t.Errorf("want problems about momentjs and gone, got %v", r.Problems)
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

func TestStopPicksEachStepAndTheTarget(t *testing.T) {
	rel, center := fixtures(t)
	s := &snapshot.Snapshot{Core: "2.346.3", Plugins: []snapshot.Plugin{{Name: "workflow-cps", Version: "2.90", Enabled: true}}}
	p, err := Build(s, "2.580.1", rel, center)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Steps) != 3 {
		t.Fatalf("2.346.3 to 2.580.1 should stop at 2.541.3: %d steps", len(p.Steps))
	}
	core, picks, java, err := p.Stop(1)
	if err != nil || core != "2.541.3" || java != 17 || len(picks) == 0 {
		t.Errorf("step 1 = %s on Java %d with %d plugins (%v), want 2.541.3 on Java 17", core, java, len(picks), err)
	}
	core, _, java, err = p.Stop(-1)
	if err != nil || core != "2.580.1" || java != 21 {
		t.Errorf("target = %s on Java %d (%v), want 2.580.1 on Java 21", core, java, err)
	}
	if _, _, _, err := p.Stop(3); err == nil {
		t.Error("a step past the end should be an error")
	}
}

// The newest trigger release that 2.60.3 can load needs a credentials release
// that 2.60.3 cannot, so the older trigger is the one that works; a removed
// dependency does not count against a release.
func TestResolvePassesOverReleasesWithDependenciesTooNew(t *testing.T) {
	rel, err := uc.ReadReleases(strings.NewReader(`{"plugins":{
 "trigger":{
  "2.0":{"requiredCore":"2.60.3","dependencies":[{"name":"credentials","version":"2.6.1.1","optional":false}]},
  "1.0":{"requiredCore":"2.60.3","dependencies":[{"name":"credentials","version":"2.1.0","optional":false},{"name":"gone","version":"1.0","optional":false}]}
 },
 "credentials":{
  "2.6.1.1":{"requiredCore":"2.138.4"},
  "2.1.19":{"requiredCore":"2.60.3"}
 },
 "report":{
  "3.0":{"requiredCore":"2.60.3","dependencies":[{"name":"token-macro","version":"2.8","optional":true}]}
 },
 "token-macro":{"2.6":{"requiredCore":"2.60.3"}}
}}`))
	if err != nil {
		t.Fatal(err)
	}
	r := Resolve(rel, []Pick{{Name: "trigger", Version: "0"}, {Name: "report", Version: "0"}}, "2.60.3")
	got := r.ByName()
	if got["trigger"].Version != "1.0" || got["credentials"].Version != "2.1.19" {
		t.Errorf("want trigger 1.0 with credentials 2.1.19, got %+v", r.Plugins)
	}
	if got["report"].Version != "3.0" {
		t.Errorf("an optional dependency that is not installed does not hold a release back, got %+v", got["report"])
	}
	if len(r.Problems) != 1 || !strings.Contains(r.Problems[0], "gone") {
		t.Errorf("want only the removed dependency as a problem, got %v", r.Problems)
	}

	// On a core that runs the newer credentials, the newer trigger wins.
	r = Resolve(rel, []Pick{{Name: "trigger", Version: "0"}}, "2.138.4")
	if got := r.ByName(); got["trigger"].Version != "2.0" || got["credentials"].Version != "2.6.1.1" {
		t.Errorf("on 2.138.4 want trigger 2.0 with credentials 2.6.1.1, got %+v", r.Plugins)
	}
}

func TestPrevious(t *testing.T) {
	p := &Plan{Steps: []Step{{Core: "2.249.3"}, {Core: "2.462.3"}, {Core: "2.541.3"}, {Core: "2.580.1"}}, Target: Resolution{Core: "2.580.1"}}
	for step, want := range map[int]string{0: "", 1: "2.249.3", 2: "2.462.3", 3: "2.541.3", -1: "2.541.3", 4: ""} {
		if got := p.Previous(step); got != want {
			t.Errorf("Previous(%d) = %q, want %q", step, got, want)
		}
	}
	if got := (&Plan{Steps: []Step{{Core: "2.541.3"}}}).Previous(-1); got != "" {
		t.Errorf("no core upgrade: %q", got)
	}
}

// A plugin built for an older core expects what core split out since; the
// war no longer carries those plugins, so the plan adds them.
func TestResolveAddsPluginsSplitOutOfCore(t *testing.T) {
	rel, err := uc.ReadReleases(strings.NewReader(`{"plugins":{
 "translation":{"1.16":{"requiredCore":"1.466"}},
 "script-security":{"1.0":{"requiredCore":"2.60.3"}},
 "modern":{"1.0":{"requiredCore":"2.479.3"}},
 "javax-mail-api":{"1.6.2-11":{"requiredCore":"2.361.4","dependencies":[{"name":"javax-activation-api","version":"1.2.0-2","optional":false}]}},
 "javax-activation-api":{"1.2.0-8":{"requiredCore":"2.361.4"}},
 "trilead-api":{"2.1.0":{"requiredCore":"2.184"}},
 "command-launcher":{"1.0":{"requiredCore":"2.86"}}
}}`))
	if err != nil {
		t.Fatal(err)
	}
	r := Resolve(rel, []Pick{{Name: "translation", Version: "1.16"}, {Name: "script-security", Version: "1.0"}, {Name: "modern", Version: "1.0"}}, "2.580.1")
	got := r.ByName()
	for _, name := range []string{"javax-mail-api", "javax-activation-api", "trilead-api", "command-launcher"} {
		if p, ok := got[name]; !ok || !p.Added || !contains(p.RequiredBy, "translation") {
			t.Errorf("%s: want added for translation, got %+v", name, p)
		}
	}
	if contains(got["command-launcher"].RequiredBy, "script-security") {
		t.Errorf("script-security is exempt from command-launcher (split-plugin-cycles.txt): %+v", got["command-launcher"])
	}
	if len(r.Problems) != 0 {
		t.Errorf("problems: %v", r.Problems)
	}

	// Before the split, core still has it: nothing to add.
	r = Resolve(rel, []Pick{{Name: "translation", Version: "1.16"}}, "2.249.3")
	got = r.ByName()
	if _, ok := got["javax-mail-api"]; ok {
		t.Errorf("2.249.3 still ships JavaMail; got %+v", r.Plugins)
	}
	if _, ok := got["trilead-api"]; !ok {
		t.Errorf("trilead-api was split out in 2.184; got %+v", r.Plugins)
	}
}

// Dependencies that are installed but too old, and ones no release on the
// core satisfies: each is raised where a release exists, and named as a
// problem where none does.
func TestResolveRaisesInstalledDependencies(t *testing.T) {
	rel, err := uc.ReadReleases(strings.NewReader(`{"plugins":{
 "web":{"1.0":{"requiredCore":"2.400","dependencies":[{"name":"util","version":"2.0","optional":false}]}},
 "util":{
  "1.0":{"requiredCore":"2.400","url":"https://u/util-1.0.hpi"},
  "2.0":{"requiredCore":"2.400","url":"https://u/util-2.0.hpi","sha256":"s2","dependencies":[{"name":"core-lib","version":"5.0","optional":false}]}
 },
 "core-lib":{"5.0":{"requiredCore":"2.500"}},
 "needy":{"1.0":{"requiredCore":"2.400","dependencies":[{"name":"old","version":"3.0","optional":false}]}},
 "old":{"1.0":{"requiredCore":"2.300"},"3.0":{"requiredCore":"2.500"}},
 "uses-private":{"1.0":{"requiredCore":"2.400","dependencies":[{"name":"private","version":"2.0","optional":false}]}},
 "app":{"1.0":{"requiredCore":"2.400","dependencies":[{"name":"lib","version":"2.0","optional":false}]}},
 "lib":{"1.0":{"requiredCore":"2.300"},"3.0":{"requiredCore":"2.500"}},
 "too-new":{"9.0":{"requiredCore":"2.500"}}
}}`))
	if err != nil {
		t.Fatal(err)
	}
	r := Resolve(rel, []Pick{
		{Name: "web", Version: "1.0"}, {Name: "util", Version: "1.0"},
		{Name: "needy", Version: "1.0"}, {Name: "old", Version: "1.0"},
		{Name: "uses-private", Version: "1.0"}, {Name: "private", Version: "1.0"},
		{Name: "app", Version: "1.0"},
		{Name: "too-new", Version: "8.0"},
	}, "2.400")
	got := r.ByName()
	// web needs util 2.0, which loads on 2.400 even though its own
	// dependency does not: util is raised, and the dependency is a problem.
	if u := got["util"]; u.Version != "2.0" || u.URL != "https://u/util-2.0.hpi" || u.SHA256 != "s2" {
		t.Errorf("util = %+v, want 2.0 with its URL and checksum", u)
	}
	// A private plugin is raised to the version its dependent asks for,
	// from the download archive.
	if p := got["private"]; p.Version != "2.0" || !p.Unknown || p.URL != uc.ArchiveURL("private", "2.0") {
		t.Errorf("private = %+v, want 2.0 from the archive", p)
	}
	if o := got["old"]; o.Version != "1.0" {
		t.Errorf("old = %+v, want kept at 1.0", o)
	}
	if n := got["too-new"]; n.Version != "8.0" {
		t.Errorf("too-new = %+v, want kept at 8.0", n)
	}
	if _, ok := got["lib"]; ok {
		t.Errorf("lib has no release that satisfies app; it must not be added: %+v", got["lib"])
	}
	want := []string{
		"app 1.0 needs lib 2.0 or newer, and no such release runs on Jenkins 2.400",
		"needy 1.0 needs old 3.0 or newer, but the newest old that runs on Jenkins 2.400 is 1.0",
		"too-new: no release runs on Jenkins 2.400; kept at 8.0",
		"util 2.0 needs core-lib 5.0 or newer, and no such release runs on Jenkins 2.400",
	}
	if !reflect.DeepEqual(r.Problems, want) {
		t.Errorf("problems =\n%q\nwant\n%q", r.Problems, want)
	}
}

// Two plugins that depend on each other resolve instead of looping.
func TestResolveDependencyCycle(t *testing.T) {
	rel, err := uc.ReadReleases(strings.NewReader(`{"plugins":{
 "a":{"1.0":{"requiredCore":"2.400","dependencies":[{"name":"b","version":"1.0","optional":false}]}},
 "b":{"1.0":{"requiredCore":"2.400","dependencies":[{"name":"a","version":"1.0","optional":false}]}}
}}`))
	if err != nil {
		t.Fatal(err)
	}
	r := Resolve(rel, []Pick{{Name: "a", Version: "0"}}, "2.400")
	got := r.ByName()
	if got["a"].Version != "1.0" || got["b"].Version != "1.0" || !got["b"].Added || len(r.Problems) != 0 {
		t.Errorf("plugins = %+v, problems = %q", r.Plugins, r.Problems)
	}
}

// A plugin built for an old core expects a split plugin; when no release
// of it runs on the core, that is a problem.
func TestResolveSplitPluginTooNew(t *testing.T) {
	rel, err := uc.ReadReleases(strings.NewReader(`{"plugins":{
 "ancient":{"1.0":{"requiredCore":"1.400"}},
 "mailer":{"500.v1":{"requiredCore":"2.500"}}
}}`))
	if err != nil {
		t.Fatal(err)
	}
	r := Resolve(rel, []Pick{{Name: "ancient", Version: "1.0"}}, "2.400")
	want := "ancient 1.0 was built for Jenkins 1.400 and expects mailer, which was split out of core, but no mailer release runs on Jenkins 2.400"
	if len(r.Problems) != 1 || r.Problems[0] != want {
		t.Errorf("problems = %q, want %q", r.Problems, want)
	}
	if _, ok := r.ByName()["mailer"]; ok {
		t.Error("mailer cannot load on 2.400; it must not be added")
	}
}

// Locate keeps every version and only says where to download it.
func TestLocate(t *testing.T) {
	rel, _ := fixtures(t)
	got := Locate(rel, []Pick{{Name: "workflow-cps", Version: "2.90"}, {Name: "private-plugin", Version: "1.0"}, {Name: "workflow-cps", Version: "1.0"}})
	if p := got[0]; p.Version != "2.90" || p.URL != "https://updates.jenkins.io/download/plugins/workflow-cps/2.90/workflow-cps.hpi" || p.SHA256 != "x" || p.Unknown {
		t.Errorf("published release = %+v", p)
	}
	if p := got[1]; p.Version != "1.0" || !p.Unknown || p.URL != uc.ArchiveURL("private-plugin", "1.0") {
		t.Errorf("private plugin = %+v, want the archive URL", p)
	}
	if p := got[2]; p.Version != "1.0" || p.URL != "" || p.Unknown {
		t.Errorf("a release the update center does not list for a known plugin = %+v, want no URL", p)
	}
}

func TestAttention(t *testing.T) {
	for _, c := range []struct {
		p    PluginPlan
		want bool
	}{
		{PluginPlan{Name: "git", Current: "4.0", Target: "5.0"}, false},
		{PluginPlan{Unknown: true}, true},
		{PluginPlan{NotDistributed: true}, true},
		{PluginPlan{Deprecated: "https://x"}, true},
		{PluginPlan{VulnerableNow: []WarningRef{{ID: "S"}}}, true},
		{PluginPlan{VulnerableAfter: []WarningRef{{ID: "S"}}}, true},
	} {
		if got := c.p.Attention(); got != c.want {
			t.Errorf("Attention(%+v) = %v, want %v", c.p, got, c.want)
		}
	}
}

func TestBuildNeedsACore(t *testing.T) {
	rel, center := fixtures(t)
	_, err := Build(&snapshot.Snapshot{}, "2.555.1", rel, center)
	if err == nil || err.Error() != "the snapshot does not say which Jenkins version it is; pass --core" {
		t.Errorf("err = %v", err)
	}
	_, err = Build(&snapshot.Snapshot{Core: "2.555.1"}, "2.541.3", rel, center)
	if err == nil || err.Error() != "the controller already runs 2.555.1, which is not older than 2.541.3" {
		t.Errorf("err = %v", err)
	}
}

// The controller's own Java, a deprecation without a link, and agents on
// HP-UX and old AIX.
func TestBuildNotesAndAgents(t *testing.T) {
	rel, _ := fixtures(t)
	center, err := uc.ReadUpdateCenter(strings.NewReader(`{"plugins":{"workflow-cps":{}},"deprecations":{"workflow-cps":{}}}`))
	if err != nil {
		t.Fatal(err)
	}
	s := &snapshot.Snapshot{Core: "2.479.3", Java: "17.0.9", Notes: []string{"agents could not be read"},
		Plugins: []snapshot.Plugin{{Name: "workflow-cps", Version: "2.90", Enabled: true}},
		Agents: []snapshot.Agent{
			{Name: "hp", OS: "HP-UX (ia64)", Java: "1.8.0_202"},
			{Name: "aix71", OS: "AIX (ppc64)", Java: "1.8.0_391"},
			{Name: "aix-unknown", OS: "AIX (ppc64)"},
			{Name: "linux", OS: "Linux (amd64)"},
		}}
	p, err := Build(s, "2.555.1", rel, center)
	if err != nil {
		t.Fatal(err)
	}
	wantNotes := []string{"agents could not be read", "the controller runs Java 17.0.9; Jenkins 2.555.1 runs on Java 21 or 25"}
	if !reflect.DeepEqual(p.Notes, wantNotes) {
		t.Errorf("notes = %q, want %q", p.Notes, wantNotes)
	}
	for _, pp := range p.Plugins {
		if pp.Name == "workflow-cps" && pp.Deprecated != "deprecated" {
			t.Errorf("a deprecation without a URL still marks the plugin: %+v", pp)
		}
	}
	problems := map[string]string{}
	for _, a := range p.Agents {
		problems[a.Name] = a.Problem
	}
	if !strings.HasPrefix(problems["hp"], "HP-UX has no Java 21 (HP's Java stops at 8)") {
		t.Errorf("hp: %q", problems["hp"])
	}
	for _, n := range []string{"aix71", "aix-unknown"} {
		if !strings.HasPrefix(problems[n], "Jenkins 2.555.1 needs Java 21 on this agent; IBM Semeru 21 needs AIX 7.2 TL5") {
			t.Errorf("%s: %q", n, problems[n])
		}
	}
	if _, ok := problems["linux"]; ok || len(p.Agents) != 3 {
		t.Errorf("an agent with no known Java on Linux is not flagged: %+v", p.Agents)
	}
	if p.Summary.Deprecated != 1 || p.Summary.Agents != 3 {
		t.Errorf("summary = %+v", p.Summary)
	}
}

func TestStopWithoutPlugins(t *testing.T) {
	p := &Plan{Steps: []Step{{Core: "2.462.3"}, {Core: "2.541.3", Upgrade: true}}}
	_, _, _, err := p.Stop(1)
	if err == nil || err.Error() != "step 1 has no plugin set; make the plan again with this version" {
		t.Errorf("err = %v", err)
	}
	core, picks, java, err := (&Plan{Target: Resolution{Core: "2.555.1"}}).Stop(-1)
	if err != nil || core != "2.555.1" || picks != nil || java != 0 {
		t.Errorf("a plan without steps: %s %v %d %v", core, picks, java, err)
	}
}
