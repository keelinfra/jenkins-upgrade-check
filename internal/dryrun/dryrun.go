// Package dryrun rehearses an upgrade: it starts a throwaway Jenkins of the
// target version in Docker with exactly the plugin set given, and reports
// which plugins failed to load and why. Nothing from the real controller goes
// into the container except the plugin list: no jobs, no credentials, no
// configuration.
package dryrun

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
	"github.com/keelinfra/jenkins-upgrade-check/internal/version"
)

// Options configure one rehearsal.
type Options struct {
	Core    string
	Java    int // preferred Java for the image; 0 picks one the core supports
	Plugins []plan.Pick
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
	Core           string    `json:"core"`
	Image          string    `json:"image"`
	Plugins        int       `json:"plugins"`
	Started        bool      `json:"started"`
	Reason         string    `json:"reason,omitempty"`
	StartupSeconds float64   `json:"startupSeconds,omitempty"`
	Failed         []Failure `json:"failed,omitempty"`
	Skipped        []string  `json:"skipped,omitempty"`
	Errors         []string  `json:"errors,omitempty"`
	LogFile        string    `json:"logFile,omitempty"`
}

// Failure is a plugin Jenkins did not load.
type Failure struct {
	Name    string   `json:"name"`
	Version string   `json:"version,omitempty"`
	Reasons []string `json:"reasons,omitempty"`
}

// OK reports whether Jenkins started and loaded every plugin.
func (r *Result) OK() bool { return r.Started && len(r.Failed) == 0 }

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
	pluginDir := filepath.Join(work, "plugins")
	if err := os.MkdirAll(pluginDir, 0o755); err != nil {
		return nil, err
	}

	var fetch []plan.Pick
	for _, p := range o.Plugins {
		if p.URL == "" {
			res.Skipped = append(res.Skipped, p.Name)
			continue
		}
		fetch = append(fetch, p)
	}
	res.Plugins = len(fetch)
	o.logf("downloading %d plugins", len(fetch))
	if err := o.download(ctx, fetch, pluginDir); err != nil {
		return nil, err
	}

	base, err := o.pickImage(ctx)
	if err != nil {
		return nil, err
	}
	res.Image = base
	dockerfile := "FROM " + base + "\n" +
		"ENV JAVA_OPTS=\"-Djenkins.install.runSetupWizard=false -Dhudson.model.UpdateCenter.never=true\"\n" +
		"COPY --chown=jenkins:jenkins plugins/ /usr/share/jenkins/ref/plugins/\n"
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

	name := strings.ReplaceAll(strings.TrimPrefix(tag, "jenkins-upgrade-check:"), ".", "-")
	name = "juc-" + name
	if out, err := o.run(ctx, "run", "-d", "--name", name, "-p", "127.0.0.1::8080", tag); err != nil {
		return nil, fmt.Errorf("docker run: %v\n%s", err, out)
	}
	if !o.Keep {
		defer o.run(context.Background(), "rm", "-f", name)
	}
	o.logf("starting Jenkins %s (container %s)", o.Core, name)

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
	if res.Started {
		// Late failures (extension loading) land just after the banner.
		time.Sleep(5 * time.Second)
		logs, _ = o.run(ctx, "logs", name)
	}

	failed, errs := ParseLog(logs)
	res.Errors = errs
	byName := map[string]*Failure{}
	for i := range failed {
		byName[failed[i].Name] = &failed[i]
	}
	if res.Started {
		if active, err := o.activePlugins(ctx, name); err == nil {
			for _, p := range fetch {
				if active[p.Name] {
					continue
				}
				if _, ok := byName[p.Name]; !ok {
					failed = append(failed, Failure{Name: p.Name, Version: p.Version, Reasons: []string{"not active after startup"}})
					byName[p.Name] = &failed[len(failed)-1]
				}
			}
		} else {
			res.Errors = append(res.Errors, "could not read the plugin list from the rehearsal controller: "+err.Error())
		}
	}
	sort.Slice(failed, func(i, j int) bool { return failed[i].Name < failed[j].Name })
	res.Failed = failed

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

func sanitize(s string) string {
	return regexp.MustCompile(`[^A-Za-z0-9_.-]`).ReplaceAllString(s, "_")
}

func (o *Options) run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, o.docker(), args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

