package plan

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/keelinfra/jenkins-upgrade-check/internal/snapshot"
	"github.com/keelinfra/jenkins-upgrade-check/internal/uc"
	"github.com/keelinfra/jenkins-upgrade-check/internal/version"
)

// Plan is the whole answer for one controller and one target.
type Plan struct {
	Tool        string    `json:"tool"`
	ToolVersion string    `json:"toolVersion"`
	GeneratedAt time.Time `json:"generatedAt"`
	Source      string    `json:"source"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	// JavaFrom and JavaTo are the Java versions each core runs on.
	JavaFrom       []int  `json:"javaFrom"`
	JavaTo         []int  `json:"javaTo"`
	ControllerJava string `json:"controllerJava,omitempty"`
	// Steps are done in order. The first updates plugins on the current
	// core; each later one is a Java switch (if needed), a core upgrade and
	// the plugin updates for that core.
	Steps []Step `json:"steps"`
	// Target is the plugin set for the final core; Dry-run starts it.
	Target       Resolution   `json:"target"`
	Plugins      []PluginPlan `json:"plugins"`
	CoreWarnings []WarningRef `json:"coreWarnings,omitempty"`
	Agents       []AgentNote  `json:"agents,omitempty"`
	Notes        []string     `json:"notes,omitempty"`
	Summary      Summary      `json:"summary"`
}

// Step is one stop on the way.
type Step struct {
	Core string `json:"core"`
	// Java is the Java version the controller and agents must run before
	// this step's core starts; 0 when the step changes no core.
	Java     int      `json:"java,omitempty"`
	Upgrade  bool     `json:"upgrade"`
	Changes  []Change `json:"changes,omitempty"`
	Problems []string `json:"problems,omitempty"`
	// Plugins is the whole plugin set once this step is done.
	Plugins []Pick `json:"plugins,omitempty"`
}

// Change is one plugin moving between versions; From is empty for a plugin
// the step adds.
type Change struct {
	Name string `json:"name"`
	From string `json:"from,omitempty"`
	To   string `json:"to"`
}

// PluginPlan is one plugin, start to finish.
type PluginPlan struct {
	Name            string       `json:"name"`
	Current         string       `json:"current,omitempty"`
	Target          string       `json:"target"`
	Added           bool         `json:"added,omitempty"`
	RequiredBy      []string     `json:"requiredBy,omitempty"`
	Disabled        bool         `json:"disabled,omitempty"`
	Unknown         bool         `json:"unknown,omitempty"`
	NotDistributed  bool         `json:"notDistributed,omitempty"`
	Deprecated      string       `json:"deprecated,omitempty"`
	VulnerableNow   []WarningRef `json:"vulnerableNow,omitempty"`
	VulnerableAfter []WarningRef `json:"vulnerableAfter,omitempty"`
}

// Attention reports whether the plugin needs a human decision.
func (p PluginPlan) Attention() bool {
	return p.Unknown || p.NotDistributed || p.Deprecated != "" || len(p.VulnerableNow) > 0 || len(p.VulnerableAfter) > 0
}

// WarningRef points at a security advisory.
type WarningRef struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	URL     string `json:"url"`
}

// AgentNote is an agent that cannot follow the controller as it is.
type AgentNote struct {
	Name    string `json:"name"`
	OS      string `json:"os,omitempty"`
	Java    string `json:"java,omitempty"`
	Offline bool   `json:"offline,omitempty"`
	Problem string `json:"problem"`
}

// Summary counts what the report lists.
type Summary struct {
	Plugins         int `json:"plugins"`
	Upgraded        int `json:"upgraded"`
	Added           int `json:"added"`
	VulnerableNow   int `json:"vulnerableNow"`
	VulnerableAfter int `json:"vulnerableAfter"`
	Deprecated      int `json:"deprecated"`
	NotDistributed  int `json:"notDistributed"`
	Unknown         int `json:"unknown"`
	Problems        int `json:"problems"`
	CoreUpgrades    int `json:"coreUpgrades"`
	Agents          int `json:"agentsAtRisk"`
}

// Build makes the plan.
func Build(s *snapshot.Snapshot, to string, rel *uc.Releases, center *uc.UpdateCenter) (*Plan, error) {
	if s.Core == "" {
		return nil, fmt.Errorf("the snapshot does not say which Jenkins version it is; pass --core")
	}
	if !version.Less(s.Core, to) {
		return nil, fmt.Errorf("the controller already runs %s, which is not older than %s", s.Core, to)
	}
	p := &Plan{
		GeneratedAt:    time.Now().UTC(),
		Source:         s.Source,
		From:           s.Core,
		To:             to,
		JavaFrom:       version.SupportedJava(s.Core),
		JavaTo:         version.SupportedJava(to),
		ControllerJava: s.Java,
		Notes:          append([]string(nil), s.Notes...),
	}

	current := make([]Pick, 0, len(s.Plugins))
	for _, pl := range s.Plugins {
		current = append(current, Pick{Name: pl.Name, Version: pl.Version})
	}

	// Step 1: bring plugins up to date on the core that runs today, as the
	// Jenkins upgrade guides ask.
	first := Resolve(rel, current, s.Core)
	p.Steps = append(p.Steps, Step{Core: s.Core, Changes: diff(current, first.Plugins), Problems: first.Problems, Plugins: first.Plugins})
	prev := first.Plugins
	last := first
	for _, hop := range version.Path(s.Core, to) {
		r := Resolve(rel, prev, hop.Core)
		p.Steps = append(p.Steps, Step{Core: hop.Core, Java: hop.Java, Upgrade: true, Changes: diff(prev, r.Plugins), Problems: r.Problems, Plugins: r.Plugins})
		prev = r.Plugins
		last = r
	}
	p.Target = last

	installed := s.Installed()
	for _, t := range last.Plugins {
		pp := PluginPlan{Name: t.Name, Target: t.Version, Added: t.Added, RequiredBy: t.RequiredBy, Unknown: t.Unknown}
		if cur, ok := installed[t.Name]; ok {
			pp.Current = cur.Version
			pp.Disabled = !cur.Enabled
			pp.Added = false
			pp.RequiredBy = nil
		}
		if center != nil && !t.Unknown {
			pp.NotDistributed = !center.Distributed(t.Name)
			if u, ok := center.Deprecated(t.Name); ok {
				pp.Deprecated = u
				if pp.Deprecated == "" {
					pp.Deprecated = "deprecated"
				}
			}
			if pp.Current != "" {
				pp.VulnerableNow = refs(center.WarningsFor("plugin", t.Name, pp.Current))
			}
			pp.VulnerableAfter = refs(center.WarningsFor("plugin", t.Name, t.Version))
		}
		p.Plugins = append(p.Plugins, pp)
	}
	if center != nil {
		p.CoreWarnings = refs(center.WarningsFor("core", "core", s.Core))
	}

	need := p.JavaTo
	if s.Java != "" {
		if major := JavaMajor(s.Java); major > 0 && !version.SupportsJava(to, major) {
			p.Notes = append(p.Notes, fmt.Sprintf("the controller runs Java %s; Jenkins %s runs on Java %s", s.Java, to, joinInts(need)))
		}
	}
	for _, a := range s.Agents {
		if prob := agentProblem(a, to); prob != "" {
			p.Agents = append(p.Agents, AgentNote{Name: a.Name, OS: a.OS, Java: a.Java, Offline: a.Offline, Problem: prob})
		}
	}
	p.summarize()
	return p, nil
}

func diff(before, after []Pick) []Change {
	old := map[string]string{}
	for _, b := range before {
		old[b.Name] = b.Version
	}
	var out []Change
	for _, a := range after {
		v, ok := old[a.Name]
		if !ok {
			out = append(out, Change{Name: a.Name, To: a.Version})
			continue
		}
		if v != a.Version {
			out = append(out, Change{Name: a.Name, From: v, To: a.Version})
		}
	}
	return out
}

func refs(ws []uc.Warning) []WarningRef {
	var out []WarningRef
	for _, w := range ws {
		out = append(out, WarningRef{ID: w.ID, Message: w.Message, URL: w.URL})
	}
	return out
}

var javaMajor = regexp.MustCompile(`^(?:1\.)?(\d+)`)

// JavaMajor turns "1.8.0_392", "17.0.9" or "21" into 8, 17, 21; 0 if unknown.
func JavaMajor(v string) int {
	m := javaMajor.FindStringSubmatch(strings.TrimSpace(v))
	if m == nil {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

func joinInts(xs []int) string {
	s := make([]string, len(xs))
	for i, x := range xs {
		s[i] = strconv.Itoa(x)
	}
	return strings.Join(s, " or ")
}

// agentProblem says why an agent cannot follow the controller, or "".
// Agents run the same remoting as the controller, so the Java requirement is
// the controller's (Jenkins Java support policy).
func agentProblem(a snapshot.Agent, to string) string {
	need := version.SupportedJava(to)
	min := 0
	if len(need) > 0 {
		min = need[0]
	}
	osName := strings.ToLower(a.OS)
	switch {
	case strings.Contains(osName, "sunos") || strings.Contains(osName, "solaris"):
		return fmt.Sprintf("Solaris has no Java %d or newer (Java 11 was the last for Solaris), so this agent cannot connect to Jenkins %s as it is. Agent8 runs the agent on the Java it has: https://keelinfra.io/jenkins/", min, to)
	case strings.Contains(osName, "hp-ux"):
		return fmt.Sprintf("HP-UX has no Java %d (HP's Java stops at 8), so this agent cannot connect to Jenkins %s as it is. Agent8 runs the agent on Java 8: https://keelinfra.io/jenkins/", min, to)
	case strings.Contains(osName, "aix"):
		if j := JavaMajor(a.Java); j > 0 && j >= min {
			return ""
		}
		return fmt.Sprintf("Jenkins %s needs Java %d on this agent; IBM Semeru %d needs AIX 7.2 TL5 or 7.3 on POWER8 or newer. On older AIX, Agent8 runs the agent on Java 8: https://keelinfra.io/jenkins/", to, min, min)
	}
	if j := JavaMajor(a.Java); j > 0 && !version.SupportsJava(to, j) {
		return fmt.Sprintf("runs Java %s; Jenkins %s needs Java %s on every agent", a.Java, to, joinInts(need))
	}
	return ""
}

func (p *Plan) summarize() {
	s := Summary{Plugins: len(p.Plugins), Agents: len(p.Agents)}
	for _, pp := range p.Plugins {
		switch {
		case pp.Added:
			s.Added++
		case pp.Current != "" && pp.Current != pp.Target:
			s.Upgraded++
		}
		if len(pp.VulnerableNow) > 0 {
			s.VulnerableNow++
		}
		if len(pp.VulnerableAfter) > 0 {
			s.VulnerableAfter++
		}
		if pp.Deprecated != "" {
			s.Deprecated++
		}
		if pp.NotDistributed {
			s.NotDistributed++
		}
		if pp.Unknown {
			s.Unknown++
		}
	}
	for _, st := range p.Steps {
		s.Problems += len(st.Problems)
		if st.Upgrade {
			s.CoreUpgrades++
		}
	}
	p.Summary = s
}
