// Package verify checks a controller right after an upgrade against what it
// was before, and says whether to keep the upgrade or roll back, and why.
package verify

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/keelinfra/jenkins-upgrade-check/internal/dryrun"
	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
	"github.com/keelinfra/jenkins-upgrade-check/internal/snapshot"
)

// Verdicts.
const (
	Go       = "GO"
	RollBack = "ROLL BACK"
)

// Check states.
const (
	OK   = "ok"
	Warn = "check"
	Fail = "fail"
)

// Input is what verify compares.
type Input struct {
	// Before is the controller as collected before the upgrade.
	Before *snapshot.Snapshot
	// After is the controller as collected now.
	After *snapshot.Snapshot
	// Core is the version the upgrade was meant to reach.
	Core string
	// Plugins is the plugin set the upgrade was meant to install, if known.
	Plugins []plan.Pick
	// Probe is what the probe script saw as Jenkins started, if it ran.
	Probe *dryrun.Probe
	// Home is JENKINS_HOME, to name the plugin a broken job needs.
	Home string
	// Log is the Jenkins log of this start, if given.
	Log string
	// Down says why the controller could not be read, when it could not:
	// then After is not used.
	Down string
}

// Check is one thing verified.
type Check struct {
	Name   string   `json:"name"`
	State  string   `json:"state"`
	Detail string   `json:"detail"`
	Items  []string `json:"items,omitempty"`
}

// Result is the verdict and what it rests on.
type Result struct {
	Verdict string  `json:"verdict"`
	Core    string  `json:"core"`
	Checks  []Check `json:"checks"`
	// Failed are plugins that worked before and do not now, by cause.
	Roots []dryrun.Root `json:"roots,omitempty"`
	// Missing and Unreadable are jobs and settings that did not come back.
	Missing    []dryrun.ConfigProblem `json:"missing,omitempty"`
	Unreadable []dryrun.ConfigProblem `json:"unreadable,omitempty"`
	Monitors   []string               `json:"monitors,omitempty"`
	Quiet      bool                   `json:"quietingDown"`
	// Errors are severe errors from the Jenkins log, when it did not start.
	Errors []string `json:"errors,omitempty"`
	// Down is set when Jenkins did not answer.
	Down bool `json:"down,omitempty"`
}

