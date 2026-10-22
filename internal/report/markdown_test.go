package report

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/keelinfra/jenkins-upgrade-check/internal/dryrun"
	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
	"github.com/keelinfra/jenkins-upgrade-check/internal/verify"
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
		Failed: []dryrun.Failure{
			{Name: "script-security", Version: "1.75", Reasons: []string{"Jenkins (2.400) or higher required"}},
			{Name: "workflow-cps", Version: "2.90", Reasons: []string{"Failed to load: Script Security (script-security 1.75)"}},
		},
		Home: &dryrun.HomeResult{Source: "/var/lib/jenkins", Files: 12, Items: 3, Loaded: 2,
			Missing: []dryrun.ConfigProblem{{Item: "deploy", Problem: "CannotResolveClassException: flow-definition", Plugin: "workflow-job 2.40", PluginState: dryrun.PluginFailed}}},
	})
	out := b.String()
	for _, want := range []string{"2 plugin(s) did not load, because of 1 cause(s)", "| script-security 1.75 | this release needs Jenkins 2.400 or newer |", "| workflow-cps |", "2 of 3 jobs and folders loaded", "| deploy | plugin workflow-job 2.40, failed to load |"} {
		if !strings.Contains(out, want) {
			t.Errorf("report lacks %q\n%s", want, out)
		}
	}
}

func contains(t *testing.T, out string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(out, want) {
			t.Errorf("report lacks %q\n%s", want, out)
		}
	}
}

func lacks(t *testing.T, out string, unwanted ...string) {
	t.Helper()
	for _, u := range unwanted {
		if strings.Contains(out, u) {
			t.Errorf("report has %q\n%s", u, out)
		}
	}
}

// Every flag of the summary, the controller's Java, core warnings, notes,
// offline agents and each reason a plugin needs a decision.
func TestPlanReportAllSections(t *testing.T) {
	p := &plan.Plan{
		Source: "https://ci.example.com", Tool: "jenkins-upgrade-check", ToolVersion: "1.0",
		GeneratedAt: time.Date(2026, 10, 1, 12, 30, 0, 0, time.UTC),
		From:        "2.479.3", To: "2.492.3", JavaFrom: []int{17, 21}, JavaTo: []int{17, 21}, ControllerJava: "17.0.9",
		Steps: []plan.Step{
			{Core: "2.479.3", Changes: []plan.Change{{Name: "git", From: "5.0", To: "5.2"}}},
			{Core: "2.492.3", Java: 0, Upgrade: true, Problems: []string{"gitlab-plugin 1.5 needs jersey3-api 2.0 or newer"}},
		},
		CoreWarnings: []plan.WarningRef{{ID: "SECURITY-3501", URL: "https://www.jenkins.io/security/advisory/2024-10-02/", Message: "Exposure of secrets"}},
		Plugins: []plan.PluginPlan{
			{Name: "both", Current: "1.0", Target: "1.1",
				VulnerableNow:   []plan.WarningRef{{ID: "S-1", URL: "https://s1"}, {ID: "S-2", URL: "https://s2"}},
				VulnerableAfter: []plan.WarningRef{{ID: "S-2", URL: "https://s2"}, {ID: "S-3", URL: "https://s3"}}},
			{Name: "dep-link", Current: "1", Target: "1", Deprecated: "https://plugins.jenkins.io/dep-link/"},
			{Name: "dep-plain", Current: "1", Target: "1", Deprecated: "deprecated", NotDistributed: true},
			{Name: "private", Current: "0.1", Target: "0.1", Unknown: true},
			{Name: "quiet", Current: "2", Target: "3"},
		},
		Agents: []plan.AgentNote{{Name: "aix", Offline: true, Problem: "needs Java 17"}},
		Notes:  []string{"agents could not be read"},
		Summary: plan.Summary{Plugins: 5, Upgraded: 2, VulnerableNow: 1, VulnerableAfter: 1, Deprecated: 2, NotDistributed: 1,
			Unknown: 1, Problems: 1, CoreUpgrades: 1, Agents: 1},
	}
	var b strings.Builder
	if err := Plan(&b, p); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	contains(t, out,
		"Source: https://ci.example.com · generated 2026-10-01 12:30 UTC by jenkins-upgrade-check 1.0",
		"**1 core upgrade(s)**, 5 plugins: 2 to upgrade, 0 new dependencies; 1 with security warnings today, 1 still with warnings after the upgrade, 2 deprecated, 1 no longer distributed, 1 not from the update center, 1 dependency problem(s), 1 agent(s) that cannot follow.",
		"Jenkins 2.479.3 runs on Java 17, 21; Jenkins 2.492.3 runs on Java 17, 21. The controller runs Java 17.0.9.",
		"## Security warnings for the running core\n\n- [SECURITY-3501](https://www.jenkins.io/security/advisory/2024-10-02/): Exposure of secrets\n",
		"1. **On 2.479.3, update plugins** (1 change(s)).",
		"2. Upgrade Jenkins to **2.492.3**, start it, then update plugins (0 change(s)).\n   - ⚠ gitlab-plugin 1.5 needs jersey3-api 2.0 or newer\n",
		"| both | 1.0 | 1.1 | [S-1](https://s1) fixed by the upgrade<br>[S-2](https://s2) **no fix yet**<br>[S-3](https://s3) affects the target version |",
		"| dep-link | 1 | 1 | [deprecated](https://plugins.jenkins.io/dep-link/) |",
		"| dep-plain | 1 | 1 | deprecated<br>no longer distributed by the update center",
		"| private | 0.1 | 0.1 | not listed by the update center",
		"| aix (offline) | – | – | needs Java 17 |",
		"## Not checked\n\n- agents could not be read\n",
	)
	lacks(t, out, "| quiet | 2 | 3 |  |", "stops at the last release of each Java generation")
}

