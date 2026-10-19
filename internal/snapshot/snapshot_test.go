package snapshot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFromPluginListTxt(t *testing.T) {
	p := filepath.Join(t.TempDir(), "plugins.txt")
	os.WriteFile(p, []byte("# comment\ngit:5.2.1\nworkflow-cps 3800.v1  # trailing\n\n"), 0o644)
	s, err := FromPluginList(p, "2.462.3")
	if err != nil {
		t.Fatal(err)
	}
	if s.Core != "2.462.3" || len(s.Plugins) != 2 || s.Plugins[0].Name != "git" || s.Plugins[1].Version != "3800.v1" {
		t.Errorf("got %+v", s)
	}
}

func TestFromPluginListRejectsLatest(t *testing.T) {
	p := filepath.Join(t.TempDir(), "plugins.txt")
	os.WriteFile(p, []byte("git:latest\n"), 0o644)
	if _, err := FromPluginList(p, "2.462.3"); err == nil {
		t.Error("a plugins.txt without exact versions cannot describe a controller")
	}
}

func TestFromPluginListJSON(t *testing.T) {
	p := filepath.Join(t.TempDir(), "pm.json")
	os.WriteFile(p, []byte(`{"plugins":[{"shortName":"git","version":"5.2.1","enabled":true,"active":true,"dependencies":[{"shortName":"scm-api","version":"676.v1","optional":false}]}]}`), 0o644)
	s, err := FromPluginList(p, "2.462.3")
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Plugins) != 1 || len(s.Plugins[0].Dependencies) != 1 || s.Plugins[0].Dependencies[0].Name != "scm-api" {
		t.Errorf("got %+v", s.Plugins)
	}
}

func fakeController(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("the collector sent %s %s", r.Method, r.URL)
		}
		w.Header().Set("X-Jenkins", "2.462.3")
		if u, p, _ := r.BasicAuth(); u != "audit" || p != "tok" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/api/json":
			w.Write([]byte(`{"nodeName":""}`))
		case "/pluginManager/api/json":
			if !strings.Contains(r.URL.RawQuery, "dependencies") {
				t.Errorf("plugin tree does not ask for dependencies: %s", r.URL.RawQuery)
			}
			w.Write([]byte(`{"plugins":[{"shortName":"git","version":"5.2.1","enabled":true,"active":true,"dependencies":[]}]}`))
		case "/computer/api/json":
			w.Write([]byte(`{"computer":[
				{"_class":"hudson.model.Hudson$MasterComputer","displayName":"Built-In Node","offline":false,"monitorData":{"hudson.node_monitors.ArchitectureMonitor":"Linux (amd64)","hudson.plugin.versioncolumn.JVMVersionMonitor":"17.0.9"}},
				{"_class":"hudson.slaves.SlaveComputer","displayName":"aix71","offline":true,"monitorData":{"hudson.node_monitors.ArchitectureMonitor":"AIX (ppc64)","hudson.node_monitors.DiskSpaceMonitor":{"size":1}}}
			]}`))
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestCollect(t *testing.T) {
	srv := fakeController(t)
	defer srv.Close()
	c := &Collector{URL: srv.URL, User: "audit", Token: "tok", Tool: "t"}
	s, err := c.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if s.Core != "2.462.3" || s.Java != "17.0.9" || s.OS != "Linux (amd64)" {
		t.Errorf("controller = %s %s %s", s.Core, s.Java, s.OS)
	}
	if len(s.Plugins) != 1 || len(s.Agents) != 1 || s.Agents[0].OS != "AIX (ppc64)" || !s.Agents[0].Offline {
		t.Errorf("got %+v", s)
	}
	if c.Requests() != 3 {
		t.Errorf("requests = %d, want 3", c.Requests())
	}
}

func TestCollectBadToken(t *testing.T) {
	srv := fakeController(t)
	defer srv.Close()
	c := &Collector{URL: srv.URL, User: "audit", Token: "wrong"}
	if _, err := c.Collect(context.Background()); err == nil || !strings.Contains(err.Error(), "refused") {
		t.Errorf("want a refused-token error, got %v", err)
	}
}

func TestCollectRefusesPlainHTTP(t *testing.T) {
	c := &Collector{URL: "http://jenkins.example.com", Token: "tok"}
	if _, err := c.Collect(context.Background()); err == nil || !strings.Contains(err.Error(), "clear text") {
		t.Errorf("want a clear-text refusal, got %v", err)
	}
}

func TestGetOnly(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "https://x/", nil)
	if _, err := (getOnly{http.DefaultTransport}).RoundTrip(req); err == nil {
		t.Error("POST must be refused")
	}
}
