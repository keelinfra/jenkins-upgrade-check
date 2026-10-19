// Package dryrun rehearses an upgrade: it starts a throwaway Jenkins of the
// target version in Docker with the plugin set given, and optionally a copy
// of a real controller's configuration, and reports what failed and why.
//
// The container runs without a network, so the rehearsal cannot reach a
// repository, an agent, a deployment target or a mail server, and the
// built-in node gets no executors, so no build runs.
package dryrun

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
	"github.com/keelinfra/jenkins-upgrade-check/internal/version"
)

//go:embed probe.groovy
var probeScript []byte

const probeFile = "jenkins-upgrade-check-probe.json"

// Options configure one rehearsal.
type Options struct {
	Core    string
	Java    int // Java for the image; 0 picks the oldest the core and plugins allow
	Plugins []plan.Pick
	// Home is a JENKINS_HOME whose configuration (not its build records)
	// is copied into the rehearsal.
	Home string
	// CasC is a Configuration as Code file or directory to load, when the
	// controller reads it from outside JENKINS_HOME.
	CasC string
	// Image overrides the base image (default jenkins/jenkins:<core>-jdk<N>).
	Image    string
	OutDir   string // where the log is kept
	CacheDir string // downloaded plugins
	Timeout  time.Duration
	Keep     bool
	Docker   string
	Client   *http.Client
	Progress io.Writer
}

// Result is what the rehearsal saw.
type Result struct {
	Core           string  `json:"core"`
	Image          string  `json:"image"`
	Java           int     `json:"java,omitempty"`
	JavaReason     string  `json:"javaReason,omitempty"`
	Plugins        int     `json:"plugins"`
	Started        bool    `json:"started"`
	Reason         string  `json:"reason,omitempty"`
	StartupSeconds float64 `json:"startupSeconds,omitempty"`
	// Failed are plugins Jenkins did not load; Roots groups them by cause.
	Failed []Failure `json:"failed,omitempty"`
	Roots  []Root    `json:"roots,omitempty"`
	// Broken are plugins that loaded but some of whose parts did not.
	Broken []Failure `json:"broken,omitempty"`
	// Skipped are plugins that could not be downloaded, with why.
	Skipped []string `json:"skipped,omitempty"`
	Errors  []string `json:"errors,omitempty"`
	LogFile string   `json:"logFile,omitempty"`
	// Home is set when a real configuration was rehearsed.
	Home *HomeResult `json:"home,omitempty"`
}

// HomeResult is what happened to a real configuration.
type HomeResult struct {
	Source   string   `json:"source"`
	Files    int      `json:"files"`
	Items    int      `json:"items"`
	Loaded   int      `json:"loaded"`
	CopyErrs []string `json:"copyErrors,omitempty"`
	// Missing are jobs and folders that did not load.
	Missing []ConfigProblem `json:"missing,omitempty"`
	// Unreadable is configuration Jenkins could not read (Manage Jenkins →
	// Manage old data), usually settings of a plugin that failed or changed.
	Unreadable []ConfigProblem `json:"unreadable,omitempty"`
	// Monitors are warnings Jenkins shows administrators after startup.
	Monitors []string `json:"monitors,omitempty"`
}

// Failure is a plugin Jenkins did not load, or loaded only in part.
type Failure struct {
	Name    string   `json:"name"`
	Version string   `json:"version,omitempty"`
	Reasons []string `json:"reasons,omitempty"`
}

// OK reports whether Jenkins started and loaded everything.
func (r *Result) OK() bool {
	if !r.Started || len(r.Failed) > 0 || len(r.Broken) > 0 {
		return false
	}
	if r.Home != nil && (len(r.Home.Missing) > 0 || len(r.Home.Unreadable) > 0) {
		return false
	}
	return true
}

func (o *Options) logf(format string, args ...any) {
	if o.Progress != nil {
		fmt.Fprintf(o.Progress, format+"\n", args...)
	}
}

func (o *Options) docker() string {
	if o.Docker != "" {
		return o.Docker
	}
	return "docker"
}

