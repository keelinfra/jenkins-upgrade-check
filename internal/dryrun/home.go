package dryrun

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// HomeCopy is what was taken from a JENKINS_HOME for the rehearsal.
type HomeCopy struct {
	Files int `json:"files"`
	// Items are the job and folder directories found, relative to
	// JENKINS_HOME (jobs/a, jobs/a/jobs/b).
	Items []string `json:"items"`
	// Skipped are files that could not be read, usually for permissions.
	Skipped []string `json:"skipped,omitempty"`
}

// copyHome copies the configuration of a JENKINS_HOME, and nothing else:
// no build records, workspaces, logs, caches or plugins. That is what Jenkins
// reads at startup, and it stays small even for a large controller.
func copyHome(src, dst string) (*HomeCopy, error) {
	st, err := os.Stat(src)
	if err != nil {
		return nil, err
	}
	if !st.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", src)
	}
	if _, err := os.Stat(filepath.Join(src, "config.xml")); err != nil {
		return nil, fmt.Errorf("%s has no config.xml; is it a JENKINS_HOME?", src)
	}
	hc := &HomeCopy{}
	copyOne := func(rel string) {
		if err := copyFile(filepath.Join(src, rel), filepath.Join(dst, rel)); err != nil {
			hc.Skipped = append(hc.Skipped, rel+": "+err.Error())
			return
		}
		hc.Files++
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		name := e.Name()
		switch {
		case e.Type().IsRegular() && rootFile(name):
			copyOne(name)
		case e.IsDir() && (name == "secrets" || name == "users" || name == "init.groovy.d" || name == "casc_configs"):
			walkFiles(src, name, func(rel string, d fs.DirEntry) { copyOne(rel) })
		case e.IsDir() && name == "nodes":
			walkFiles(src, name, func(rel string, d fs.DirEntry) {
				if d.Name() == "config.xml" {
					copyOne(rel)
				}
			})
		case e.IsDir() && name == "jobs":
			walkFiles(src, name, func(rel string, d fs.DirEntry) {
				if d.Name() != "config.xml" || inBuildRecords(rel) {
					return
				}
				copyOne(rel)
				hc.Items = append(hc.Items, filepath.ToSlash(filepath.Dir(rel)))
			})
		}
	}
	sort.Strings(hc.Items)
	if err := quietBuiltInNode(filepath.Join(dst, "config.xml")); err != nil {
		return nil, err
	}
	return hc, nil
}

// rootFile picks the files at the top of JENKINS_HOME that hold
// configuration: every *.xml (core and plugin settings, credentials), the
// keys that decrypt secrets, and Configuration as Code or Groovy files.
func rootFile(name string) bool {
	switch {
	case strings.HasSuffix(name, ".xml"):
		return !strings.HasPrefix(name, "queue") // pending builds are not configuration
	case strings.HasPrefix(name, "secret.key"), name == "identity.key.enc":
		return true
	case strings.HasSuffix(name, ".yaml"), strings.HasSuffix(name, ".yml"), strings.HasSuffix(name, ".groovy"):
		return true
	}
	return false
}

// inBuildRecords is true for files under a builds directory, at any depth.
func inBuildRecords(rel string) bool {
	for _, part := range strings.Split(filepath.ToSlash(rel), "/") {
		if part == "builds" || strings.HasPrefix(part, "workspace") {
			return true
		}
	}
	return false
}

// walkFiles calls fn for each regular file under root/dir, not following
// symbolic links (JENKINS_HOME has many, pointing at builds).
func walkFiles(root, dir string, fn func(rel string, d fs.DirEntry)) {
	_ = filepath.WalkDir(filepath.Join(root, dir), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if d.Name() == "builds" || strings.HasPrefix(d.Name(), "workspace") {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		fn(rel, d)
		return nil
	})
}

var numExecutors = regexp.MustCompile(`<numExecutors>\s*\d+\s*</numExecutors>`)

// quietBuiltInNode sets the built-in node's executors to zero so the
// rehearsal never runs a build, even one a timer starts.
func quietBuiltInNode(configXML string) error {
	b, err := os.ReadFile(configXML)
	if err != nil {
		return err
	}
	out := numExecutors.ReplaceAll(b, []byte("<numExecutors>0</numExecutors>"))
	return os.WriteFile(configXML, out, 0o644)
}

var (
	pluginAttr = regexp.MustCompile(`\splugin="([^"@]+)(?:@([^"]*))?"`)
	xmlTag     = regexp.MustCompile(`<[A-Za-z_][^<>]*>`)
	cannotRead = regexp.MustCompile(`(?:CannotResolveClassException|ClassNotFoundException|NoClassDefFoundError):\s*(\S+)`)
)

// pluginOf finds the plugin that a config.xml credits with a class: XStream
// writes plugin="name@version" on the element that names the class, either
// as the element itself or in its class attribute. With an empty class it
// returns the plugin of the root element.
func pluginOf(configXML []byte, class string) (name, ver string) {
	tags := xmlTag.FindAll(configXML, -1)
	for _, t := range tags {
		s := string(t)
		if strings.HasPrefix(s, "<?") || strings.HasPrefix(s, "<!") {
			continue
		}
		if class != "" {
			el := strings.TrimPrefix(strings.Fields(strings.TrimSuffix(strings.TrimSuffix(s, ">"), "/"))[0], "<")
			if el != xstreamName(class) && !strings.Contains(s, `class="`+class+`"`) {
				continue
			}
		}
		if m := pluginAttr.FindStringSubmatch(s); m != nil {
			return m[1], m[2]
		}
		if class == "" {
			return "", "" // the root element names no plugin
		}
	}
	return "", ""
}

// xstreamName is how XStream writes a class name as an element name.
func xstreamName(class string) string {
	return strings.ReplaceAll(strings.ReplaceAll(class, "_", "__"), "$", "_-")
}

// ConfigProblem is a job, folder or setting that did not come back whole.
type ConfigProblem struct {
	Item    string `json:"item"`
	Path    string `json:"path,omitempty"`
	Problem string `json:"problem"`
	// Plugin is the plugin the configuration says the missing part came
	// from, when it says.
	Plugin string `json:"plugin,omitempty"`
	// PluginState says what happened to that plugin in the rehearsal.
	PluginState string `json:"pluginState,omitempty"`
}

// Plugin states for ConfigProblem.
const (
	PluginNotInstalled = "not installed"
	PluginFailed       = "failed to load"
	PluginLoaded       = "loaded, but no longer provides it"
)

// blame fills in the plugin a problem comes from, reading the item's
// config.xml from the rehearsal copy.
func blame(p *ConfigProblem, ref string, failed, active map[string]bool) {
	if p.Path == "" {
		return
	}
	b, err := os.ReadFile(filepath.Join(ref, filepath.FromSlash(p.Path), "config.xml"))
	if err != nil {
		return
	}
	class := ""
	if m := cannotRead.FindStringSubmatch(p.Problem); m != nil {
		class = strings.ReplaceAll(m[1], "/", ".")
	}
	name, ver := pluginOf(b, class)
	if name == "" && class != "" {
		// The class was the root element.
		name, ver = pluginOf(b, "")
	}
	if name == "" {
		return
	}
	p.Plugin = name
	if ver != "" {
		p.Plugin += " " + ver
	}
	switch {
	case failed[name]:
		p.PluginState = PluginFailed
	case active[name]:
		p.PluginState = PluginLoaded
	default:
		p.PluginState = PluginNotInstalled
	}
}
