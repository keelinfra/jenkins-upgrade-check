package snapshot

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Collector reads a live controller. It only ever sends GET requests: the
// transport refuses anything else, so a bug here cannot change the controller.
type Collector struct {
	URL               string
	User              string
	Token             string
	CAFile            string
	AllowHTTP         bool
	Timeout           time.Duration
	Tool, ToolVersion string

	client   *http.Client
	requests atomic.Int64
}

type getOnly struct{ next http.RoundTripper }

func (g getOnly) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Method != http.MethodGet {
		return nil, fmt.Errorf("refusing %s %s: this tool only reads", r.Method, r.URL)
	}
	return g.next.RoundTrip(r)
}

func (c *Collector) init() error {
	if c.client != nil {
		return nil
	}
	u, err := url.Parse(c.URL)
	if err != nil || u.Host == "" {
		return fmt.Errorf("--url %q is not a URL", c.URL)
	}
	if u.Scheme != "https" && !c.AllowHTTP && !isLoopback(u.Hostname()) {
		return fmt.Errorf("--url %s would send the API token in clear text; use https, or pass --allow-http", c.URL)
	}
	tr := http.DefaultTransport.(*http.Transport).Clone()
	if c.CAFile != "" {
		pem, err := os.ReadFile(c.CAFile)
		if err != nil {
			return err
		}
		pool, err := x509.SystemCertPool()
		if err != nil || pool == nil {
			pool = x509.NewCertPool()
		}
		if !pool.AppendCertsFromPEM(pem) {
			return fmt.Errorf("%s holds no PEM certificate", c.CAFile)
		}
		tr.TLSClientConfig = &tls.Config{RootCAs: pool, MinVersion: tls.VersionTLS12}
	}
	timeout := c.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	c.client = &http.Client{
		Transport: getOnly{tr},
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Never carry the token to another host.
			if req.URL.Host != via[0].URL.Host {
				return fmt.Errorf("refusing redirect to %s", req.URL.Host)
			}
			if len(via) >= 5 {
				return errors.New("too many redirects")
			}
			return nil
		},
	}
	return nil
}

func isLoopback(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func (c *Collector) get(ctx context.Context, path string, into any) (http.Header, error) {
	resp, err := c.do(ctx, path)
	if err != nil {
		if resp != nil {
			return resp.Header, err
		}
		return nil, err
	}
	defer resp.Body.Close()
	if into != nil {
		if err := json.NewDecoder(io.LimitReader(resp.Body, 256<<20)).Decode(into); err != nil {
			return resp.Header, fmt.Errorf("GET %s: %w", path, err)
		}
	}
	return resp.Header, nil
}

// getText reads a page as text, at most its last limit bytes.
func (c *Collector) getText(ctx context.Context, path string, limit int64) (string, error) {
	resp, err := c.do(ctx, path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
	if err != nil {
		return "", fmt.Errorf("GET %s: %w", path, err)
	}
	if int64(len(b)) > limit {
		b = b[int64(len(b))-limit:]
	}
	return string(b), nil
}

// do sends one GET. On a status other than 200 the body is closed and the
// response comes back with the error, for its headers.
func (c *Collector) do(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(c.URL, "/")+path, nil)
	if err != nil {
		return nil, err
	}
	if c.User != "" || c.Token != "" {
		req.SetBasicAuth(c.User, c.Token)
	}
	ua := c.Tool
	if c.ToolVersion != "" {
		ua += "/" + c.ToolVersion
	}
	req.Header.Set("User-Agent", ua)
	c.requests.Add(1)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<16))
		resp.Body.Close()
		err := fmt.Errorf("GET %s: %s", path, resp.Status)
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			err = fmt.Errorf("%w: %w", ErrRefused, err)
		}
		return resp, err
	}
	return resp, nil
}

// ErrRefused marks a request the controller refused (401 or 403): the user
// or the API token is wrong, or lacks a permission.
var ErrRefused = errors.New("refused")

// Requests is how many GET requests the collection sent.
func (c *Collector) Requests() int { return int(c.requests.Load()) }

const (
	pluginTree   = "plugins[shortName,version,enabled,active,dependencies[shortName,version,optional]]"
	computerTree = "computer[_class,displayName,offline,monitorData[*]]"

	architectureMonitor = "hudson.node_monitors.ArchitectureMonitor"
	jvmVersionMonitor   = "hudson.plugin.versioncolumn.JVMVersionMonitor"
	builtInClass        = "hudson.model.Hudson$MasterComputer"

	// itemTree reads job and folder names four levels deep, which covers
	// folders of folders; deeper items are not counted.
	itemTree = "quietingDown,jobs[fullName,jobs[fullName,jobs[fullName,jobs[fullName]]]]"
)

type item struct {
	FullName string `json:"fullName"`
	Jobs     []item `json:"jobs"`
}

func flatten(items []item, out []string) []string {
	for _, it := range items {
		if it.FullName != "" {
			out = append(out, it.FullName)
		}
		out = flatten(it.Jobs, out)
	}
	return out
}

