// Package report renders plans and rehearsals for people.
package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/keelinfra/jenkins-upgrade-check/internal/dryrun"
	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
	"github.com/keelinfra/jenkins-upgrade-check/internal/verify"
)

const footer = "\n---\n" +
	"This report is a plan, not a guarantee: plugins that break on an API change only show it when they are started. " +
	"Run `jenkins-upgrade-check dry-run` to start your plugin set on the target version in a throwaway container. " +
	"If you would rather not do the upgrade alone, keelinfra does Jenkins upgrade checks and upgrades: https://keelinfra.io/jenkins/pricing/\n"

// Plan writes a Markdown report.
func Plan(w io.Writer, p *plan.Plan) error {
	b := &strings.Builder{}
	s := p.Summary
	fmt.Fprintf(b, "# Jenkins upgrade plan: %s → %s\n\n", p.From, p.To)
	fmt.Fprintf(b, "Source: %s · generated %s by %s %s\n\n", p.Source, p.GeneratedAt.Format("2006-01-02 15:04 MST"), p.Tool, p.ToolVersion)

	fmt.Fprintf(b, "**%d core upgrade(s)**, %d plugins: %d to upgrade, %d new dependencies", s.CoreUpgrades, s.Plugins, s.Upgraded, s.Added)
	var flags []string
	if s.VulnerableNow > 0 {
		flags = append(flags, fmt.Sprintf("%d with security warnings today", s.VulnerableNow))
	}
	if s.VulnerableAfter > 0 {
		flags = append(flags, fmt.Sprintf("%d still with warnings after the upgrade", s.VulnerableAfter))
	}
	if s.Deprecated > 0 {
		flags = append(flags, fmt.Sprintf("%d deprecated", s.Deprecated))
	}
	if s.NotDistributed > 0 {
		flags = append(flags, fmt.Sprintf("%d no longer distributed", s.NotDistributed))
	}
	if s.Unknown > 0 {
		flags = append(flags, fmt.Sprintf("%d not from the update center", s.Unknown))
	}
	if s.Problems > 0 {
		flags = append(flags, fmt.Sprintf("%d dependency problem(s)", s.Problems))
	}
	if s.Agents > 0 {
		flags = append(flags, fmt.Sprintf("%d agent(s) that cannot follow", s.Agents))
	}
	if len(flags) > 0 {
		fmt.Fprintf(b, "; %s", strings.Join(flags, ", "))
	}
	b.WriteString(".\n\n")

	fmt.Fprintf(b, "Jenkins %s runs on Java %s; Jenkins %s runs on Java %s.", p.From, ints(p.JavaFrom), p.To, ints(p.JavaTo))
	if p.ControllerJava != "" {
		fmt.Fprintf(b, " The controller runs Java %s.", p.ControllerJava)
	}
	b.WriteString("\n\n")

	if len(p.CoreWarnings) > 0 {
		b.WriteString("## Security warnings for the running core\n\n")
		for _, w := range p.CoreWarnings {
			fmt.Fprintf(b, "- [%s](%s): %s\n", w.ID, w.URL, w.Message)
		}
		b.WriteString("\n")
	}

	b.WriteString("## Steps\n\n")
	for i, st := range p.Steps {
		switch {
		case !st.Upgrade:
			fmt.Fprintf(b, "%d. **On %s, update plugins** (%d change(s)). Back up `$JENKINS_HOME` first.\n", i+1, st.Core, len(st.Changes))
		case st.Java > 0:
			fmt.Fprintf(b, "%d. **Run the controller and every agent on Java %d**, upgrade Jenkins to **%s**, start it, then update plugins (%d change(s)).\n", i+1, st.Java, st.Core, len(st.Changes))
		default:
			fmt.Fprintf(b, "%d. Upgrade Jenkins to **%s**, start it, then update plugins (%d change(s)).\n", i+1, st.Core, len(st.Changes))
		}
		for _, pr := range st.Problems {
			fmt.Fprintf(b, "   - ⚠ %s\n", pr)
		}
	}
	if len(p.Steps) > 2 {
		b.WriteString("\nNo Java version runs both the current and the target core, so the upgrade stops at the last release of each Java generation on the way.\n")
	}
	b.WriteString("\n")

	var attention []plan.PluginPlan
	for _, pp := range p.Plugins {
		if pp.Attention() {
			attention = append(attention, pp)
		}
	}
	if len(attention) > 0 {
		b.WriteString("## Needs a decision\n\n| Plugin | Now | After | Why |\n|---|---|---|---|\n")
		for _, pp := range attention {
			fmt.Fprintf(b, "| %s | %s | %s | %s |\n", pp.Name, dash(pp.Current), pp.Target, strings.Join(reasons(pp), "<br>"))
		}
		b.WriteString("\n")
	}

	if len(p.Agents) > 0 {
		b.WriteString("## Agents\n\n| Agent | OS | Java | Problem |\n|---|---|---|---|\n")
		for _, a := range p.Agents {
			name := a.Name
			if a.Offline {
				name += " (offline)"
			}
			fmt.Fprintf(b, "| %s | %s | %s | %s |\n", name, dash(a.OS), dash(a.Java), a.Problem)
		}
		b.WriteString("\n")
	}

	if len(p.Notes) > 0 {
		b.WriteString("## Not checked\n\n")
		for _, n := range p.Notes {
			fmt.Fprintf(b, "- %s\n", n)
		}
		b.WriteString("\n")
	}

	b.WriteString("<details><summary>Every plugin</summary>\n\n| Plugin | Now | After |\n|---|---|---|\n")
	for _, pp := range p.Plugins {
		after := pp.Target
		if pp.Added {
			after += " (new, for " + strings.Join(pp.RequiredBy, ", ") + ")"
		}
		fmt.Fprintf(b, "| %s | %s | %s |\n", pp.Name, dash(pp.Current), after)
	}
	b.WriteString("\n</details>\n")
	b.WriteString(footer)
	_, err := io.WriteString(w, b.String())
	return err
}