// Run performs the rehearsal.
func Run(ctx context.Context, o Options) (*Result, error) {
	if o.Core == "" {
		return nil, errors.New("no Jenkins version to rehearse")
	}
	if _, err := exec.LookPath(o.docker()); err != nil {
		return nil, fmt.Errorf("dry-run needs Docker: %w", err)
	}
	if o.Timeout == 0 {
		o.Timeout = 10 * time.Minute
	}
	res := &Result{Core: o.Core}

	work, err := os.MkdirTemp("", "jenkins-upgrade-check-")
	if err != nil {
		return nil, err
	}
	if !o.Keep {
		defer os.RemoveAll(work)
	}
	ref := filepath.Join(work, "ref")
	pluginDir := filepath.Join(ref, "plugins")
	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		return nil, err
	}

	var expectedItems []string
	if o.Home != "" {
		o.logf("copying the configuration of %s (no build records)", o.Home)
		hc, err := copyHome(o.Home, ref)
		if err != nil {
			return nil, err
		}
		expectedItems = hc.Items
		res.Home = &HomeResult{Source: o.Home, Files: hc.Files, Items: len(hc.Items), CopyErrs: hc.Skipped}
	}
	casc := ""
	if o.CasC != "" {
		if err := copyTree(o.CasC, filepath.Join(ref, "casc_configs")); err != nil {
			return nil, fmt.Errorf("--casc: %w", err)
		}
		casc = "/var/jenkins_home/casc_configs"
	}
	if err := os.MkdirAll(filepath.Join(ref, "init.groovy.d"), 0o755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(ref, "init.groovy.d", "zz-jenkins-upgrade-check-probe.groovy"), probeScript, 0o644); err != nil {
		return nil, err
	}

	var fetch []plan.Pick
	for _, p := range o.Plugins {
		if p.URL == "" {
			res.Skipped = append(res.Skipped, fmt.Sprintf("%s %s: no download known", p.Name, p.Version))
			continue
		}
		fetch = append(fetch, p)
	}
	o.logf("downloading %d plugins", len(fetch))
	got, skipped := o.download(ctx, fetch, pluginDir)
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	res.Skipped = append(res.Skipped, skipped...)
	res.Plugins = len(got)

	minJava, who := minimumJava(pluginDir)
	base, java, err := o.pickImage(ctx, minJava)
	if err != nil {
		return nil, err
	}
	res.Image, res.Java = base, java
	if minJava > 0 && o.Java == 0 && o.Image == "" && java >= minJava {
		if oldest := version.SupportedJava(o.Core); len(oldest) > 0 && oldest[0] < minJava {
			res.JavaReason = fmt.Sprintf("%s needs Java %d or newer", who, minJava)
		}
	}

	dockerfile := "FROM " + base + "\n" +
		"ENV JAVA_OPTS=\"-Djenkins.install.runSetupWizard=false -Dhudson.model.UpdateCenter.never=true\"\n"
	if casc != "" {
		dockerfile += "ENV CASC_JENKINS_CONFIG=" + casc + "\n"
	}
	dockerfile += "COPY --chown=jenkins:jenkins ref/ /usr/share/jenkins/ref/\n"
	if err := os.WriteFile(filepath.Join(work, "Dockerfile"), []byte(dockerfile), 0o644); err != nil {
		return nil, err
	}
	tag := fmt.Sprintf("jenkins-upgrade-check:%s-%d", sanitize(o.Core), time.Now().UnixNano())
	o.logf("building %s", tag)
	if out, err := o.run(ctx, "build", "-q", "-t", tag, work); err != nil {
		return nil, fmt.Errorf("docker build: %v\n%s", err, out)
	}
	if !o.Keep {
		defer o.run(context.Background(), "rmi", "-f", tag)
	}

	name := "juc-" + strings.ReplaceAll(strings.TrimPrefix(tag, "jenkins-upgrade-check:"), ".", "-")
	if out, err := o.run(ctx, "run", "-d", "--name", name, "--network", "none", tag); err != nil {
		return nil, fmt.Errorf("docker run: %v\n%s", err, out)
	}
	if !o.Keep {
		defer o.run(context.Background(), "rm", "-f", name)
	}
	o.logf("starting Jenkins %s without network (container %s)", o.Core, name)

	start := time.Now()
	deadline := start.Add(o.Timeout)
	var logs string
	for {
		logs, _ = o.run(ctx, "logs", name)
		if strings.Contains(logs, "Jenkins is fully up and running") {
			res.Started = true
			res.StartupSeconds = time.Since(start).Seconds()
			break
		}
		running, _ := o.run(ctx, "inspect", "-f", "{{.State.Running}} {{.State.ExitCode}}", name)
		if strings.HasPrefix(strings.TrimSpace(running), "false") {
			res.Reason = "the container exited (state: " + strings.TrimSpace(running) + ")"
			break
		}
		if time.Now().After(deadline) {
			res.Reason = fmt.Sprintf("not up after %s", o.Timeout)
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	var probe *probeData
	if res.Started {
		// The probe writes its final report a few seconds after startup.
		for i := 0; i < 30 && probe == nil; i++ {
			if p, err := o.readProbe(ctx, name, work); err == nil && p.Completed {
				probe = p
				break
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Second):
			}
		}
		logs, _ = o.run(ctx, "logs", name)
	}
	if probe == nil {
		if p, err := o.readProbe(ctx, name, work); err == nil {
			probe = p
		}
	}

	versions := map[string]string{}
	for _, p := range o.Plugins {
		versions[p.Name] = p.Version
	}
	failed, errs := ParseLog(logs)
	if probe != nil {
		failed = mergeProbe(failed, probe, got)
		for _, e := range probe.ProbeErrors {
			errs = append(errs, "probe: "+e)
		}
	} else if res.Started {
		errs = append(errs, "the rehearsal did not report its plugin list; only the log was read")
	}
	byName := map[string]bool{}
	for i := range failed {
		byName[failed[i].Name] = true
		if failed[i].Version == "" {
			failed[i].Version = versions[failed[i].Name]
		}
	}
	res.Failed = failed
	res.Roots = Roots(failed)
	// Errors that name a class: find the plugin that ships it.
	res.Broken, res.Errors = attribute(errs, pluginDir, byName, versions)

	if res.Home != nil && probe != nil {
		active := map[string]bool{}
		for _, p := range probe.Plugins {
			active[p.Name] = p.Active
		}
		var missing []string
		res.Home.Loaded, missing = compareItems(expectedItems, probe.Items)
		itemErrs := map[string]string{}
		res.Errors, itemErrs = itemErrors(res.Errors)
		for _, m := range missing {
			p := ConfigProblem{Item: itemName(m), Path: m, Problem: shortClass(itemErrs[path.Base(m)])}
			if p.Problem == "" {
				p.Problem = "not loaded; the log has no reason"
			}
			blame(&p, ref, byName, active)
			res.Home.Missing = append(res.Home.Missing, p)
		}
		for _, u := range probe.Unreadable {
			u.Problem = shortClass(u.Problem)
			blame(&u, ref, byName, active)
			res.Home.Unreadable = append(res.Home.Unreadable, u)
		}
		for _, m := range probe.Monitors {
			if !isNoise(m.ID) {
				res.Home.Monitors = append(res.Home.Monitors, m.Name)
			}
		}
	}

	logDir := o.OutDir
	if logDir == "" {
		logDir = work
	}
	if err := os.MkdirAll(logDir, 0o755); err == nil {
		lf := filepath.Join(logDir, "dry-run-"+sanitize(o.Core)+".log")
		if os.WriteFile(lf, []byte(logs), 0o644) == nil && (o.OutDir != "" || o.Keep) {
			res.LogFile = lf
		}
	}
	return res, nil
}

