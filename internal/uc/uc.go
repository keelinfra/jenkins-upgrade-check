// Package uc reads the Jenkins update center: which plugins are distributed,
// every release of each with the core it needs and what it depends on, and the
// security warnings and deprecations Jenkins shows administrators.
//
// Two files are enough for everything this tool says:
//
//   - current/update-center.actual.json: the plugins still distributed,
//     security warnings and deprecations.
//   - current/plugin-versions.json: every release of every plugin, with its
//     required core and dependencies. Tiered update sites only exist for about
//     a year of cores, so a controller older than that cannot ask the update
//     center what its own core can run; this file answers for any core.
package uc

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/keelinfra/jenkins-upgrade-check/internal/version"
)

// Dependency is one plugin another plugin needs.
type Dependency struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Optional bool   `json:"optional"`
}

// Release is one published version of a plugin.
type Release struct {
	Name             string       `json:"name"`
	Version          string       `json:"version"`
	RequiredCore     string       `json:"requiredCore"`
	URL              string       `json:"url"`
	SHA256           string       `json:"sha256"`
	ReleaseTimestamp string       `json:"releaseTimestamp"`
	Dependencies     []Dependency `json:"dependencies"`

	v version.Version
}

// Warning is a security warning as the update center publishes it.
type Warning struct {
	ID       string         `json:"id"`
	Message  string         `json:"message"`
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	URL      string         `json:"url"`
	Versions []WarningRange `json:"versions"`
}

// WarningRange is one affected range. Jenkins matches Pattern against the
// whole version string; LastVersion is the last affected release.
type WarningRange struct {
	FirstVersion string `json:"firstVersion"`
	LastVersion  string `json:"lastVersion"`
	Pattern      string `json:"pattern"`

	re *regexp.Regexp
}

// Deprecation marks a plugin its maintainers or the project consider
// deprecated.
type Deprecation struct {
	URL string `json:"url"`
}

// Plugin is a plugin's entry in update-center.actual.json.
type Plugin struct {
	Name         string       `json:"name"`
	Title        string       `json:"title"`
	Version      string       `json:"version"`
	RequiredCore string       `json:"requiredCore"`
	Dependencies []Dependency `json:"dependencies"`
}

// UpdateCenter is update-center.actual.json.
type UpdateCenter struct {
	Core struct {
		Version string `json:"version"`
	} `json:"core"`
	Plugins      map[string]Plugin      `json:"plugins"`
	Warnings     []Warning              `json:"warnings"`
	Deprecations map[string]Deprecation `json:"deprecations"`
}

// ReadUpdateCenter decodes update-center.actual.json, or update-center.json
// with its JSONP wrapper.
func ReadUpdateCenter(r io.Reader) (*UpdateCenter, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	b = stripJSONP(b)
	var u UpdateCenter
	if err := json.Unmarshal(b, &u); err != nil {
		return nil, fmt.Errorf("update center JSON: %w", err)
	}
	for i := range u.Warnings {
		for j := range u.Warnings[i].Versions {
			wr := &u.Warnings[i].Versions[j]
			if wr.Pattern != "" {
				// A pattern Go cannot compile falls back to lastVersion.
				wr.re, _ = regexp.Compile("^(?:" + wr.Pattern + ")$")
			}
		}
	}
	return &u, nil
}

func stripJSONP(b []byte) []byte {
	s := strings.TrimSpace(string(b))
	if strings.HasPrefix(s, "updateCenter.post(") {
		s = strings.TrimPrefix(s, "updateCenter.post(")
		s = strings.TrimSuffix(strings.TrimSpace(s), ";")
		s = strings.TrimSuffix(strings.TrimSpace(s), ")")
		return []byte(s)
	}
	return b
}

// Distributed reports whether the update center still offers the plugin.
func (u *UpdateCenter) Distributed(name string) bool {
	_, ok := u.Plugins[name]
	return ok
}

// Deprecated returns the deprecation notice URL, if the plugin is deprecated.
func (u *UpdateCenter) Deprecated(name string) (string, bool) {
	d, ok := u.Deprecations[name]
	return d.URL, ok
}