func TestVerifyReport(t *testing.T) {
	for _, c := range []struct {
		name  string
		r     verify.Result
		wants []string
		not   []string
	}{
		{"go, still quiet", verify.Result{Verdict: verify.Go, Core: "2.580.1", Quiet: true,
			Checks: []verify.Check{{Name: "Plugins", State: verify.OK, Detail: "all 3 plugins that should run are active"}}},
			[]string{"# Upgrade check: Jenkins 2.580.1", "## GO: what worked before works", "Jenkins is still in quiet mode", "| Plugins | ✅ | all 3 plugins that should run are active |"},
			[]string{"ROLL BACK"}},
		{"go, open", verify.Result{Verdict: verify.Go, Core: "2.580.1"},
			[]string{"## GO"}, []string{"quiet mode"}},
		{"down", verify.Result{Verdict: verify.RollBack, Core: "2.580.1", Down: true,
			Checks: []verify.Check{{Name: "Jenkins is up", State: verify.Fail, Detail: "Jenkins does not answer (connection refused)"}},
			Errors: []string{"hudson.WebAppMain$3#run: Failed to initialize Jenkins"}},
			[]string{"## ROLL BACK\n\nJenkins did not come up (below).", "| Jenkins is up | ❌ | Jenkins does not answer (connection refused) |", "## Errors in the Jenkins log\n\n- hudson.WebAppMain$3#run: Failed to initialize Jenkins\n"},
			[]string{"Something that worked before"}},
		{"broken, quiet", verify.Result{Verdict: verify.RollBack, Core: "2.580.1", Quiet: true,
			Checks: []verify.Check{
				{Name: "Plugins", State: verify.Fail, Detail: "1 of 2 plugins that should run are not active", Items: []string{"workflow-cps"}},
				{Name: "Plugin versions", State: verify.Warn, Detail: "running, but | not planned", Items: []string{"git is 5.8, the plan says 5.9"}},
			},
			Roots: []dryrun.Root{{Failure: dryrun.Failure{Name: "workflow-cps", Version: "3000"}, Diagnosis: dryrun.Diagnosis{Why: "needs script-security 1336.v1 or newer", Fix: "update script-security"}, Cascade: []string{"pipeline-model"}},
				{Failure: dryrun.Failure{Name: "x"}, Diagnosis: dryrun.Diagnosis{Why: "Jenkins did not say why"}}},
			Missing:    []dryrun.ConfigProblem{{Item: "team/deploy", Problem: "did not load"}},
			Unreadable: []dryrun.ConfigProblem{{Item: "app", Problem: "CannotResolveClassException: x", Plugin: "foo 1.0", PluginState: dryrun.PluginNotInstalled}},
			Monitors:   []string{"Builds run as SYSTEM"}},
			[]string{"Something that worked before does not work now (below). Follow the rollback section of your runbook; Jenkins is still in quiet mode, so no build ran on the new version and restoring the backup loses nothing. Keep this report",
				"| Plugins | ❌ | 1 of 2 plugins that should run are not active: workflow-cps |",
				"| Plugin versions | ⚠️ | running, but \\| not planned: git is 5.8, the plan says 5.9 |",
				"| workflow-cps 3000 | needs script-security 1336.v1 or newer | update script-security | pipeline-model |",
				"| x | Jenkins did not say why | – | – |",
				"## Jobs and folders that did not load\n\n| Where | Comes from | Problem |\n|---|---|---|\n| team/deploy | – | did not load |",
				"## Settings Jenkins could not read\n\n| Where | Comes from | Problem |\n|---|---|---|\n| app | plugin foo 1.0, not installed | CannotResolveClassException: x |",
				"## Warnings Jenkins shows administrators\n\n- Builds run as SYSTEM\n"},
			nil},
		{"broken, open", verify.Result{Verdict: verify.RollBack, Core: "2.580.1"},
			[]string{"Follow the rollback section of your runbook. Keep this report"},
			[]string{"still in quiet mode"}},
	} {
		var b strings.Builder
		if err := Verify(&b, &c.r); err != nil {
			t.Fatal(err)
		}
		out := b.String()
		for _, want := range c.wants {
			if !strings.Contains(out, want) {
				t.Errorf("%s: report lacks %q\n%s", c.name, want, out)
			}
		}
		for _, u := range c.not {
			if strings.Contains(out, u) {
				t.Errorf("%s: report has %q\n%s", c.name, u, out)
			}
		}
	}
}