// rehearsalNoise are administrator warnings a rehearsal shows because it has
// no network, no URL, no executors and a fresh Java, or that the report
// already covers; they say nothing new about the upgrade. Matched as
// prefixes, since some monitors add details to their ID.
var rehearsalNoise = []string{
	"OldData", // reported as settings Jenkins could not read
	"hudson.diagnosis.OldDataMonitor",
	"hudson.PluginWrapper$PluginWrapperAdministrativeMonitor", // reported as failed plugins
	"hudson.diagnosis.ReverseProxySetupMonitor",
	"jenkins.diagnostics.RootUrlNotSetMonitor",
	"jenkins.diagnostics.URICheckEncodingMonitor",
	"jenkins.diagnostics.ControllerExecutorsNoAgents",
	"jenkins.diagnostics.ControllerExecutorsAgents",
	"jenkins.diagnostics.SecurityIsOffMonitor",
	"jenkins.security.csrf.CSRFAdministrativeMonitor",
	"hudson.diagnosis.TooManyJobsButNoView",
	"jenkins.security.UpdateSiteWarningsMonitor",
	"hudson.model.UpdateCenter$CoreUpdateMonitor",
	"hudson.PluginManager$PluginUpdateMonitor",
	"jenkins.monitor.JavaVersionRecommendationAdminMonitor",
	"jenkins.security.ResourceDomainRecommendation",
	"jenkins.security.csp.impl.CspRecommendation",
}