// WarningsFor returns the warnings of the given type ("plugin" or "core")
// that apply to name at version v.
func (u *UpdateCenter) WarningsFor(kind, name, v string) []Warning {
	var out []Warning
	for _, w := range u.Warnings {
		if w.Type != kind || w.Name != name {
			continue
		}
		if w.Affects(v) {
			out = append(out, w)
		}
	}
	return out
}

// Affects reports whether the warning applies to version v.
func (w Warning) Affects(v string) bool {
	for _, r := range w.Versions {
		if r.re != nil {
			if r.re.MatchString(v) {
				return true
			}
			continue
		}
		if r.LastVersion == "" {
			continue
		}
		if r.FirstVersion != "" && version.Less(v, r.FirstVersion) {
			continue
		}
		if !version.Less(r.LastVersion, v) {
			return true
		}
	}
	return false
}

// pluginVersionsFile is plugin-versions.json.
type pluginVersionsFile struct {
	Plugins map[string]map[string]*Release `json:"plugins"`
}

// Releases holds every release of every plugin.
type Releases struct {
	// byName lists each plugin's releases, oldest first.
	byName map[string][]*Release
}

// ReadReleases decodes plugin-versions.json.
func ReadReleases(r io.Reader) (*Releases, error) {
	var f pluginVersionsFile
	if err := json.NewDecoder(r).Decode(&f); err != nil {
		return nil, fmt.Errorf("plugin versions JSON: %w", err)
	}
	return newReleases(f.Plugins), nil
}

func newReleases(m map[string]map[string]*Release) *Releases {
	db := &Releases{byName: make(map[string][]*Release, len(m))}
	for name, rels := range m {
		list := make([]*Release, 0, len(rels))
		for v, r := range rels {
			if r == nil {
				continue
			}
			if r.Version == "" {
				r.Version = v
			}
			if r.Name == "" {
				r.Name = name
			}
			r.v = version.Parse(r.Version)
			list = append(list, r)
		}
		sort.Slice(list, func(i, j int) bool { return list[i].v.Compare(list[j].v) < 0 })
		db.byName[name] = list
	}
	return db
}

// Known reports whether any release of the plugin was ever published.
func (db *Releases) Known(name string) bool {
	return len(db.byName[name]) > 0
}

// Get returns one release, or nil.
func (db *Releases) Get(name, v string) *Release {
	for _, r := range db.byName[name] {
		if r.Version == v {
			return r
		}
	}
	return nil
}

// LatestFor returns the newest release of the plugin that the given core can
// load, skipping pre-releases, or nil if there is none.
func (db *Releases) LatestFor(name, core string) *Release {
	list := db.byName[name]
	for i := len(list) - 1; i >= 0; i-- {
		r := list[i]
		if version.IsPrerelease(r.Version) {
			continue
		}
		if r.RequiredCore == "" || version.AtLeast(core, r.RequiredCore) {
			return r
		}
	}
	return nil
}

// LoadableFor lists the releases of the plugin that the given core can load,
// newest first, skipping pre-releases.
func (db *Releases) LoadableFor(name, core string) []*Release {
	list := db.byName[name]
	var out []*Release
	for i := len(list) - 1; i >= 0; i-- {
		r := list[i]
		if version.IsPrerelease(r.Version) {
			continue
		}
		if r.RequiredCore == "" || version.AtLeast(core, r.RequiredCore) {
			out = append(out, r)
		}
	}
	return out
}

// Latest returns the newest non-pre-release of the plugin, whatever core it
// needs.
func (db *Releases) Latest(name string) *Release {
	list := db.byName[name]
	for i := len(list) - 1; i >= 0; i-- {
		if !version.IsPrerelease(list[i].Version) {
			return list[i]
		}
	}
	return nil
}

// ArchiveURL is where the Jenkins download server keeps every release ever
// published, including plugins the update center no longer lists. There is
// no checksum for these.
func ArchiveURL(name, v string) string {
	return fmt.Sprintf("%s/download/plugins/%s/%s/%s.hpi", DefaultBase, name, v, name)
}
