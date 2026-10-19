package dryrun

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Cause kinds, from the most to the least specific.
const (
	CauseCoreTooOld        = "core-too-old"
	CauseJavaTooOld        = "java-too-old"
	CauseMissingDependency = "missing-dependency"
	CauseOldDependency     = "old-dependency"
	CauseDisabledDep       = "disabled-dependency"
	CauseDependencyFailed  = "dependency-failed"
	CauseJakartaEE9        = "jakarta-ee9"
	CauseRemovedAPI        = "removed-api"
	CauseNewerJava         = "built-for-newer-java"
	CauseInactive          = "not-active"
	CauseOther             = "other"
)

// Diagnosis is a plain-language reading of why a plugin failed.
type Diagnosis struct {
	Kind string `json:"kind"`
	// Because names the failed plugin this one depends on, for a failure
	// that only follows from another.
	Because string `json:"because,omitempty"`
	Why     string `json:"why"`
	Fix     string `json:"fix,omitempty"`
}

var (
	reMissing     = regexp.MustCompile(`Plugin is missing: (\S+) \(([^)]+)\)`)
	reUpdate      = regexp.MustCompile(`Update required: (.+?) \((\S+) (\S+)\) to be updated to (\S+) or higher`)
	reCore        = regexp.MustCompile(`Jenkins \((\S+)\) or higher required`)
	reJava        = regexp.MustCompile(`Java \((\d+)\) or higher required`)
	reFailedDep   = regexp.MustCompile(`Failed to load: (.+?) \((\S+) (\S+)\)`)
	reFailedDepV1 = regexp.MustCompile(`(.+?) \((\S+) (\S+)\) failed to load\. Fix this plugin first`)
	reDisabled    = regexp.MustCompile(`Plugin is disabled: (.+?)(?: \((\S+)\))?$`)
	reClassFile   = regexp.MustCompile(`class file version (\d+)(?:\.\d+)?`)
	reMissingCls  = regexp.MustCompile(`(?:NoClassDefFoundError|ClassNotFoundException):\s*(\S+)`)
)