// pickImage pulls the first image tag that exists for the core.
func (o *Options) pickImage(ctx context.Context) (string, error) {
	var candidates []string
	if o.Image != "" {
		candidates = []string{o.Image}
	} else {
		javas := version.SupportedJava(o.Core)
		order := []int{}
		if o.Java > 0 {
			order = append(order, o.Java)
		}
		// Oldest first: that is what a controller of this version most
		// likely runs today.
		for _, j := range javas {
			if j != o.Java {
				order = append(order, j)
			}
		}
		for _, j := range order {
			candidates = append(candidates, fmt.Sprintf("jenkins/jenkins:%s-jdk%d", o.Core, j))
		}
		candidates = append(candidates, "jenkins/jenkins:"+o.Core)
	}
	var last string
	for _, c := range candidates {
		o.logf("pulling %s", c)
		out, err := o.run(ctx, "pull", "-q", c)
		if err == nil {
			return c, nil
		}
		last = strings.TrimSpace(out)
	}
	return "", fmt.Errorf("no Jenkins image found for %s (tried %s): %s", o.Core, strings.Join(candidates, ", "), last)
}

func (o *Options) client() *http.Client {
	if o.Client != nil {
		return o.Client
	}
	return &http.Client{Timeout: 5 * time.Minute}
}

func (o *Options) download(ctx context.Context, picks []plan.Pick, dir string) error {
	jobs := make(chan plan.Pick)
	errc := make(chan error, len(picks))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				if err := o.fetchOne(ctx, p, dir); err != nil {
					errc <- fmt.Errorf("%s %s: %w", p.Name, p.Version, err)
				}
			}
		}()
	}
	for _, p := range picks {
		jobs <- p
	}
	close(jobs)
	wg.Wait()
	close(errc)
	var errs []error
	for err := range errc {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
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

func (o *Options) activePlugins(ctx context.Context, container string) (map[string]bool, error) {
	out, err := o.run(ctx, "port", container, "8080")
	if err != nil {
		return nil, fmt.Errorf("docker port: %v", err)
	}
	addr := strings.TrimSpace(strings.Split(strings.TrimSpace(out), "\n")[0])
	if addr == "" {
		return nil, errors.New("no published port")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+addr+"/pluginManager/api/json?tree=plugins[shortName,active]", nil)
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("plugin list: %s", resp.Status)
	}
	var pm struct {
		Plugins []struct {
			ShortName string `json:"shortName"`
			Active    bool   `json:"active"`
		} `json:"plugins"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pm); err != nil {
		return nil, err
	}
	active := map[string]bool{}
	for _, p := range pm.Plugins {
		active[p.ShortName] = p.Active
	}
	return active, nil
}

var (
	failedLoading = regexp.MustCompile(`Failed Loading plugin (.*?) v(\S+) \(([^)\s]+)\)`)
	logHeader     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}[ T][\d:.]+\S*\s+\[id=\d+\]\s+(SEVERE|WARNING)\s+(.*)$`)
)

// ParseLog finds the plugins Jenkins reported as failed to load, with the
// reasons it gave, and other severe errors worth a look.
func ParseLog(log string) ([]Failure, []string) {
	var failed []Failure
	var errs []string
	seenErr := map[string]bool{}
	sc := bufio.NewScanner(strings.NewReader(log))
	sc.Buffer(make([]byte, 1024*1024), 16*1024*1024)
	var cur *Failure
	var pendingErr string
	for sc.Scan() {
		line := sc.Text()
		if m := failedLoading.FindStringSubmatch(line); m != nil {
			failed = append(failed, Failure{Name: m[3], Version: m[2]})
			cur = &failed[len(failed)-1]
			pendingErr = ""
			continue
		}
		if cur != nil {
			t := strings.TrimSpace(line)
			switch {
			case strings.HasPrefix(t, "- "):
				cur.Reasons = append(cur.Reasons, strings.TrimPrefix(t, "- "))
				continue
			case strings.HasPrefix(t, "java.io.IOException: Failed to load"):
				continue
			case strings.HasPrefix(t, "at ") || strings.HasPrefix(t, "Caused by") || strings.HasPrefix(t, "..."):
				continue
			default:
				cur = nil
			}
		}
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
			if t := strings.TrimSpace(line); t != "" && !strings.HasPrefix(t, "at ") {
				msg += " — " + t
			}
			if !seenErr[msg] && len(errs) < 40 {
				seenErr[msg] = true
				errs = append(errs, msg)
			}
			pendingErr = ""
		}
	}
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
