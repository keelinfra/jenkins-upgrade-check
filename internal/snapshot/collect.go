package snapshot

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
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
	requests int
}

type getOnly struct{ next http.RoundTripper }

func (g getOnly) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Method != http.MethodGet {
		return nil, fmt.Errorf("refusing %s %s: this tool only reads", r.Method, r.URL)
	}
	return g.next.RoundTrip(r)
}

func (c *Collector) init() error {
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
	c.requests++
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<16))
		err := fmt.Errorf("GET %s: %s", path, resp.Status)
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			err = fmt.Errorf("%w: %w", ErrRefused, err)
		}
		return resp.Header, err
	}
	if into != nil {
		if err := json.NewDecoder(io.LimitReader(resp.Body, 256<<20)).Decode(into); err != nil {
			return resp.Header, fmt.Errorf("GET %s: %w", path, err)
		}
	}
	return resp.Header, nil
}

// ErrRefused marks a request the controller refused (401 or 403): the user
// or the API token is wrong, or lacks a permission.
var ErrRefused = errors.New("refused")

// Requests is how many GET requests the collection sent.
func (c *Collector) Requests() int { return c.requests }

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
	if !sawJava && len(nodes.Computer) > 0 {
		s.Notes = append(s.Notes, "Java versions are not reported: Jenkins core does not publish them; install the Versions Node Monitors plugin (versioncolumn) to have them read")
	}
	return s, nil
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