// Explain reads the reasons Jenkins gave for one plugin.
func Explain(f Failure) Diagnosis {
	all := strings.Join(f.Reasons, "\n")
	// Dependency problems first: they are what Jenkins checks before it
	// loads any code, and a cascade is not a cause of its own.
	for _, r := range f.Reasons {
		if m := reCore.FindStringSubmatch(r); m != nil {
			return Diagnosis{Kind: CauseCoreTooOld,
				Why: fmt.Sprintf("this release needs Jenkins %s or newer", m[1]),
				Fix: "use an older release of the plugin, or upgrade Jenkins first"}
		}
		if m := reJava.FindStringSubmatch(r); m != nil {
			return Diagnosis{Kind: CauseJavaTooOld,
				Why: fmt.Sprintf("this release needs Java %s or newer", m[1]),
				Fix: fmt.Sprintf("run the controller on Java %s or newer, or use an older release of the plugin", m[1])}
		}
	}
	for _, r := range f.Reasons {
		if m := reMissing.FindStringSubmatch(r); m != nil {
			return Diagnosis{Kind: CauseMissingDependency,
				Why: fmt.Sprintf("needs the %s plugin (%s or newer), which is not installed", m[1], m[2]),
				Fix: fmt.Sprintf("install %s %s or newer", m[1], m[2])}
		}
		if m := reUpdate.FindStringSubmatch(r); m != nil {
			return Diagnosis{Kind: CauseOldDependency,
				Why: fmt.Sprintf("needs %s %s or newer; %s is installed", m[2], m[4], m[3]),
				Fix: fmt.Sprintf("update %s to %s or newer", m[2], m[4])}
		}
		if m := reDisabled.FindStringSubmatch(r); m != nil {
			return Diagnosis{Kind: CauseDisabledDep,
				Why: fmt.Sprintf("needs %s, which is disabled", m[1]),
				Fix: fmt.Sprintf("enable %s", m[1])}
		}
	}
	for _, r := range f.Reasons {
		var m []string
		if m = reFailedDep.FindStringSubmatch(r); m == nil {
			m = reFailedDepV1.FindStringSubmatch(r)
		}
		if m != nil {
			return Diagnosis{Kind: CauseDependencyFailed, Because: m[2],
				Why: fmt.Sprintf("needs %s, which failed to load", m[2]),
				Fix: fmt.Sprintf("fix %s first", m[2])}
		}
	}
	if m := reClassFile.FindStringSubmatch(all); m != nil && strings.Contains(all, "UnsupportedClassVersionError") {
		n, _ := strconv.Atoi(m[1])
		j := n - 44
		return Diagnosis{Kind: CauseNewerJava,
			Why: fmt.Sprintf("built for Java %d, newer than the Java this Jenkins runs on", j),
			Fix: fmt.Sprintf("run Jenkins on Java %d or newer, or use an older release of the plugin", j)}
	}
	if m := reMissingCls.FindStringSubmatch(all); m != nil {
		cls := strings.ReplaceAll(m[1], "/", ".")
		if strings.HasPrefix(cls, "javax.servlet.") {
			return Diagnosis{Kind: CauseJakartaEE9,
				Why: fmt.Sprintf("built for the old servlet API (%s); Jenkins 2.475 and newer moved to Jakarta EE 9", cls),
				Fix: "update the plugin to a release made for Jenkins 2.475 or newer, or remove it"}
		}
		return Diagnosis{Kind: CauseRemovedAPI,
			Why: fmt.Sprintf("uses %s, which this Jenkins and its plugins no longer provide", cls),
			Fix: "update the plugin; if no release fixes it, the plugin is abandoned and needs a replacement"}
	}
	for _, r := range f.Reasons {
		if r == reasonInactive {
			return Diagnosis{Kind: CauseInactive, Why: "installed but not active after startup", Fix: "read the log for this plugin"}
		}
	}
	why := "Jenkins did not say why"
	if len(f.Reasons) > 0 {
		why = f.Reasons[0]
	}
	return Diagnosis{Kind: CauseOther, Why: why, Fix: "read the log for this plugin"}
}

const reasonInactive = "not active after startup"

// Root is one cause and everything that failed because of it.
type Root struct {
	Failure
	Diagnosis Diagnosis `json:"diagnosis"`
	// Cascade lists the plugins that failed only because this one did.
	Cascade []string `json:"cascade,omitempty"`
}

// Roots groups failures by their root cause: a plugin that failed only
// because a dependency failed is listed under that dependency.
func Roots(failed []Failure) []Root {
	diag := map[string]Diagnosis{}
	byName := map[string]Failure{}
	for _, f := range failed {
		diag[f.Name] = Explain(f)
		byName[f.Name] = f
	}
	rootOf := func(name string) string {
		seen := map[string]bool{}
		for {
			d, ok := diag[name]
			if !ok || d.Kind != CauseDependencyFailed || seen[name] {
				return name
			}
			seen[name] = true
			if _, failedToo := diag[d.Because]; !failedToo {
				return name
			}
			name = d.Because
		}
	}
	idx := map[string]int{}
	var out []Root
	for _, f := range failed {
		r := rootOf(f.Name)
		if r == f.Name {
			continue
		}
		i, ok := idx[r]
		if !ok {
			i = len(out)
			idx[r] = i
			out = append(out, Root{Failure: byName[r], Diagnosis: diag[r]})
		}
		out[i].Cascade = append(out[i].Cascade, f.Name)
	}
	for _, f := range failed {
		if rootOf(f.Name) != f.Name {
			continue
		}
		if _, ok := idx[f.Name]; ok {
			continue
		}
		idx[f.Name] = len(out)
		out = append(out, Root{Failure: f, Diagnosis: diag[f.Name]})
	}
	for i := range out {
		sort.Strings(out[i].Cascade)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if len(out[i].Cascade) != len(out[j].Cascade) {
			return len(out[i].Cascade) > len(out[j].Cascade)
		}
		return out[i].Name < out[j].Name
	})
	return out
}
