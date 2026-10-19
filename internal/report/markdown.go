// Package report renders plans and rehearsals for people.
package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/keelinfra/jenkins-upgrade-check/internal/dryrun"
	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
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
		out = append(out, "not published on the update center (a private build?); test it yourself")
	}
	return out
}

// Rehearsal writes a Markdown report of a dry run.
func Rehearsal(w io.Writer, r *dryrun.Result) error {
	b := &strings.Builder{}
	fmt.Fprintf(b, "# Dry run: Jenkins %s with %d plugins\n\n", r.Core, r.Plugins)
	fmt.Fprintf(b, "Image `%s`.\n\n", r.Image)
	switch {
	case !r.Started:
		fmt.Fprintf(b, "**Jenkins did not finish starting** (%s).\n\n", r.Reason)
	case len(r.Failed) == 0:
		fmt.Fprintf(b, "**Jenkins started in %.0fs and every plugin loaded.**\n\n", r.StartupSeconds)
	default:
		fmt.Fprintf(b, "**Jenkins started in %.0fs; %d plugin(s) failed to load.**\n\n", r.StartupSeconds, len(r.Failed))
	}
	if len(r.Failed) > 0 {
		b.WriteString("| Plugin | Version | Why |\n|---|---|---|\n")
		for _, f := range r.Failed {
			fmt.Fprintf(b, "| %s | %s | %s |\n", f.Name, dash(f.Version), dash(strings.Join(f.Reasons, "<br>")))
		}
		b.WriteString("\n")
	}
	if len(r.Skipped) > 0 {
		fmt.Fprintf(b, "Not rehearsed (no download on the update center): %s.\n\n", strings.Join(r.Skipped, ", "))
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
