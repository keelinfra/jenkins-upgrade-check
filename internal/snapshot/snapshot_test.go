package snapshot

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
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
			w.Write([]byte(`{"quietingDown":true,"jobs":[{"fullName":"app"},{"fullName":"team","jobs":[{"fullName":"team/deploy"}]}]}`))
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
	if !s.QuietingDown || strings.Join(s.Items, ",") != "app,team,team/deploy" {
		t.Errorf("quietingDown = %v, items = %q", s.QuietingDown, s.Items)
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

// A controller without the versioncolumn plugin: the Java of an online
// agent comes from its System Information page. The markup is as 2.580.1
// writes it (values hidden behind a button) and as 2.346.3 does (plain).
func TestCollectAgentJavaFromSystemInfo(t *testing.T) {
	var asked []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Jenkins", "2.346.3")
		asked = append(asked, r.URL.Path)
		switch r.URL.Path {
		case "/api/json":
			w.Write([]byte(`{"jobs":[]}`))
		case "/pluginManager/api/json":
			w.Write([]byte(`{"plugins":[]}`))
		case "/computer/api/json":
			w.Write([]byte(`{"computer":[
				{"_class":"hudson.model.Hudson$MasterComputer","displayName":"Built-In Node","offline":false,"monitorData":{}},
				{"_class":"hudson.slaves.SlaveComputer","displayName":"linux 1","offline":false,"monitorData":{}},
				{"_class":"hudson.slaves.SlaveComputer","displayName":"gone","offline":true,"monitorData":{}}
			]}`))
		case "/computer/linux 1/systemInfo":
			w.Write([]byte(`<table class="jenkins-table sortable" id="t1"><tbody>
<tr>
  <td>
    file.encoding
  </td>
  <td style="white-space: normal" class="property-table__value"><div>UTF-8</div></td>
</tr>
<tr><td>java.version</td><td style="white-space: normal" class="property-table__value"><div class="app-hidden-info-reveal"><button class="jenkins-button jenkins-button--transparent jenkins-!-color-blue"><svg xmlns="http://www.w3.org/2000/svg" aria-hidden="true" viewBox="0 0 512 512"><title/><circle cx="256" cy="256" r="80"/></svg> Hidden value, click to show this value
                  </button></div><div class="app-hidden-info-hide jenkins-hidden"><button class="jenkins-button jenkins-button--transparent"><span class="jenkins-!-color-blue"><svg xmlns="http://www.w3.org/2000/svg" aria-hidden="true" viewBox="0 0 512 512"><title/><path d="M432 448"/></svg> 17<wbr>.0<wbr>.12</span></button></div></td></tr>
<tr><td>os.name</td><td style="white-space: normal">Linux</td></tr>
</tbody></table>`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	c := &Collector{URL: srv.URL, Tool: "t"}
	s, err := c.Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Agents) != 2 || s.Agents[0].Java != "17.0.12" || s.Agents[0].OS != "Linux" {
		t.Errorf("agents = %+v", s.Agents)
	}
	if s.Agents[1].Java != "" {
		t.Errorf("an offline agent was read: %+v", s.Agents[1])
	}
	if len(s.Notes) != 0 {
		t.Errorf("notes = %q", s.Notes)
	}
	for _, p := range asked {
		if strings.Contains(p, "gone") {
			t.Errorf("asked the offline agent: %s", p)
		}
	}
}

func TestCollectAgentJavaUnknown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Jenkins", "2.346.3")
		switch r.URL.Path {
		case "/api/json":
			w.Write([]byte(`{"jobs":[]}`))
		case "/pluginManager/api/json":
			w.Write([]byte(`{"plugins":[]}`))
		case "/computer/api/json":
			w.Write([]byte(`{"computer":[{"_class":"hudson.slaves.SlaveComputer","displayName":"a1","offline":false,"monitorData":{}}]}`))
		default:
			w.WriteHeader(http.StatusForbidden)
		}
	}))
	defer srv.Close()
	s, err := (&Collector{URL: srv.URL}).Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Notes) != 1 || !strings.Contains(s.Notes[0], "Overall/Administer") || !strings.Contains(s.Notes[0], "a1") {
		t.Errorf("notes = %q", s.Notes)
	}
}

// minimalController answers the three reads of a collection.
func minimalController(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Jenkins", "2.462.3")
		switch r.URL.Path {
		case "/api/json":
			w.Write([]byte(`{"jobs":[]}`))
		case "/pluginManager/api/json":
			w.Write([]byte(`{"plugins":[{"shortName":"git","version":"5.2.1","enabled":true,"active":true}]}`))
		case "/computer/api/json":
			w.Write([]byte(`{"computer":[]}`))
		default:
			http.NotFound(w, r)
		}
	})
}

