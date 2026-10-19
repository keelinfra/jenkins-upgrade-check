package uc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Default locations on the Jenkins update center.
const (
	DefaultBase = "https://updates.jenkins.io"
)

// Source fetches update center files, from the network or from a directory
// the user filled by hand (for controllers that cannot reach the internet).
type Source struct {
	// Base is the update center root, DefaultBase unless overridden.
	Base string
	// Dir, when set, is read instead of the network: it holds
	// update-center.actual.json and plugin-versions.json.
	Dir string
	// CacheDir keeps downloads for MaxAge so repeated runs stay fast.
	CacheDir string
	MaxAge   time.Duration
	Client   *http.Client
	// UserAgent identifies the tool to the update center.
	UserAgent string
}

func (s *Source) base() string {
	if s.Base != "" {
		return strings.TrimRight(s.Base, "/")
	}
	return DefaultBase
}

func (s *Source) client() *http.Client {
	if s.Client != nil {
		return s.Client
	}
	return &http.Client{Timeout: 10 * time.Minute}
}

// Open returns the named file: "update-center.actual.json",
// "plugin-versions.json", or a "stable/latestCore.txt"-style path.
func (s *Source) Open(ctx context.Context, name string) (io.ReadCloser, error) {
	if s.Dir != "" {
		f, err := os.Open(filepath.Join(s.Dir, filepath.Base(name)))
		if err != nil {
			return nil, fmt.Errorf("%w (with --update-center-dir, download it from %s/current/%s)", err, s.base(), filepath.Base(name))
		}
		return f, nil
	}
	path := name
	if !strings.Contains(name, "/") {
		path = "current/" + name
	}
	if s.CacheDir != "" {
		cached := filepath.Join(s.CacheDir, strings.ReplaceAll(path, "/", "_"))
		if st, err := os.Stat(cached); err == nil && time.Since(st.ModTime()) < s.MaxAge {
			return os.Open(cached)
		}
		if err := s.download(ctx, s.base()+"/"+path, cached); err != nil {
			return nil, err
		}
		return os.Open(cached)
	}
	resp, err := s.get(ctx, s.base()+"/"+path)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (s *Source) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if s.UserAgent != "" {
		req.Header.Set("User-Agent", s.UserAgent)
	}
	resp, err := s.client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return resp, nil
}

func (s *Source) download(ctx context.Context, url, dest string) error {
	resp, err := s.get(ctx, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(dest), ".download-*")
	if err != nil {
		return err
	}
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return fmt.Errorf("GET %s: %w", url, err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return os.Rename(tmp.Name(), dest)
}

// UpdateCenter loads update-center.actual.json.
func (s *Source) UpdateCenter(ctx context.Context) (*UpdateCenter, error) {
	r, err := s.Open(ctx, "update-center.actual.json")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ReadUpdateCenter(r)
}

// Releases loads plugin-versions.json.
func (s *Source) Releases(ctx context.Context) (*Releases, error) {
	r, err := s.Open(ctx, "plugin-versions.json")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ReadReleases(r)
}

// ResolveCore turns "lts" or "weekly" into a version number; anything else is
// returned as given.
func (s *Source) ResolveCore(ctx context.Context, core string) (string, error) {
	var path string
	switch core {
	case "lts", "stable":
		path = "stable/latestCore.txt"
	case "weekly", "latest":
		path = "current/latestCore.txt"
	default:
		return core, nil
	}
	if s.Dir != "" {
		return "", fmt.Errorf("--to %s needs the network; give a version number with --update-center-dir", core)
	}
	resp, err := s.get(ctx, s.base()+"/"+path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return "", err
	}
	v := strings.TrimSpace(string(b))
	if v == "" {
		return "", fmt.Errorf("%s is empty", path)
	}
	return v, nil
}
