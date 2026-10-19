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
	// Unknown is true for a plugin the update center has never published (a
	// private build, or one removed long ago). It is kept as it is.
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
// is now too old. A plugin already newer than anything the core supports is
// left alone.
func Resolve(rel *uc.Releases, current []Pick, core string) Resolution {
	res := Resolution{Core: core}
	picks := map[string]*Pick{}
	order := []string{}
	set := func(p Pick) {
		if _, ok := picks[p.Name]; !ok {
			order = append(order, p.Name)
		}
		cp := p
		picks[p.Name] = &cp
	}
	for _, cur := range current {
		p := Pick{Name: cur.Name, Version: cur.Version, Added: cur.Added, RequiredBy: cur.RequiredBy}
		if !rel.Known(cur.Name) {
			p.Unknown = true
			set(p)
			continue
		}
		latest := rel.LatestFor(cur.Name, core)
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
				cand := rel.LatestFor(d.Name, core)
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
			cand := rel.LatestFor(d.Name, core)
			if cand == nil || version.Less(cand.Version, d.Version) {
				res.Problems = append(res.Problems, fmt.Sprintf("%s %s needs %s %s or newer, but the newest %s that runs on Jenkins %s is %s", name, p.Version, d.Name, d.Version, d.Name, core, have.Version))
				continue
			}
			have.Version = cand.Version
			queue = append(queue, d.Name)
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
		}
		out[i] = p
	}
	return out
}