// A controller behind an internal CA: --ca-file makes it trusted.
func TestCollectWithCAFile(t *testing.T) {
	srv := httptest.NewTLSServer(minimalController(t))
	defer srv.Close()
	dir := t.TempDir()
	ca := filepath.Join(dir, "ca.pem")
	if err := os.WriteFile(ca, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw}), 0o644); err != nil {
		t.Fatal(err)
	}

	// Without the CA the certificate is not trusted.
	if _, err := (&Collector{URL: srv.URL}).Collect(context.Background()); err == nil || !strings.Contains(err.Error(), "could not reach the controller: tls: failed to verify certificate") {
		t.Errorf("err = %v", err)
	}
	s, err := (&Collector{URL: srv.URL, CAFile: ca, Timeout: 5 * time.Second}).Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if s.Core != "2.462.3" || len(s.Plugins) != 1 {
		t.Errorf("snapshot = %+v", s)
	}

	if _, err := (&Collector{URL: srv.URL, CAFile: filepath.Join(dir, "missing.pem")}).Collect(context.Background()); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a missing CA file: %v", err)
	}
	notPEM := filepath.Join(dir, "not.pem")
	os.WriteFile(notPEM, []byte("hello"), 0o644)
	if _, err := (&Collector{URL: srv.URL, CAFile: notPEM}).Collect(context.Background()); err == nil || err.Error() != notPEM+" holds no PEM certificate" {
		t.Errorf("a file without a certificate: %v", err)
	}
}

func TestCollectBadURL(t *testing.T) {
	for _, u := range []string{"jenkins.example.com", "://x"} {
		_, err := (&Collector{URL: u}).Collect(context.Background())
		if want := fmt.Sprintf("--url %q is not a URL", u); err == nil || err.Error() != want {
			t.Errorf("Collect(%q) = %v, want %s", u, err, want)
		}
	}
	if _, err := (&Collector{URL: "jenkins.example.com"}).AgentLog(context.Background(), "a"); err == nil {
		t.Error("AgentLog with a bad URL must fail")
	}
}

// Plain http is allowed to localhost and with --allow-http.
func TestCollectPlainHTTP(t *testing.T) {
	srv := httptest.NewServer(minimalController(t))
	defer srv.Close()
	_, port, _ := net.SplitHostPort(srv.Listener.Addr().String())
	if _, err := (&Collector{URL: "http://localhost:" + port}).Collect(context.Background()); err != nil {
		t.Errorf("localhost: %v", err)
	}
	if !isLoopback("::1") || isLoopback("jenkins.example.com") || isLoopback("10.0.0.1") {
		t.Error("isLoopback")
	}
	c := &Collector{URL: "http://jenkins.example.com", AllowHTTP: true}
	if err := c.init(); err != nil {
		t.Errorf("--allow-http: %v", err)
	}
}

// The token is never carried to another host, and redirect loops end.
func TestCollectRedirects(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Jenkins", "2.462.3")
		switch r.URL.Path {
		case "/elsewhere/api/json":
			// Same port, but "localhost" is another host than 127.0.0.1.
			http.Redirect(w, r, strings.Replace(srv.URL, "127.0.0.1", "localhost", 1)+"/api/json", http.StatusFound)
		case "/loop/api/json":
			http.Redirect(w, r, "/loop/api/json", http.StatusFound)
		case "/moved/api/json":
			http.Redirect(w, r, "/api/json", http.StatusMovedPermanently)
		default:
			minimalController(t).ServeHTTP(w, r)
		}
	}))
	defer srv.Close()
	ctx := context.Background()
	if _, err := (&Collector{URL: srv.URL + "/elsewhere", Token: "t"}).Collect(ctx); err == nil || !strings.Contains(err.Error(), "refusing redirect to localhost:") {
		t.Errorf("redirect to another host: %v", err)
	}
	if _, err := (&Collector{URL: srv.URL + "/loop"}).Collect(ctx); err == nil || !strings.Contains(err.Error(), "too many redirects") {
		t.Errorf("redirect loop: %v", err)
	}
	// A redirect on the same host is followed; the plugin list then lives
	// under the old prefix, which this controller does not serve.
	if _, err := (&Collector{URL: srv.URL + "/moved"}).Collect(ctx); err == nil || !strings.Contains(err.Error(), "reading the plugin list needs Overall/SystemRead or Overall/Administer: GET /pluginManager/api/json") || !strings.Contains(err.Error(), "404 Not Found") {
		t.Errorf("same-host redirect: %v", err)
	}
}