func reasons(pp plan.PluginPlan) []string {
	var out []string
	for _, w := range pp.VulnerableNow {
		fixed := "fixed by the upgrade"
		for _, a := range pp.VulnerableAfter {
			if a.ID == w.ID {
				fixed = "**no fix yet**"
			}
		}
		out = append(out, fmt.Sprintf("[%s](%s) %s", w.ID, w.URL, fixed))
	}
	for _, w := range pp.VulnerableAfter {
		dup := false
		for _, n := range pp.VulnerableNow {
			dup = dup || n.ID == w.ID
		}
		if !dup {
			out = append(out, fmt.Sprintf("[%s](%s) affects the target version", w.ID, w.URL))
		}
	}
	if pp.Deprecated != "" {
		if strings.HasPrefix(pp.Deprecated, "http") {
			out = append(out, fmt.Sprintf("[deprecated](%s)", pp.Deprecated))
		} else {
			out = append(out, "deprecated")
		}
	}
	if pp.NotDistributed {
		out = append(out, "no longer distributed by the update center (suspended or removed); find a replacement or remove it")
	}
	if pp.Unknown {
		out = append(out, "not listed by the update center (removed from distribution, or a private build); it stays at this version, so find a replacement or test it yourself")
	}
	return out
}

// Verify writes the verdict on an upgraded controller.
func Verify(w io.Writer, r *verify.Result) error {
	b := &strings.Builder{}
	fmt.Fprintf(b, "# Upgrade check: Jenkins %s\n\n", r.Core)
	if r.Verdict == verify.Go {
		b.WriteString("## GO: what worked before works\n\n")
		if r.Quiet {
			b.WriteString("Jenkins is still in quiet mode, so no build has run on the new version yet. Open it to users: Manage Jenkins > Cancel Shutdown (or `POST /cancelQuietDown`), and remove `init.groovy.d/zz-jenkins-upgrade-check.groovy` from JENKINS_HOME.\n\n")
		}
	} else {
		b.WriteString("## ROLL BACK\n\n")
		if r.Down {
			b.WriteString("Jenkins did not come up (below). Follow the rollback section of your runbook. Keep this report and the Jenkins log: they say what to fix before the next attempt.\n\n")
		} else {
			b.WriteString("Something that worked before does not work now (below). Follow the rollback section of your runbook")
		}
		if !r.Down {
			if r.Quiet {
				b.WriteString("; Jenkins is still in quiet mode, so no build ran on the new version and restoring the backup loses nothing")
			}
			b.WriteString(". Keep this report and the Jenkins log: they say what to fix before the next attempt.\n\n")
		}
	}
	b.WriteString("| Check | | Result |\n|---|---|---|\n")
	for _, c := range r.Checks {
		mark := map[string]string{verify.OK: "✅", verify.Warn: "⚠️", verify.Fail: "❌"}[c.State]
		detail := c.Detail
		if len(c.Items) > 0 {
			detail += ": " + list(c.Items, 15)
		}
		fmt.Fprintf(b, "| %s | %s | %s |\n", c.Name, mark, cell(detail))
	}
	b.WriteString("\n")
	if len(r.Roots) > 0 {
		b.WriteString("## Why plugins did not load\n\n| Plugin | Why | What to do | Also fails because of it |\n|---|---|---|---|\n")
		for _, rt := range r.Roots {
			name := rt.Name
			if rt.Version != "" {
				name += " " + rt.Version
			}
			fmt.Fprintf(b, "| %s | %s | %s | %s |\n", cell(name), cell(rt.Diagnosis.Why), cell(dash(rt.Diagnosis.Fix)), cell(dash(list(rt.Cascade, 12))))
		}
		b.WriteString("\n")
	}
	if len(r.Errors) > 0 {
		b.WriteString("## Errors in the Jenkins log\n\n")
		for _, e := range r.Errors {
			fmt.Fprintf(b, "- %s\n", e)
		}
		b.WriteString("\n")
	}
	if len(r.Missing) > 0 {
		b.WriteString("## Jobs and folders that did not load\n\n")
		configTable(b, r.Missing)
	}
	if len(r.Unreadable) > 0 {
		b.WriteString("## Settings Jenkins could not read\n\n")
		configTable(b, r.Unreadable)
	}
	if len(r.Monitors) > 0 {
		b.WriteString("## Warnings Jenkins shows administrators\n\n")
		for _, m := range r.Monitors {
			fmt.Fprintf(b, "- %s\n", m)
		}
		b.WriteString("\n")
	}
	_, err := io.WriteString(w, b.String())
	return err
}