// Collect reads the core version, the plugins and the agents.
func (c *Collector) Collect(ctx context.Context) (*Snapshot, error) {
	if err := c.init(); err != nil {
		return nil, err
	}
	s := &Snapshot{Tool: c.Tool, ToolVersion: c.ToolVersion, CollectedAt: time.Now().UTC(), Source: strings.TrimRight(c.URL, "/")}

	// The core version arrives as a header on every response, a refusal
	// included.
	var top struct {
		QuietingDown bool   `json:"quietingDown"`
		Jobs         []item `json:"jobs"`
	}
	h, err := c.get(ctx, "/api/json?tree="+url.QueryEscape(itemTree), &top)
	if h != nil {
		s.Core = h.Get("X-Jenkins")
	}
	if s.Core == "" {
		if err != nil {
			// Without the request URL, whose long tree query says nothing
			// to the reader: "connection refused" is the news.
			var ue *url.Error
			if errors.As(err, &ue) {
				err = ue.Err
			}
			return nil, fmt.Errorf("could not reach the controller: %w", err)
		}
		return nil, errors.New("the response carries no X-Jenkins header; is --url a Jenkins controller?")
	}
	if errors.Is(err, ErrRefused) {
		return nil, fmt.Errorf("the controller answered but refused the API token: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("the controller answered with an error: %w", err)
	}
	s.QuietingDown = top.QuietingDown
	s.Items = flatten(top.Jobs, nil)
	sort.Strings(s.Items)

	var pm pluginManager
	if _, err := c.get(ctx, "/pluginManager/api/json?tree="+url.QueryEscape(pluginTree), &pm); err != nil {
		return nil, fmt.Errorf("reading the plugin list needs Overall/SystemRead or Overall/Administer: %w", err)
	}
	s.Plugins = pm.plugins()
	s.sortPlugins()

	var nodes struct {
		Computer []struct {
			Class       string                     `json:"_class"`
			DisplayName string                     `json:"displayName"`
			Offline     bool                       `json:"offline"`
			MonitorData map[string]json.RawMessage `json:"monitorData"`
		} `json:"computer"`
	}
	if _, err := c.get(ctx, "/computer/api/json?tree="+url.QueryEscape(computerTree), &nodes); err != nil {
		s.Notes = append(s.Notes, fmt.Sprintf("agents could not be read (%v); their Java readiness is not checked", err))
	}
	sawJava := false
	for _, n := range nodes.Computer {
		osName := monitorString(n.MonitorData[architectureMonitor])
		java := monitorString(n.MonitorData[jvmVersionMonitor])
		if java != "" {
			sawJava = true
		}
		if n.Class == builtInClass {
			s.OS, s.Java = osName, java
			continue
		}
		s.Agents = append(s.Agents, Agent{Name: n.DisplayName, OS: osName, Java: java, Offline: n.Offline})
	}
	if !sawJava && len(s.Agents) > 0 {
		if note := c.agentJava(ctx, s.Agents); note != "" {
			s.Notes = append(s.Notes, note)
		}
	}
	return s, nil
}

// agentJava fills in the Java of online agents from their System
// Information page, since core reports no Java version: one GET per agent,
// a few at a time. It returns a note when some could not be read.
func (c *Collector) agentJava(ctx context.Context, agents []Agent) string {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		refused bool
		missed  []string
		sem     = make(chan struct{}, 8)
	)
	for i := range agents {
		a := &agents[i]
		if a.Offline {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			page, err := c.getText(ctx, "/computer/"+url.PathEscape(a.Name)+"/systemInfo", 8<<20)
			props := systemProperties(page)
			mu.Lock()
			defer mu.Unlock()
			if err != nil || props["java.version"] == "" {
				refused = refused || errors.Is(err, ErrRefused)
				missed = append(missed, a.Name)
				return
			}
			a.Java = props["java.version"]
			if a.OS == "" {
				a.OS = props["os.name"]
			}
		}()
	}
	wg.Wait()
	if len(missed) == 0 {
		return ""
	}
	sort.Strings(missed)
	why := "their System Information page did not say"
	if refused {
		why = "reading their System Information page needs Overall/Administer"
	}
	return fmt.Sprintf("the Java version of %d online agent(s) is not known (%s): %s; install the Versions Node Monitors plugin (versioncolumn) to have it reported", len(missed), why, strings.Join(missed, ", "))
}

var propertyRow = regexp.MustCompile(`(?s)<td[^>]*>\s*([A-Za-z0-9_.]+)\s*</td>\s*<td[^>]*>(.*?)</td>`)
var tags = regexp.MustCompile(`<[^>]*>`)

// systemProperties reads name and value pairs from the property tables of
// a System Information page. Values come with <wbr> and other markup, which
// is dropped; newer versions hide each value behind a button whose label
// comes first, so only what follows the hidden part's opening tag counts.
// The first table with a name wins: system properties come before
// environment variables.
func systemProperties(page string) map[string]string {
	out := map[string]string{}
	for _, m := range propertyRow.FindAllStringSubmatch(page, -1) {
		if _, seen := out[m[1]]; seen {
			continue
		}
		cell := m[2]
		if i := strings.Index(cell, "app-hidden-info-hide"); i >= 0 {
			cell = cell[i:]
			if j := strings.IndexByte(cell, '>'); j >= 0 {
				cell = cell[j+1:]
			}
		}
		out[m[1]] = strings.Join(strings.Fields(html.UnescapeString(tags.ReplaceAllString(cell, ""))), " ")
	}
	return out
}

// AgentLog reads the end of an agent's connection log on the controller.
func (c *Collector) AgentLog(ctx context.Context, name string) (string, error) {
	if err := c.init(); err != nil {
		return "", err
	}
	return c.getText(ctx, "/computer/"+url.PathEscape(name)+"/logText/progressiveText?start=0", 256<<10)
}

func monitorString(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var obj map[string]any
	if json.Unmarshal(raw, &obj) == nil {
		for _, k := range []string{"version", "value", "displayName"} {
			if v, ok := obj[k].(string); ok {
				return v
			}
		}
	}
	return ""
}