// Each way a collection can fail says what went wrong in words.
func TestCollectErrors(t *testing.T) {
	ctx := context.Background()
	gone := httptest.NewServer(http.NotFoundHandler())
	goneURL := gone.URL
	gone.Close()
	if _, err := (&Collector{URL: goneURL}).Collect(ctx); err == nil || !strings.HasPrefix(err.Error(), "could not reach the controller: dial tcp") || strings.Contains(err.Error(), "tree=") {
		t.Errorf("unreachable: %v", err)
	}

	for _, c := range []struct {
		name    string
		handler http.HandlerFunc
		want    string
	}{
		{"not Jenkins", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("{}")) },
			"the response carries no X-Jenkins header; is --url a Jenkins controller?"},
		{"server error", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Jenkins", "2.462.3")
			w.WriteHeader(http.StatusInternalServerError)
		}, "the controller answered with an error: GET /api/json?tree=quietingDown%2Cjobs%5BfullName%2Cjobs%5BfullName%2Cjobs%5BfullName%2Cjobs%5BfullName%5D%5D%5D%5D: 500 Internal Server Error"},
		{"bad JSON", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Jenkins", "2.462.3")
			w.Write([]byte("<html>"))
		}, "the controller answered with an error: GET /api/json?tree=quietingDown%2Cjobs%5BfullName%2Cjobs%5BfullName%2Cjobs%5BfullName%2Cjobs%5BfullName%5D%5D%5D%5D: invalid character '<' looking for beginning of value"},
		{"no plugin permission", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Jenkins", "2.462.3")
			if r.URL.Path == "/pluginManager/api/json" {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			w.Write([]byte("{}"))
		}, "reading the plugin list needs Overall/SystemRead or Overall/Administer: refused: GET /pluginManager/api/json?tree=plugins%5BshortName%2Cversion%2Cenabled%2Cactive%2Cdependencies%5BshortName%2Cversion%2Coptional%5D%5D: 403 Forbidden"},
	} {
		srv := httptest.NewServer(c.handler)
		_, err := (&Collector{URL: srv.URL}).Collect(ctx)
		srv.Close()
		if err == nil || err.Error() != c.want {
			t.Errorf("%s: err = %v\nwant %s", c.name, err, c.want)
		}
	}
}

// Agents that cannot be read are a note, not a failure.
func TestCollectAgentsUnreadable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/computer/api/json" {
			w.Header().Set("X-Jenkins", "2.462.3")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		minimalController(t).ServeHTTP(w, r)
	}))
	defer srv.Close()
	s, err := (&Collector{URL: srv.URL}).Collect(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Notes) != 1 || !strings.HasPrefix(s.Notes[0], "agents could not be read (refused: GET /computer/api/json") || !strings.HasSuffix(s.Notes[0], "403 Forbidden); their Java readiness is not checked") {
		t.Errorf("notes = %q", s.Notes)
	}
}

