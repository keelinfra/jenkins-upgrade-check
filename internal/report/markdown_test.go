package report

import (
	"strings"
	"testing"

	"github.com/keelinfra/jenkins-upgrade-check/internal/dryrun"
	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
)

func TestPlanReport(t *testing.T) {
	p := &plan.Plan{
		From: "2.462.3", To: "2.555.1", JavaFrom: []int{11, 17}, JavaTo: []int{21, 25},
		Steps: []plan.Step{{Core: "2.462.3"}, {Core: "2.541.3", Java: 17, Upgrade: true}, {Core: "2.555.1", Java: 21, Upgrade: true}},
		Plugins: []plan.PluginPlan{
			{Name: "script-security", Current: "1.75", Target: "1350.v1", VulnerableNow: []plan.WarningRef{{ID: "SECURITY-1", URL: "https://x"}}},
			{Name: "workflow-api", Target: "1300.v1", Added: true, RequiredBy: []string{"workflow-cps"}},
		},
		Agents:  []plan.AgentNote{{Name: "sol", OS: "SunOS (sparc)", Problem: "Solaris has no Java 21"}},
		Summary: plan.Summary{Plugins: 2, Upgraded: 1, Added: 1, VulnerableNow: 1, CoreUpgrades: 2, Agents: 1},
	}
	var b strings.Builder
	if err := Plan(&b, p); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	for _, want := range []string{"2.462.3 → 2.555.1", "Run the controller and every agent on Java 17", "SECURITY-1", "fixed by the upgrade", "sol", "new, for workflow-cps", "keelinfra.io/jenkins/pricing"} {
		if !strings.Contains(out, want) {
			t.Errorf("report lacks %q\n%s", want, out)
		}
	}
}

func TestRehearsalReport(t *testing.T) {
	var b strings.Builder
	Rehearsal(&b, &dryrun.Result{Core: "2.555.1", Image: "jenkins/jenkins:2.555.1-jdk21", Plugins: 3, Started: true, StartupSeconds: 40,
		Failed: []dryrun.Failure{{Name: "workflow-cps", Version: "2.90", Reasons: []string{"Update required"}}}})
	if !strings.Contains(b.String(), "1 plugin(s) failed to load") || !strings.Contains(b.String(), "workflow-cps") {
		t.Error(b.String())
	}
}