// The headline says which of the ways a rehearsal can end it was.
func TestRehearsalHeadlines(t *testing.T) {
	failed := []dryrun.Failure{{Name: "a", Version: "1", Reasons: []string{"Plugin is missing: b (1.0)"}}}
	for _, c := range []struct {
		name string
		r    dryrun.Result
		want string
	}{
		{"did not start", dryrun.Result{Reason: "not up after 10m0s"}, "**Jenkins did not finish starting** (not up after 10m0s).\n"},
		{"did not start, with failures", dryrun.Result{Reason: "the container exited", Failed: failed},
			"**Jenkins did not finish starting** (the container exited). The plugin failures below are the likely cause; the log has the rest."},
		{"clean", dryrun.Result{Started: true, StartupSeconds: 41.6}, "**Jenkins started in 42s and every plugin loaded.**\n"},
		{"clean with home", dryrun.Result{Started: true, StartupSeconds: 30, Home: &dryrun.HomeResult{Items: 7, Loaded: 7}},
			"**Jenkins started in 30s and every plugin loaded.** All 7 jobs and folders of your configuration came back."},
		{"failed and broken", dryrun.Result{Started: true, StartupSeconds: 30, Failed: failed, Broken: []dryrun.Failure{{Name: "email-ext", Reasons: []string{"Failed to scout x"}}}},
			"but 1 plugin(s) did not load, because of 1 cause(s).** 1 more loaded only in part. Upgrading"},
		{"broken only", dryrun.Result{Started: true, StartupSeconds: 30, Broken: []dryrun.Failure{{Name: "email-ext", Reasons: []string{"Failed to scout x"}}}},
			"**Jenkins started in 30s and every plugin loaded, but 1 loaded only in part.**"},
		{"configuration only", dryrun.Result{Started: true, StartupSeconds: 30, Home: &dryrun.HomeResult{Items: 2, Loaded: 1, Missing: []dryrun.ConfigProblem{{Item: "x", Problem: "y"}}}},
			"**Jenkins started in 30s and every plugin loaded, but your configuration did not come back whole** (below)."},
	} {
		c.r.Core, c.r.Image = "2.580.1", "jenkins/jenkins:2.580.1-jdk21"
		var b strings.Builder
		if err := Rehearsal(&b, &c.r); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(b.String(), c.want) {
			t.Errorf("%s: report lacks %q\n%s", c.name, c.want, b.String())
		}
	}
}