// The agent's connection log: the end of it, read with the same client.
func TestAgentLog(t *testing.T) {
	long := strings.Repeat("x", 300<<10) + "\nSEVERE: the end"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ua := r.Header.Get("User-Agent"); ua != "jenkins-upgrade-check/1.0" {
			t.Errorf("User-Agent = %q", ua)
		}
		switch r.URL.Path {
		case "/computer/linux 1/logText/progressiveText":
			if r.URL.Query().Get("start") != "0" {
				t.Errorf("query = %s", r.URL.RawQuery)
			}
			w.Write([]byte(long))
		case "/computer/cut/logText/progressiveText":
			w.Header().Set("Content-Length", "100")
			w.Write([]byte("partial"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	c := &Collector{URL: srv.URL, Tool: "jenkins-upgrade-check", ToolVersion: "1.0"}
	got, err := c.AgentLog(context.Background(), "linux 1")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 256<<10 || !strings.HasSuffix(got, "\nSEVERE: the end") {
		t.Errorf("log: %d bytes, ends %q", len(got), got[len(got)-20:])
	}
	if _, err := c.AgentLog(context.Background(), "cut"); err == nil || err.Error() != "GET /computer/cut/logText/progressiveText?start=0: unexpected EOF" {
		t.Errorf("a cut-off log: %v", err)
	}
	if _, err := c.AgentLog(context.Background(), "nobody"); err == nil || !strings.Contains(err.Error(), "404 Not Found") {
		t.Errorf("an unknown agent: %v", err)
	}
	if c.Requests() != 3 {
		t.Errorf("requests = %d", c.Requests())
	}
}

func TestMonitorString(t *testing.T) {
	for raw, want := range map[string]string{
		``:                              "",
		`null`:                          "",
		`"Linux (amd64)"`:               "Linux (amd64)",
		`{"version":"17.0.9"}`:          "17.0.9",
		`{"value":"21"}`:                "21",
		`{"displayName":"AIX (ppc64)"}`: "AIX (ppc64)",
		`{"size":1}`:                    "",
		`42`:                            "",
	} {
		if got := monitorString(json.RawMessage(raw)); got != want {
			t.Errorf("monitorString(%s) = %q, want %q", raw, got, want)
		}
	}
}

// System properties come before environment variables of the same name.
func TestSystemPropertiesFirstWins(t *testing.T) {
	page := `<table><tr><td>java.version</td><td>17.0.12</td></tr></table>
<table><tr><td>java.version</td><td>from the environment</td></tr><tr><td>os.name</td><td>Linux &amp; more</td></tr></table>`
	got := systemProperties(page)
	if got["java.version"] != "17.0.12" || got["os.name"] != "Linux & more" {
		t.Errorf("properties = %q", got)
	}
}

func TestWriteAndRead(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "snapshot.json")
	s := &Snapshot{Tool: "jenkins-upgrade-check", Core: "2.462.3", CollectedAt: time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC),
		Plugins: []Plugin{{Name: "git", Version: "5.2.1", Enabled: true, Active: true}, {Name: "ant", Version: "1.0"}},
		Agents:  []Agent{{Name: "a1", Java: "17.0.9"}}}
	if err := s.Write(p); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(p)
	if !strings.HasSuffix(string(b), "}\n") || !strings.Contains(string(b), "\n  \"core\": \"2.462.3\",") {
		t.Errorf("not indented JSON with a final newline:\n%s", b)
	}
	got, err := Read(p)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, s) {
		t.Errorf("read back\n%+v\nwant\n%+v", got, s)
	}
	inst := got.Installed()
	if len(inst) != 2 || inst["git"].Version != "5.2.1" || inst["ant"].Enabled {
		t.Errorf("Installed = %+v", inst)
	}

	if _, err := Read(filepath.Join(dir, "missing.json")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a missing file: %v", err)
	}
	bad := filepath.Join(dir, "bad.json")
	os.WriteFile(bad, []byte("plugins.txt?"), 0o644)
	if _, err := Read(bad); err == nil || !strings.HasPrefix(err.Error(), bad+": invalid character") {
		t.Errorf("not JSON: %v", err)
	}
	if err := s.Write(filepath.Join(dir, "no", "such", "dir.json")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("an unwritable path: %v", err)
	}
	// A time JSON cannot hold is refused, not written half.
	s.CollectedAt = time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := s.Write(filepath.Join(dir, "y10k.json")); err == nil {
		t.Error("a year JSON cannot hold must be refused")
	}
	if _, err := os.Stat(filepath.Join(dir, "y10k.json")); err == nil {
		t.Error("a refused snapshot must not be written")
	}
}

func TestFromPluginListErrors(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		os.WriteFile(p, []byte(body), 0o644)
		return p
	}
	if _, err := FromPluginList(filepath.Join(dir, "missing.txt"), "2.462.3"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("missing: %v", err)
	}
	p := write("pm.json", `{"plugins":"no"}`)
	if _, err := FromPluginList(p, "2.462.3"); err == nil || !strings.HasPrefix(err.Error(), p+": json: cannot unmarshal") {
		t.Errorf("bad JSON: %v", err)
	}
	p = write("three.txt", "git 5.2.1 extra\n")
	if _, err := FromPluginList(p, "2.462.3"); err == nil || err.Error() != p+`: line 1: want name:version, got "git 5.2.1 extra"` {
		t.Errorf("three words: %v", err)
	}
	p = write("empty.txt", "# nothing here\n\n")
	if _, err := FromPluginList(p, "2.462.3"); err == nil || err.Error() != p+" lists no plugins" {
		t.Errorf("empty: %v", err)
	}
	p = write("empty.json", `{"plugins":[]}`)
	if _, err := FromPluginList(p, "2.462.3"); err == nil || err.Error() != p+" lists no plugins" {
		t.Errorf("empty JSON: %v", err)
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("disk full") }

func TestWritePluginsTxt(t *testing.T) {
	var b strings.Builder
	if err := WritePluginsTxt(&b, []Plugin{{Name: "git", Version: "5.2.1"}, {Name: "ant", Version: "1.0"}}); err != nil {
		t.Fatal(err)
	}
	if b.String() != "git:5.2.1\nant:1.0\n" {
		t.Errorf("plugins.txt = %q", b.String())
	}
	if err := WritePluginsTxt(failingWriter{}, []Plugin{{Name: "git", Version: "5.2.1"}}); err == nil || err.Error() != "disk full" {
		t.Errorf("err = %v", err)
	}
}