// Rehearsal writes a Markdown report of a dry run.
func Rehearsal(w io.Writer, r *dryrun.Result) error {
	b := &strings.Builder{}
	fmt.Fprintf(b, "# Dry run: Jenkins %s with %d plugins\n\n", r.Core, r.Plugins)
	fmt.Fprintf(b, "Image `%s`", r.Image)
	if r.JavaReason != "" {
		fmt.Fprintf(b, " (Java %d because %s)", r.Java, r.JavaReason)
	}
	b.WriteString(". The container had no network and no executors, so nothing outside it was reached and no build ran.\n\n")

	roots := r.Roots
	if len(roots) == 0 && len(r.Failed) > 0 {
		roots = dryrun.Roots(r.Failed)
	}
	switch {
	case !r.Started:
		fmt.Fprintf(b, "**Jenkins did not finish starting** (%s).", r.Reason)
		if len(roots) > 0 {
			b.WriteString(" The plugin failures below are the likely cause; the log has the rest.")
		}
		b.WriteString("\n\n")
	case r.OK():
		fmt.Fprintf(b, "**Jenkins started in %.0fs and every plugin loaded.**", r.StartupSeconds)
		if r.Home != nil {
			fmt.Fprintf(b, " All %d jobs and folders of your configuration came back.", r.Home.Items)
		}
		b.WriteString("\n\n")
	case len(r.Failed) > 0:
		fmt.Fprintf(b, "**Jenkins started in %.0fs, but %d plugin(s) did not load", r.StartupSeconds, len(r.Failed))
		if len(roots) > 0 {
			fmt.Fprintf(b, ", because of %d cause(s)", len(roots))
		}
		b.WriteString(".**")
		if len(r.Broken) > 0 {
			fmt.Fprintf(b, " %d more loaded only in part.", len(r.Broken))
		}
		b.WriteString(" Upgrading a controller with this set would break the jobs that use them.\n\n")
	case len(r.Broken) > 0:
		fmt.Fprintf(b, "**Jenkins started in %.0fs and every plugin loaded, but %d loaded only in part.**\n\n", r.StartupSeconds, len(r.Broken))
	default:
		fmt.Fprintf(b, "**Jenkins started in %.0fs and every plugin loaded, but your configuration did not come back whole** (below).\n\n", r.StartupSeconds)
	}

	if len(roots) > 0 {
		b.WriteString("## What failed, and what to do\n\n")
		b.WriteString("Each row is one cause. Plugins that failed only because another one did are listed with it: fix the first and they load too.\n\n")
		b.WriteString("| Plugin | Why | What to do | Also fails because of it |\n|---|---|---|---|\n")
		for _, rt := range roots {
			name := rt.Name
			if rt.Version != "" {
				name += " " + rt.Version
			}
			fmt.Fprintf(b, "| %s | %s | %s | %s |\n", cell(name), cell(rt.Diagnosis.Why), cell(dash(rt.Diagnosis.Fix)), cell(dash(list(rt.Cascade, 12))))
		}
		b.WriteString("\n<details><summary>What Jenkins said for each plugin</summary>\n\n| Plugin | Version | Jenkins said |\n|---|---|---|\n")
		for _, f := range r.Failed {
			fmt.Fprintf(b, "| %s | %s | %s |\n", f.Name, dash(f.Version), cell(dash(strings.Join(f.Reasons, "\n"))))
		}
		b.WriteString("\n</details>\n\n")
	}
	if len(r.Broken) > 0 {
		b.WriteString("## Loaded only in part\n\nThese plugins started, but some of their parts did not; the features those parts provide will be missing.\n\n| Plugin | Version | Error |\n|---|---|---|\n")
		for _, f := range r.Broken {
			fmt.Fprintf(b, "| %s | %s | %s |\n", f.Name, dash(f.Version), cell(strings.Join(f.Reasons, "\n")))
		}
		b.WriteString("\n")
	}
	if h := r.Home; h != nil {
		b.WriteString("## Your configuration\n\n")
		fmt.Fprintf(b, "Copied %d configuration files from `%s` (no build records, workspaces or logs).", h.Files, h.Source)
		if r.Started {
			fmt.Fprintf(b, " %d of %d jobs and folders loaded.", h.Loaded, h.Items)
		}
		b.WriteString("\n\n")
		if len(h.Missing) > 0 {
			fmt.Fprintf(b, "**Jobs and folders that did not load (%d).** Jenkins leaves them out after the upgrade, as if deleted, until the plugin they need loads:\n\n", len(h.Missing))
			configTable(b, h.Missing)
		}
		if len(h.Unreadable) > 0 {
			fmt.Fprintf(b, "**Settings Jenkins could not read (%d).** The item loads without them, and Jenkins drops them for good when the item is next saved:\n\n", len(h.Unreadable))
			configTable(b, h.Unreadable)
		}
		if len(h.Monitors) > 0 {
			b.WriteString("Warnings Jenkins shows administrators after the upgrade:\n\n")
			for _, m := range h.Monitors {
				fmt.Fprintf(b, "- %s\n", m)
			}
			b.WriteString("\n")
		}
		if len(h.CopyErrs) > 0 {
			fmt.Fprintf(b, "Could not copy %d file(s), so they were not rehearsed: %s\n\n", len(h.CopyErrs), list(h.CopyErrs, 10))
		}
	}
	if len(r.Skipped) > 0 {
		fmt.Fprintf(b, "Not rehearsed (could not be downloaded): %s.\n\n", strings.Join(r.Skipped, "; "))
	}
	if len(r.Errors) > 0 {
		b.WriteString("<details><summary>Other errors in the log</summary>\n\n```\n")
		for _, e := range r.Errors {
			b.WriteString(e + "\n")
		}
		b.WriteString("```\n\n</details>\n")
	}
	if r.LogFile != "" {
		fmt.Fprintf(b, "\nFull log: `%s`\n", r.LogFile)
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func configTable(b *strings.Builder, ps []dryrun.ConfigProblem) {
	b.WriteString("| Where | Comes from | Problem |\n|---|---|---|\n")
	for i, p := range ps {
		if i == 30 {
			fmt.Fprintf(b, "| … | | %d more |\n", len(ps)-30)
			break
		}
		from := "–"
		if p.Plugin != "" {
			from = "plugin " + p.Plugin + ", " + p.PluginState
		}
		fmt.Fprintf(b, "| %s | %s | %s |\n", cell(p.Item), cell(from), cell(p.Problem))
	}
	b.WriteString("\n")
}

// cell keeps text inside one Markdown table cell.
func cell(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	return strings.ReplaceAll(s, "\n", "<br>")
}

func list(xs []string, max int) string {
	if len(xs) > max {
		return strings.Join(xs[:max], ", ") + fmt.Sprintf(" and %d more", len(xs)-max)
	}
	return strings.Join(xs, ", ")
}

func dash(s string) string {
	if s == "" {
		return "–"
	}
	return s
}

func ints(xs []int) string {
	s := make([]string, len(xs))
	for i, x := range xs {
		s[i] = fmt.Sprint(x)
	}
	return strings.Join(s, ", ")
}
