// Package plan works out what an upgrade involves: the path through Java
// generations, the plugin versions for every stop on the way, and what the
// update center already knows is wrong with them.
package plan

import (
	"fmt"
	"sort"

	"github.com/keelinfra/jenkins-upgrade-check/internal/uc"
	"github.com/keelinfra/jenkins-upgrade-check/internal/version"
)

// Pick is one plugin at one version.
type Pick struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	// URL and SHA256 locate the release on the update center; empty for a
	// plugin the update center never published.
	URL    string `json:"url,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
	// Added is true for a plugin that is not installed today but becomes a
	// required dependency.
	Added      bool     `json:"added,omitempty"`
	RequiredBy []string `json:"requiredBy,omitempty"`
	// Unknown is true for a plugin the update center does not list: a
	// private build, or one removed from distribution. It is kept as it is,
	// and URL points at the download archive, which still has every public
	// release (but no checksum).
	Unknown bool `json:"unknown,omitempty"`
}

// Resolution is a consistent plugin set for one core.
type Resolution struct {
	Core     string   `json:"core"`
	Plugins  []Pick   `json:"plugins"`
	Problems []string `json:"problems,omitempty"`
}

// ByName indexes the plugins.
func (r Resolution) ByName() map[string]Pick {
	m := make(map[string]Pick, len(r.Plugins))
	for _, p := range r.Plugins {
		m[p.Name] = p
	}
	return m
}

// Resolve picks, for each plugin, the newest release the core can load, then
// adds whatever those releases newly depend on and raises any dependency that
// is now too old. A release whose dependencies need a newer core is passed
// over for an older one, as the update center for that core does. A plugin
// already newer than anything the core supports is left alone.
func Resolve(rel *uc.Releases, current []Pick, core string) Resolution {
	res := Resolution{Core: core}
	l := &loadable{rel: rel, core: core, memo: map[string]*uc.Release{}, busy: map[string]bool{}}
	picks := map[string]*Pick{}
	order := []string{}
	set := func(p Pick) {
		if _, ok := picks[p.Name]; !ok {
			order = append(order, p.Name)
		}
		cp := p
		picks[p.Name] = &cp
	}
	var anyVersion []string // unknown plugins asked for by name only
	for _, cur := range current {
		p := Pick{Name: cur.Name, Version: cur.Version, Added: cur.Added, RequiredBy: cur.RequiredBy}
		if !rel.Known(cur.Name) && (cur.Version == "" || cur.Version == "0") {
			// No version to take it at: only a plugin that depends on it
			// can say which.
			anyVersion = append(anyVersion, cur.Name)
			continue
		}
		if !rel.Known(cur.Name) {
			p.Unknown = true
			p.URL = uc.ArchiveURL(cur.Name, cur.Version)
			set(p)
			continue
		}
		latest := l.newest(cur.Name, "")
		switch {
		case latest == nil:
			res.Problems = append(res.Problems, fmt.Sprintf("%s: no release runs on Jenkins %s; kept at %s", cur.Name, core, cur.Version))
		case version.Less(cur.Version, latest.Version):
			p.Version = latest.Version
		}
		set(p)
	}

	queue := append([]string(nil), order...)
	seen := 0
	for len(queue) > 0 && seen < 100000 {
		seen++
		name := queue[0]
		queue = queue[1:]
		p := picks[name]
		r := rel.Get(name, p.Version)
		if r == nil {
			continue
		}
		for _, d := range r.Dependencies {
			have, installed := picks[d.Name]
			if !installed {
				if d.Optional {
					continue
				}
				if !rel.Known(d.Name) {
					// Removed from the update center since; an installation
					// of that age has it, so take the version asked for.
					set(Pick{Name: d.Name, Version: d.Version, Added: true, RequiredBy: []string{name}, Unknown: true, URL: uc.ArchiveURL(d.Name, d.Version)})
					res.Problems = append(res.Problems, fmt.Sprintf("%s %s needs %s %s, which the update center no longer distributes", name, p.Version, d.Name, d.Version))
					continue
				}
				cand := l.newest(d.Name, d.Version)
				if cand == nil || version.Less(cand.Version, d.Version) {
					res.Problems = append(res.Problems, fmt.Sprintf("%s %s needs %s %s or newer, and no such release runs on Jenkins %s", name, p.Version, d.Name, d.Version, core))
					continue
				}
				set(Pick{Name: d.Name, Version: cand.Version, Added: true, RequiredBy: []string{name}})
				queue = append(queue, d.Name)
				continue
			}
			if have.Added && !contains(have.RequiredBy, name) {
				have.RequiredBy = append(have.RequiredBy, name)
			}
			if !version.Less(have.Version, d.Version) {
				continue
			}
			if have.Unknown {
				// Not in the update center any more: take the version asked
				// for from the download archive.
				have.Version, have.URL = d.Version, uc.ArchiveURL(d.Name, d.Version)
				continue
			}
			cand := l.newest(d.Name, d.Version)
			if cand == nil || version.Less(cand.Version, d.Version) {
				res.Problems = append(res.Problems, fmt.Sprintf("%s %s needs %s %s or newer, but the newest %s that runs on Jenkins %s is %s", name, p.Version, d.Name, d.Version, d.Name, core, have.Version))
				continue
			}
			have.Version = cand.Version
			queue = append(queue, d.Name)
		}
		// Plugins built for an older core still expect what was split out
		// of core since; Jenkins treats that as a dependency.
		for _, d := range implied(r, core) {
			if !rel.Known(d.Name) {
				continue
			}
			if have, installed := picks[d.Name]; installed {
				if have.Added && !contains(have.RequiredBy, name) {
					have.RequiredBy = append(have.RequiredBy, name)
				}
				continue
			}
			cand := l.newest(d.Name, d.Version)
			if cand == nil {
				res.Problems = append(res.Problems, fmt.Sprintf("%s %s was built for Jenkins %s and expects %s, which was split out of core, but no %s release runs on Jenkins %s", name, p.Version, r.RequiredCore, d.Name, d.Name, core))
				continue
			}
			set(Pick{Name: d.Name, Version: cand.Version, Added: true, RequiredBy: []string{name}})
			queue = append(queue, d.Name)
		}
	}

	for _, name := range anyVersion {
		if _, ok := picks[name]; !ok {
			res.Problems = append(res.Problems, fmt.Sprintf("%s: the update center no longer distributes it and no version was given; left out", name))
		}
	}
	for _, name := range order {
		p := picks[name]
		if r := rel.Get(name, p.Version); r != nil {
			p.URL, p.SHA256 = r.URL, r.SHA256
		}
		sort.Strings(p.RequiredBy)
		res.Plugins = append(res.Plugins, *p)
	}
	sort.Slice(res.Plugins, func(i, j int) bool { return res.Plugins[i].Name < res.Plugins[j].Name })
	sort.Strings(res.Problems)
	res.Problems = dedupe(res.Problems)
	return res
}

// loadable finds, for each plugin, the newest release the core can load
// together with everything it requires.
type loadable struct {
	rel  *uc.Releases
	core string
	memo map[string]*uc.Release
	busy map[string]bool
}

// newest is the newest release of the plugin, at least atLeast, whose
// required dependencies all run on the core; failing that, the newest release
// the core can load, so that the missing dependency shows up as a problem.
func (l *loadable) newest(name, atLeast string) *uc.Release {
	if b := l.best(name); b != nil && (atLeast == "" || !version.Less(b.Version, atLeast)) {
		return b
	}
	return l.rel.LatestFor(name, l.core)
}

func (l *loadable) best(name string) *uc.Release {
	if r, ok := l.memo[name]; ok {
		return r
	}
	if l.busy[name] {
		// A dependency cycle: take the newest and let the rest decide.
		return l.rel.LatestFor(name, l.core)
	}
	l.busy[name] = true
	defer delete(l.busy, name)
	var found *uc.Release
	for _, r := range l.rel.LoadableFor(name, l.core) {
		if l.requirementsMet(r) {
			found = r
			break
		}
	}
	l.memo[name] = found
	return found
}

func (l *loadable) requirementsMet(r *uc.Release) bool {
	for _, d := range r.Dependencies {
		if d.Optional || !l.rel.Known(d.Name) {
			continue
		}
		b := l.best(d.Name)
		if b == nil || version.Less(b.Version, d.Version) {
			return false
		}
	}
	return true
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

func dedupe(list []string) []string {
	out := list[:0]
	for i, s := range list {
		if i == 0 || s != list[i-1] {
			out = append(out, s)
		}
	}
	return out
}

// Locate fills URL and SHA256 for picks at the versions given, without
// changing any version. It is how the current plugin set is downloaded to
// rehearse an upgrade that leaves plugins alone.
func Locate(rel *uc.Releases, picks []Pick) []Pick {
	out := make([]Pick, len(picks))
	for i, p := range picks {
		if r := rel.Get(p.Name, p.Version); r != nil {
			p.URL, p.SHA256 = r.URL, r.SHA256
		} else if !rel.Known(p.Name) {
			p.Unknown = true
			p.URL = uc.ArchiveURL(p.Name, p.Version)
		}
		out[i] = p
	}
	return out
}