// Everything else a rehearsal can report.
func TestRehearsalDetails(t *testing.T) {
	var missing []dryrun.ConfigProblem
	for i := 0; i < 32; i++ {
		missing = append(missing, dryrun.ConfigProblem{Item: fmt.Sprintf("job%02d", i), Problem: "did not load"})
	}
	var copyErrs []string
	for i := 0; i < 12; i++ {
		copyErrs = append(copyErrs, fmt.Sprintf("secrets/f%d: permission denied", i))
	}
	r := &dryrun.Result{Core: "2.555.1", Image: "jenkins/jenkins:2.555.1-jdk25", Java: 25, JavaReason: "gradle needs Java 25 or newer",
		Plugins: 2, Upgrading: "2.346.3", Detached: []string{"sshd", "instance-identity"},
		Started: false, Reason: "the container exited",
		Broken: []dryrun.Failure{{Name: "email-ext", Reasons: []string{"Failed to scout A", "Failed to scout B"}}},
		Home: &dryrun.HomeResult{Source: "/srv/home", Files: 40, Items: 32, Missing: missing,
			Unreadable: []dryrun.ConfigProblem{{Item: "(global configuration)", Problem: "old field"}},
			Monitors:   []string{"Builds run as SYSTEM"}, CopyErrs: copyErrs},
		Skipped: []string{"private 1.0: no download known", "gone 2.0: GET x: 404 Not Found"},
		Errors:  []string{"some.Thing#run: boom"},
		LogFile: "/tmp/out/dry-run-2.555.1.log",
	}
	var b strings.Builder
	if err := Rehearsal(&b, r); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	contains(t, out,
		"Image `jenkins/jenkins:2.555.1-jdk25` (Java 25 because gradle needs Java 25 or newer).",
		"Jenkins started as an upgrade from 2.346.3 and installed 2 plugin(s) split out of core since then from its own war: sshd, instance-identity.\n",
		"## Loaded only in part\n",
		"| email-ext | – | Failed to scout A<br>Failed to scout B |",
		"Copied 40 configuration files from `/srv/home` (no build records, workspaces or logs).\n",
		"**Jobs and folders that did not load (32).**",
		"| job29 | – | did not load |\n| … | | 2 more |\n",
		"**Settings Jenkins could not read (1).**",
		"Warnings Jenkins shows administrators after the upgrade:\n\n- Builds run as SYSTEM\n",
		"Could not copy 12 file(s), so they were not rehearsed: secrets/f0: permission denied,",
		"secrets/f9: permission denied and 2 more\n",
		"Not rehearsed (could not be downloaded): private 1.0: no download known; gone 2.0: GET x: 404 Not Found.",
		"<details><summary>Other errors in the log</summary>\n\n```\nsome.Thing#run: boom\n```",
		"Full log: `/tmp/out/dry-run-2.555.1.log`",
	)
	// Jenkins did not start: no count of jobs that loaded.
	lacks(t, out, "| job30 |", "jobs and folders loaded")

	// Failures without roots get them worked out for the report; an
	// upgrade that installed no split plugins says only where it came from.
	b.Reset()
	r = &dryrun.Result{Core: "2.555.1", Started: true, Upgrading: "2.541.3",
		Failed: []dryrun.Failure{{Name: "a", Reasons: []string{"Plugin is missing: b (1.0)"}}, {Name: "c", Reasons: []string{"Failed to load: A (a 1)"}}}}
	if err := Rehearsal(&b, r); err != nil {
		t.Fatal(err)
	}
	contains(t, b.String(), "Jenkins started as an upgrade from 2.541.3.\n", "| a | needs the b plugin (1.0 or newer), which is not installed | install b 1.0 or newer | c |", "| c | – | Failed to load: A (a 1) |")
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("disk full") }

func TestWriteErrors(t *testing.T) {
	if err := Plan(failingWriter{}, &plan.Plan{}); err == nil || err.Error() != "disk full" {
		t.Errorf("Plan: %v", err)
	}
	if err := Verify(failingWriter{}, &verify.Result{}); err == nil || err.Error() != "disk full" {
		t.Errorf("Verify: %v", err)
	}
	if err := Rehearsal(failingWriter{}, &dryrun.Result{}); err == nil || err.Error() != "disk full" {
		t.Errorf("Rehearsal: %v", err)
	}
}