func isNoise(id string) bool {
	for _, n := range rehearsalNoise {
		if id == n || strings.HasPrefix(id, n+"-") {
			return true
		}
	}
	return false
}

type probeData struct {
	Completed bool   `json:"completed"`
	Version   string `json:"version"`
	Java      string `json:"java"`
	Plugins   []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Active  bool   `json:"active"`
		Enabled bool   `json:"enabled"`
	} `json:"plugins"`
	Failed []struct {
		Name   string   `json:"name"`
		Causes []string `json:"causes"`
	} `json:"failed"`
	Items      []string        `json:"items"`
	Unreadable []ConfigProblem `json:"unreadable"`
	Monitors   []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"monitors"`
	ProbeErrors []string `json:"probeErrors"`
}

func (o *Options) readProbe(ctx context.Context, container, work string) (*probeData, error) {
	dest := filepath.Join(work, probeFile)
	os.Remove(dest)
	if out, err := o.run(ctx, "cp", container+":/var/jenkins_home/"+probeFile, dest); err != nil {
		return nil, fmt.Errorf("%v: %s", err, out)
	}
	b, err := os.ReadFile(dest)
	if err != nil {
		return nil, err
	}
	var p probeData
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// mergeProbe adds what the probe saw to what the log said: plugins Jenkins
// lists as failed, and plugins installed but not active.
func mergeProbe(failed []Failure, probe *probeData, installed []plan.Pick) []Failure {
	idx := map[string]int{}
	for i, f := range failed {
		idx[f.Name] = i
	}
	for _, f := range probe.Failed {
		reasons := probeReasons(f.Causes)
		if i, ok := idx[f.Name]; ok {
			if len(failed[i].Reasons) == 0 {
				failed[i].Reasons = reasons
			}
			continue
		}
		idx[f.Name] = len(failed)
		failed = append(failed, Failure{Name: f.Name, Reasons: reasons})
	}
	active := map[string]bool{}
	for _, p := range probe.Plugins {
		active[p.Name] = p.Active
	}
	for _, p := range installed {
		if active[p.Name] {
			continue
		}
		if _, ok := idx[p.Name]; !ok {
			idx[p.Name] = len(failed)
			failed = append(failed, Failure{Name: p.Name, Version: p.Version, Reasons: []string{reasonInactive}})
		}
	}
	sort.Slice(failed, func(i, j int) bool { return failed[i].Name < failed[j].Name })
	return failed
}

var failedItem = regexp.MustCompile(`Failed Loading item (\S+) — (.*)$`)

// itemErrors takes the "Failed Loading item" errors out of the log errors,
// by item name.
func itemErrors(errs []string) ([]string, map[string]string) {
	out := map[string]string{}
	var rest []string
	for _, e := range errs {
		if m := failedItem.FindStringSubmatch(e); m != nil {
			out[m[1]] = m[2]
			continue
		}
		rest = append(rest, e)
	}
	return rest, out
}

// shortClass drops the package of XStream's exception, which says nothing
// to an administrator.
func shortClass(s string) string {
	return strings.TrimPrefix(s, "com.thoughtworks.xstream.mapper.")
}

// itemName turns jobs/a/jobs/b into a/b, the name Jenkins shows.
func itemName(rel string) string {
	parts := strings.Split(rel, "/")
	var out []string
	for i, p := range parts {
		if i%2 == 1 {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return rel
	}
	return strings.Join(out, "/")
}

// compareItems counts the jobs and folders that came back after startup and
// lists the ones that did not.
func compareItems(expected, loaded []string) (int, []string) {
	got := map[string]bool{}
	for _, it := range loaded {
		got[filepath.ToSlash(it)] = true
	}
	n := 0
	var missing []string
	for _, it := range expected {
		if got[it] {
			n++
		} else {
			missing = append(missing, it)
		}
	}
	return n, missing
}

// probeReasons turns an exception chain into the reasons Jenkins lists: the
// " - " lines of a dependency failure, or the exceptions themselves.
func probeReasons(causes []string) []string {
	var out []string
	for _, c := range causes {
		for _, line := range strings.Split(c, "\n") {
			if t := logLine(line); strings.HasPrefix(t, "- ") {
				out = append(out, strings.TrimPrefix(t, "- "))
			}
		}
	}
	if len(out) > 0 {
		return out
	}
	for _, c := range causes {
		first, _, _ := strings.Cut(c, "\n")
		if strings.HasPrefix(first, "java.io.IOException: Failed to load") {
			continue
		}
		out = append(out, strings.TrimSpace(first))
	}
	return out
}

var failedToLoadClass = regexp.MustCompile(`(?:Failed to load|Failed to instantiate|Failed to resolve class|Failed to inject)\s+(?:class\s+)?([A-Za-z_$][\w$]*(?:\.[\w$]+)+)`)

// attribute moves log errors that name a class to the plugin that ships the
// class; the rest stay general errors.
func attribute(errs []string, pluginDir string, failed map[string]bool, versions map[string]string) ([]Failure, []string) {
	var idx classIndex
	broken := map[string]*Failure{}
	var order []string
	var rest []string
	for _, e := range errs {
		m := failedToLoadClass.FindStringSubmatch(e)
		if m == nil {
			rest = append(rest, e)
			continue
		}
		if idx == nil {
			idx = buildClassIndex(pluginDir)
		}
		owner := idx.owner(m[1])
		if owner == "" {
			rest = append(rest, e)
			continue
		}
		if failed[owner] {
			continue // already explained by the plugin's own failure
		}
		b, ok := broken[owner]
		if !ok {
			b = &Failure{Name: owner, Version: versions[owner]}
			broken[owner] = b
			order = append(order, owner)
		}
		if len(b.Reasons) < 5 {
			b.Reasons = append(b.Reasons, e)
		}
	}
	sort.Strings(order)
	out := make([]Failure, 0, len(order))
	for _, n := range order {
		out = append(out, *broken[n])
	}
	return out, rest
}

var unsafeChars = regexp.MustCompile(`[^A-Za-z0-9_.-]`)

func sanitize(s string) string { return unsafeChars.ReplaceAllString(s, "_") }

func (o *Options) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, o.docker(), args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

// pickImage pulls the first image that exists for the core, on the oldest
// Java that both the core and the plugins support: that is what a
// controller of this version most likely runs today.
func (o *Options) pickImage(ctx context.Context, minJava int) (string, int, error) {
	if o.Image != "" {
		o.logf("pulling %s", o.Image)
		if out, err := o.run(ctx, "pull", "-q", o.Image); err != nil {
			return "", 0, fmt.Errorf("docker pull %s: %s", o.Image, strings.TrimSpace(out))
		}
		return o.Image, o.Java, nil
	}
	javas := version.SupportedJava(o.Core)
	var order []int
	if o.Java > 0 {
		order = append(order, o.Java)
	}
	for _, j := range javas {
		if j != o.Java && j >= minJava {
			order = append(order, j)
		}
	}
	if len(order) == 0 || (o.Java > 0 && len(order) == 1) {
		// No Java the core supports satisfies every plugin, so the ones
		// that need more fail whatever runs; rehearse what runs today.
		for _, j := range javas {
			if j != o.Java {
				order = append(order, j)
			}
		}
	}
	var tried []string
	var last string
	for _, j := range order {
		c := fmt.Sprintf("jenkins/jenkins:%s-jdk%d", o.Core, j)
		tried = append(tried, c)
		o.logf("pulling %s", c)
		out, err := o.run(ctx, "pull", "-q", c)
		if err == nil {
			return c, j, nil
		}
		last = strings.TrimSpace(out)
	}
	c := "jenkins/jenkins:" + o.Core
	tried = append(tried, c)
	out, err := o.run(ctx, "pull", "-q", c)
	if err == nil {
		return c, 0, nil
	}
	last = strings.TrimSpace(out)
	return "", 0, fmt.Errorf("no Jenkins image found for %s (tried %s): %s", o.Core, strings.Join(tried, ", "), last)
}

func (o *Options) client() *http.Client {
	if o.Client != nil {
		return o.Client
	}
	return &http.Client{Timeout: 5 * time.Minute}
}

// download fetches the plugins into dir. A plugin that cannot be fetched is
// reported, not fatal: the rehearsal still shows what the rest does.
func (o *Options) download(ctx context.Context, picks []plan.Pick, dir string) ([]plan.Pick, []string) {
	type result struct {
		p   plan.Pick
		err error
	}
	jobs := make(chan plan.Pick)
	results := make(chan result, len(picks))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				results <- result{p, o.fetchOne(ctx, p, dir)}
			}
		}()
	}
	for _, p := range picks {
		jobs <- p
	}
	close(jobs)
	wg.Wait()
	close(results)
	var got []plan.Pick
	var skipped []string
	for r := range results {
		if r.err != nil {
			skipped = append(skipped, fmt.Sprintf("%s %s: %v", r.p.Name, r.p.Version, r.err))
			continue
		}
		got = append(got, r.p)
	}
	sort.Slice(got, func(i, j int) bool { return got[i].Name < got[j].Name })
	sort.Strings(skipped)
	return got, skipped
}

func (o *Options) fetchOne(ctx context.Context, p plan.Pick, dir string) error {
	cache := o.CacheDir
	if cache == "" {
		cache = filepath.Join(os.TempDir(), "jenkins-upgrade-check-plugins")
	}
	cached := filepath.Join(cache, p.Name, p.Version, p.Name+".hpi")
	if _, err := os.Stat(cached); err != nil {
		if err := o.get(ctx, p.URL, cached); err != nil {
			return err
		}
	}
	if err := verify(cached, p.SHA256); err != nil {
		os.Remove(cached)
		return err
	}
	return copyFile(cached, filepath.Join(dir, p.Name+".jpi"))
}

func (o *Options) get(ctx context.Context, url, dest string) error {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt*2) * time.Second)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := o.client().Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("GET %s: %s", url, resp.Status)
			if resp.StatusCode == http.StatusNotFound {
				return lastErr
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			resp.Body.Close()
			return err
		}
		tmp, err := os.CreateTemp(filepath.Dir(dest), ".part-*")
		if err != nil {
			resp.Body.Close()
			return err
		}
		_, err = io.Copy(tmp, resp.Body)
		resp.Body.Close()
		if cerr := tmp.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			os.Remove(tmp.Name())
			lastErr = err
			continue
		}
		return os.Rename(tmp.Name(), dest)
	}
	return lastErr
}

// verify checks a download against the update center's SHA-256, which it
// publishes base64-encoded.
func verify(path, want string) error {
	if want == "" {
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	sum := h.Sum(nil)
	if want == base64.StdEncoding.EncodeToString(sum) || strings.EqualFold(want, hex.EncodeToString(sum)) {
		return nil
	}
	return fmt.Errorf("checksum mismatch for %s", filepath.Base(path))
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// copyTree copies a file, or a directory of regular files.
func copyTree(src, dst string) error {
	st, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !st.IsDir() {
		return copyFile(src, filepath.Join(dst, filepath.Base(src)))
	}
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !d.Type().IsRegular() {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		return copyFile(path, filepath.Join(dst, rel))
	})
}

var (
	failedLoading = regexp.MustCompile(`Failed Loading plugin (.*?) v(\S+) \(([^)\s]+)\)`)
	logHeader     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}[ T][\d:.]+\S*\s+\[id=\d+\]\s+(SEVERE|WARNING)\s+(.*)$`)
)