// Run compares.
func Run(in Input) *Result {
	if in.Down != "" {
		return down(in)
	}
	r := &Result{Core: in.After.Core, Quiet: in.After.QuietingDown}
	add := func(c Check) { r.Checks = append(r.Checks, c) }

	// The version.
	switch {
	case in.Core == "":
		add(Check{Name: "Jenkins version", State: OK, Detail: fmt.Sprintf("%s (was %s)", in.After.Core, in.Before.Core)})
	case in.After.Core == in.Core:
		add(Check{Name: "Jenkins version", State: OK, Detail: fmt.Sprintf("%s, as planned (was %s)", in.After.Core, in.Before.Core)})
	default:
		add(Check{Name: "Jenkins version", State: Fail, Detail: fmt.Sprintf("%s, but the upgrade was to %s: the new version did not start", in.After.Core, in.Core)})
	}

	// Plugins: everything active before, and everything the plan adds,
	// must be active now.
	now := map[string]snapshot.Plugin{}
	for _, p := range in.After.Plugins {
		now[p.Name] = p
	}
	want := map[string]string{} // name → planned version ("" if unknown)
	for _, p := range in.Before.Plugins {
		if p.Active {
			want[p.Name] = ""
		}
	}
	for _, p := range in.Plugins {
		if _, ok := want[p.Name]; ok || p.Added {
			want[p.Name] = p.Version
		}
	}
	var down, wrongVersion []string
	for name, v := range want {
		p, ok := now[name]
		if !ok || !p.Active {
			down = append(down, name)
			continue
		}
		if v != "" && p.Version != v {
			wrongVersion = append(wrongVersion, fmt.Sprintf("%s is %s, the plan says %s", name, p.Version, v))
		}
	}
	sort.Strings(down)
	sort.Strings(wrongVersion)
	if len(down) == 0 {
		add(Check{Name: "Plugins", State: OK, Detail: fmt.Sprintf("all %d plugins that should run are active", len(want))})
	} else {
		add(Check{Name: "Plugins", State: Fail, Detail: fmt.Sprintf("%d of %d plugins that should run are not active", len(down), len(want)), Items: down})
		r.Roots = dryrun.Roots(reasons(down, in))
	}
	if len(wrongVersion) > 0 {
		add(Check{Name: "Plugin versions", State: Warn, Detail: "running, but not the planned version: was the new file copied into plugins/?", Items: wrongVersion})
	}

	// Jobs and folders.
	if len(in.Before.Items) > 0 {
		have := map[string]bool{}
		for _, it := range in.After.Items {
			have[it] = true
		}
		var missing []string
		for _, it := range in.Before.Items {
			if !have[it] {
				missing = append(missing, it)
			}
		}
		if len(missing) == 0 {
			add(Check{Name: "Jobs and folders", State: OK, Detail: fmt.Sprintf("all %d are back", len(in.Before.Items))})
		} else {
			add(Check{Name: "Jobs and folders", State: Fail, Detail: fmt.Sprintf("%d of %d did not load", len(missing), len(in.Before.Items)), Items: missing})
			r.Missing = blameItems(missing, in)
		}
	} else {
		add(Check{Name: "Jobs and folders", State: Warn, Detail: "not compared: the snapshot from before the upgrade has no job list (collect it again with this version)"})
	}

	// Settings Jenkins could not read.
	if in.Probe != nil {
		failed := map[string]bool{}
		for _, d := range down {
			failed[d] = true
		}
		active := in.Probe.Active()
		wasActive := map[string]bool{}
		for _, p := range in.Before.Plugins {
			wasActive[p.Name] = p.Active
		}
		lost := 0
		for _, u := range in.Probe.Unreadable {
			dryrun.Blame(&u, in.Home, failed, active)
			r.Unreadable = append(r.Unreadable, u)
			// A setting from a plugin that ran before and does not now is
			// lost by the upgrade; one from a plugin that was already gone
			// was unreadable before too.
			if name, _, _ := strings.Cut(u.Plugin, " "); name != "" && wasActive[name] && !active[name] {
				lost++
			}
		}
		switch {
		case len(r.Unreadable) == 0:
			add(Check{Name: "Settings", State: OK, Detail: "Jenkins read all configuration"})
		case lost > 0:
			add(Check{Name: "Settings", State: Fail, Detail: fmt.Sprintf("%d setting(s) come from plugins that ran before and do not now; Jenkins drops them when the item is next saved", lost)})
		default:
			add(Check{Name: "Settings", State: Warn, Detail: fmt.Sprintf("%d setting(s) Jenkins cannot read, from plugins that were not running before the upgrade either, or old fields this version ignores; see Manage Jenkins > Manage Old Data", len(r.Unreadable))})
		}
		for _, m := range in.Probe.Monitors {
			if !dryrun.IsNoise(m.ID) {
				r.Monitors = append(r.Monitors, m.Name)
			}
		}
	}

	// Agents.
	wasOnline := map[string]bool{}
	for _, a := range in.Before.Agents {
		if !a.Offline {
			wasOnline[a.Name] = true
		}
	}
	var gone []string
	for _, a := range in.After.Agents {
		if wasOnline[a.Name] && a.Offline {
			gone = append(gone, a.Name)
		}
	}
	sort.Strings(gone)
	switch {
	case len(wasOnline) == 0:
	case len(gone) == 0:
		add(Check{Name: "Agents", State: OK, Detail: fmt.Sprintf("all %d agents that were online are online", len(wasOnline))})
	default:
		add(Check{Name: "Agents", State: Warn, Detail: fmt.Sprintf("%d agent(s) online before are offline now. Agents can take a few minutes to reconnect; run verify again. If they stay offline, check their Java version and logs", len(gone)), Items: gone})
	}

	r.Verdict = Go
	for _, c := range r.Checks {
		if c.State == Fail {
			r.Verdict = RollBack
		}
	}
	return r
}

