// Package snapshot is what the tool knows about one controller: its core,
// its plugins and its agents. It holds no secret, so it can be attached to a
// support request as it is.
package snapshot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/keelinfra/jenkins-upgrade-check/internal/uc"
)

// Snapshot is one controller.
type Snapshot struct {
	Tool        string    `json:"tool"`
	ToolVersion string    `json:"toolVersion"`
	CollectedAt time.Time `json:"collectedAt"`
	// Source is the controller URL, or the file the plugin list came from.
	Source string `json:"source"`
	Core   string `json:"core"`
	// Java is the controller's Java version, when the Versions Node Monitors
	// plugin reports it.
	Java    string   `json:"java,omitempty"`
	OS      string   `json:"os,omitempty"`
	Plugins []Plugin `json:"plugins"`
	Agents  []Agent  `json:"agents,omitempty"`
	// Items are the full names of jobs and folders, four levels deep.
	Items []string `json:"items,omitempty"`
	// QuietingDown is true while the controller starts no new builds.
	QuietingDown bool `json:"quietingDown,omitempty"`
	// Notes are things the collection could not read, and why.
	Notes []string `json:"notes,omitempty"`
}

// Plugin is one installed plugin.
type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Enabled bool   `json:"enabled"`
	Active  bool   `json:"active"`
	// Dependencies are what the installed version declares.
	Dependencies []uc.Dependency `json:"dependencies,omitempty"`
}

// Agent is one node other than the built-in one.
type Agent struct {
	Name    string `json:"name"`
	OS      string `json:"os,omitempty"`
	Java    string `json:"java,omitempty"`
	Offline bool   `json:"offline"`
}

// Read loads a snapshot written by Write.
func Read(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &s, nil
}

// Write saves the snapshot as indented JSON.
func (s *Snapshot) Write(path string) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

// Installed returns the plugins by name.
func (s *Snapshot) Installed() map[string]Plugin {
	m := make(map[string]Plugin, len(s.Plugins))
	for _, p := range s.Plugins {
		m[p.Name] = p
	}
	return m
}

func (s *Snapshot) sortPlugins() {
	sort.Slice(s.Plugins, func(i, j int) bool { return s.Plugins[i].Name < s.Plugins[j].Name })
}

// FromPluginList builds a snapshot from a plugin list. It accepts:
//
//   - plugins.txt, one "name:version" per line, as jenkins-plugin-cli and the
//     official Docker image use (and as the script console one-liner in the
//     README prints);
//   - the JSON of GET /pluginManager/api/json?depth=1, saved with curl.
//
// The core version is not in either, so the caller supplies it.
func FromPluginList(path, core string) (*Snapshot, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := &Snapshot{Source: path, Core: core}
	trimmed := strings.TrimSpace(string(b))
	if strings.HasPrefix(trimmed, "{") {
		var pm pluginManager
		if err := json.Unmarshal(b, &pm); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		s.Plugins = pm.plugins()
	} else {
		plugins, err := parsePluginsTxt(strings.NewReader(string(b)))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		s.Plugins = plugins
	}
	if len(s.Plugins) == 0 {
		return nil, fmt.Errorf("%s lists no plugins", path)
	}
	s.sortPlugins()
	return s, nil
}

func parsePluginsTxt(r io.Reader) ([]Plugin, error) {
	var out []Plugin
	sc := bufio.NewScanner(r)
	n := 0
	for sc.Scan() {
		n++
		line := sc.Text()
		if i := strings.Index(line, "#"); i >= 0 {
			line = line[:i]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		name, ver, ok := strings.Cut(line, ":")
		if !ok {
			// "name version" works too.
			f := strings.Fields(line)
			if len(f) != 2 {
				return nil, fmt.Errorf("line %d: want name:version, got %q", n, line)
			}
			name, ver = f[0], f[1]
		}
		name, ver = strings.TrimSpace(name), strings.TrimSpace(ver)
		if name == "" || ver == "" || ver == "latest" || ver == "experimental" {
			return nil, fmt.Errorf("line %d: %q needs an exact version", n, line)
		}
		out = append(out, Plugin{Name: name, Version: ver, Enabled: true, Active: true})
	}
	return out, sc.Err()
}

// WritePluginsTxt writes "name:version" lines.
func WritePluginsTxt(w io.Writer, plugins []Plugin) error {
	for _, p := range plugins {
		if _, err := fmt.Fprintf(w, "%s:%s\n", p.Name, p.Version); err != nil {
			return err
		}
	}
	return nil
}

// pluginManager is GET /pluginManager/api/json.
type pluginManager struct {
	Plugins []struct {
		ShortName    string `json:"shortName"`
		Version      string `json:"version"`
		Enabled      bool   `json:"enabled"`
		Active       bool   `json:"active"`
		Dependencies []struct {
			ShortName string `json:"shortName"`
			Version   string `json:"version"`
			Optional  bool   `json:"optional"`
		} `json:"dependencies"`
	} `json:"plugins"`
}

func (pm pluginManager) plugins() []Plugin {
	out := make([]Plugin, 0, len(pm.Plugins))
	for _, p := range pm.Plugins {
		pl := Plugin{Name: p.ShortName, Version: p.Version, Enabled: p.Enabled, Active: p.Active}
		for _, d := range p.Dependencies {
			pl.Dependencies = append(pl.Dependencies, uc.Dependency{Name: d.ShortName, Version: d.Version, Optional: d.Optional})
		}
		out = append(out, pl)
	}
	return out
}