// logLine trims a log line and the "[LF]>" marker newer Jenkins puts in
// front of each continuation line of a multi-line message.
func logLine(s string) string {
	t := strings.TrimSpace(s)
	for strings.HasPrefix(t, "[LF]>") {
		t = strings.TrimSpace(strings.TrimPrefix(t, "[LF]>"))
	}
	return t
}

// ParseLog finds the plugins Jenkins reported as failed to load, with the
// reasons it gave, and other severe errors worth a look.
func ParseLog(log string) ([]Failure, []string) {
	var failed []Failure
	var errs []string
	seenErr := map[string]bool{}
	sc := bufio.NewScanner(strings.NewReader(log))
	sc.Buffer(make([]byte, 1024*1024), 16*1024*1024)
	var cur *Failure
	var exc []string // exception lines, for a failure with no " - " reasons
	var pendingErr string
	done := func() {
		if cur != nil && len(cur.Reasons) == 0 {
			cur.Reasons = exc
		}
		cur, exc = nil, nil
	}
	for sc.Scan() {
		line := sc.Text()
		if m := failedLoading.FindStringSubmatch(line); m != nil {
			done()
			failed = append(failed, Failure{Name: m[3], Version: m[2]})
			cur = &failed[len(failed)-1]
			pendingErr = ""
			continue
		}
		if cur != nil && !logHeader.MatchString(line) {
			t := logLine(line)
			switch {
			case t == "", strings.HasPrefix(t, "at "), strings.HasPrefix(t, "..."):
			case strings.HasPrefix(t, "- "):
				cur.Reasons = append(cur.Reasons, strings.TrimPrefix(t, "- "))
			case strings.HasPrefix(t, "java.io.IOException: Failed to load"),
				strings.HasPrefix(t, "Caused: java.io.IOException: Failed to load"):
			default:
				t = strings.TrimPrefix(strings.TrimPrefix(t, "Caused: "), "Caused by: ")
				if len(exc) < 3 {
					exc = append(exc, t)
				}
			}
			continue
		}
		done()
		if m := logHeader.FindStringSubmatch(line); m != nil {
			pendingErr = ""
			if m[1] == "SEVERE" {
				pendingErr = m[2]
			} else if strings.Contains(m[2], "NoClassDefFoundError") || strings.Contains(m[2], "ClassNotFoundException") || strings.Contains(m[2], "LinkageError") {
				pendingErr = m[2]
			}
			if strings.Contains(m[2], "optional component") {
				pendingErr = ""
			}
			continue
		}
		if pendingErr != "" {
			msg := pendingErr
			if t := logLine(line); t != "" && !strings.HasPrefix(t, "at ") {
				msg += " — " + t
			}
			if !seenErr[msg] && len(errs) < 40 {
				seenErr[msg] = true
				errs = append(errs, msg)
			}
			pendingErr = ""
		}
	}
	done()
	// A plugin can be reported more than once; keep the first with reasons.
	uniq := map[string]int{}
	var out []Failure
	for _, f := range failed {
		if i, ok := uniq[f.Name]; ok {
			if len(out[i].Reasons) == 0 {
				out[i].Reasons = f.Reasons
			}
			continue
		}
		uniq[f.Name] = len(out)
		out = append(out, f)
	}
	return out, errs
}