// down is the verdict on a controller that does not answer: it did not
// start, or it is stuck. The log says why.
func down(in Input) *Result {
	r := &Result{Core: in.Core, Verdict: RollBack, Down: true}
	r.Checks = append(r.Checks, Check{Name: "Jenkins is up", State: Fail, Detail: "Jenkins does not answer (" + in.Down + "). If it is still starting, wait and run verify again; if it failed, the log below says why"})
	if in.Log == "" {
		r.Checks[0].Detail += "; pass the Jenkins log with --log to see why"
		return r
	}
	failed, errs := dryrun.ParseLog(in.Log)
	r.Roots = dryrun.Roots(failed)
	if in.Home != "" {
		byName := map[string]bool{}
		for _, f := range failed {
			byName[f.Name] = true
		}
		var broken []dryrun.Failure
		broken, errs = dryrun.Attribute(errs, filepath.Join(in.Home, "plugins"), byName)
		r.Roots = append(r.Roots, dryrun.Roots(broken)...)
	}
	if len(errs) > 10 {
		errs = errs[:10]
	}
	r.Errors = errs
	return r
}

// reasons finds why each plugin is down: from the probe, else the log.
func reasons(down []string, in Input) []dryrun.Failure {
	known := map[string]dryrun.Failure{}
	if in.Log != "" {
		fs, _ := dryrun.ParseLog(in.Log)
		for _, f := range fs {
			known[f.Name] = f
		}
	}
	if in.Probe != nil {
		for _, f := range in.Probe.Failures() {
			if k, ok := known[f.Name]; !ok || len(k.Reasons) == 0 {
				known[f.Name] = f
			}
		}
	}
	installed := map[string]string{}
	for _, p := range in.After.Plugins {
		installed[p.Name] = p.Version
	}
	var out []dryrun.Failure
	for _, name := range down {
		f, ok := known[name]
		if !ok {
			f = dryrun.Failure{Name: name}
			if v, ok := installed[name]; ok {
				f.Version = v
				f.Reasons = []string{"installed but not active; pass the Jenkins log with --log, or --home after a start with the probe script, to see why"}
			} else {
				f.Reasons = []string{"not installed any more, or failed to load; pass the Jenkins log with --log, or --home after a start with the probe script, to see why"}
			}
		}
		if f.Version == "" {
			f.Version = installed[name]
		}
		out = append(out, f)
	}
	return out
}

// blameItems names the plugin each missing job needs, from its config.xml.
func blameItems(missing []string, in Input) []dryrun.ConfigProblem {
	failed := map[string]bool{}
	active := map[string]bool{}
	for _, p := range in.After.Plugins {
		if p.Active {
			active[p.Name] = true
		}
	}
	for _, c := range in.Before.Plugins {
		if c.Active && !active[c.Name] {
			failed[c.Name] = true
		}
	}
	var itemErrs map[string]string
	if in.Log != "" {
		_, errs := dryrun.ParseLog(in.Log)
		_, itemErrs = dryrun.ItemErrors(errs)
	}
	var out []dryrun.ConfigProblem
	for _, m := range missing {
		p := dryrun.ConfigProblem{Item: m, Path: dryrun.ItemPath(m), Problem: "did not load"}
		base := m[strings.LastIndex(m, "/")+1:]
		if e := itemErrs[base]; e != "" {
			p.Problem = dryrun.ShortClass(e)
		}
		if in.Home != "" {
			dryrun.Blame(&p, in.Home, failed, active)
		}
		out = append(out, p)
	}
	return out
}
