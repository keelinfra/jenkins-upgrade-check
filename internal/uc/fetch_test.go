package uc

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestOpenRetries(t *testing.T) {
	retryWait = 0
	calls, failures := 0, 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls <= failures {
			http.NotFound(w, r)
			return
		}
		io.WriteString(w, "2.580.1")
	}))
	defer srv.Close()
	s := &Source{Base: srv.URL}

	// A mirror that 404s twice and then answers.
	failures = 2
	r, err := s.Open(context.Background(), "stable/latestCore.txt")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := io.ReadAll(r)
	r.Close()
	if string(b) != "2.580.1" || calls != 3 {
		t.Fatalf("got %q after %d calls", b, calls)
	}

	// One that never answers: give up after three tries.
	calls, failures = 0, 99
	if _, err := s.Open(context.Background(), "stable/latestCore.txt"); err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("want a 404, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("want three tries, got %d", calls)
	}
}

// withRetryWait sets how long get waits between tries for one test.
func withRetryWait(t *testing.T, d time.Duration) {
	old := retryWait
	retryWait = d
	t.Cleanup(func() { retryWait = old })
}

func TestBase(t *testing.T) {
	if got := (&Source{}).base(); got != DefaultBase {
		t.Errorf("default base = %q", got)
	}
	if got := (&Source{Base: "https://mirror.example.com/jenkins/"}).base(); got != "https://mirror.example.com/jenkins" {
		t.Errorf("base keeps no trailing slash: %q", got)
	}
}

// With --update-center-dir the files come from the directory, and a
// missing one says where to download it.
func TestOpenDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "plugin-versions.json"), []byte(`{"plugins":{"git":{"5.2.1":{"requiredCore":"2.426.3"}}}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	s := &Source{Dir: dir}
	rel, err := s.Releases(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if r := rel.Get("git", "5.2.1"); r == nil || r.RequiredCore != "2.426.3" || r.Name != "git" {
		t.Errorf("git 5.2.1 = %+v", r)
	}
	_, err = s.UpdateCenter(context.Background())
	if err == nil || !strings.Contains(err.Error(), "(with --update-center-dir, download it from https://updates.jenkins.io/current/update-center.actual.json)") {
		t.Errorf("err = %v", err)
	}
	if _, err := s.ResolveCore(context.Background(), "lts"); err == nil || err.Error() != "--to lts needs the network; give a version number with --update-center-dir" {
		t.Errorf("err = %v", err)
	}
	_, err = (&Source{Dir: t.TempDir()}).Releases(context.Background())
	if err == nil || !strings.Contains(err.Error(), "download it from https://updates.jenkins.io/current/plugin-versions.json") {
		t.Errorf("err = %v", err)
	}
}

// Files go to current/ on the update center, under the tool's name, and a
// cached copy younger than MaxAge is used without asking again.
func TestOpenCache(t *testing.T) {
	var mu sync.Mutex
	var paths []string
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		if ua := r.Header.Get("User-Agent"); ua != "jenkins-upgrade-check/1.2.3" {
			t.Errorf("User-Agent = %q", ua)
		}
		switch r.URL.Path {
		case "/current/update-center.actual.json":
			io.WriteString(w, `{"core":{"version":"2.580"},"plugins":{"git":{"name":"git"}},"deprecations":{"old":{"url":"https://x"}}}`)
		case "/current/plugin-versions.json":
			io.WriteString(w, `{"plugins":{"git":{"5.2.1":{}}}}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	cache := filepath.Join(t.TempDir(), "cache")
	s := &Source{Base: srv.URL, CacheDir: cache, MaxAge: time.Hour, Client: srv.Client(), UserAgent: "jenkins-upgrade-check/1.2.3"}
	for i := 0; i < 2; i++ {
		u, err := s.UpdateCenter(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if u.Core.Version != "2.580" || !u.Distributed("git") {
			t.Errorf("update center = %+v", u)
		}
	}
	if len(paths) != 1 {
		t.Errorf("a fresh cache must not be fetched again: %q", paths)
	}
	if _, err := os.Stat(filepath.Join(cache, "current_update-center.actual.json")); err != nil {
		t.Errorf("not cached: %v", err)
	}

	// A copy older than MaxAge is fetched again.
	s.MaxAge = 0
	if _, err := s.UpdateCenter(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 {
		t.Errorf("a stale cache must be fetched again: %q", paths)
	}

	// Without a cache, straight from the network.
	s.CacheDir = ""
	rel, err := s.Releases(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !rel.Known("git") {
		t.Error("git not read")
	}
}

func TestOpenCacheErrors(t *testing.T) {
	withRetryWait(t, 0)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/current/truncated.json":
			// Promise more than is sent: the copy fails half way.
			w.Header().Set("Content-Length", "1000")
			io.WriteString(w, `{"plugins":`)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	cache := t.TempDir()
	s := &Source{Base: srv.URL, CacheDir: cache, MaxAge: time.Hour}
	_, err := s.Open(context.Background(), "plugin-versions.json")
	if want := "GET " + srv.URL + "/current/plugin-versions.json: 500 Internal Server Error"; err == nil || err.Error() != want {
		t.Errorf("err = %v, want %s", err, want)
	}
	_, err = s.Open(context.Background(), "truncated.json")
	if err == nil || !strings.Contains(err.Error(), "GET "+srv.URL+"/current/truncated.json: unexpected EOF") {
		t.Errorf("err = %v", err)
	}
	if left, _ := filepath.Glob(filepath.Join(cache, "*")); len(left) != 0 {
		t.Errorf("a failed download leaves files behind: %q", left)
	}

	// The cache directory cannot be made under a file.
	file := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	s.CacheDir = filepath.Join(file, "cache")
	srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "{}") }))
	defer srv2.Close()
	s.Base = srv2.URL
	if _, err := s.Open(context.Background(), "plugin-versions.json"); err == nil || !strings.Contains(err.Error(), "not a directory") {
		t.Errorf("err = %v", err)
	}
}

// A server that is not there: three tries, then the connection error.
func TestGetUnreachable(t *testing.T) {
	withRetryWait(t, 0)
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close()
	_, err := (&Source{Base: url}).Open(context.Background(), "plugin-versions.json")
	if err == nil || !strings.HasPrefix(err.Error(), "GET "+url+"/current/plugin-versions.json: ") || !strings.Contains(err.Error(), "refused") {
		t.Errorf("err = %v", err)
	}
	if _, err := (&Source{Base: "http://[::1"}).Open(context.Background(), "plugin-versions.json"); err == nil || !strings.Contains(err.Error(), "missing ']'") {
		t.Errorf("a malformed base URL: %v", err)
	}
}

// A cancelled context stops the wait between tries.
func TestGetCancelledWhileWaiting(t *testing.T) {
	withRetryWait(t, time.Hour)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := (&Source{Base: srv.URL}).Open(ctx, "plugin-versions.json")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want the context's", err)
	}
	if d := time.Since(start); d > 10*time.Second {
		t.Errorf("took %s: the wait did not stop", d)
	}
}

func TestResolveCore(t *testing.T) {
	withRetryWait(t, 0)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/stable/latestCore.txt":
			io.WriteString(w, "2.580.1\n")
		case "/current/latestCore.txt":
			io.WriteString(w, "  \n")
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	s := &Source{Base: srv.URL}
	ctx := context.Background()
	for _, in := range []string{"lts", "stable"} {
		if v, err := s.ResolveCore(ctx, in); err != nil || v != "2.580.1" {
			t.Errorf("ResolveCore(%s) = %q, %v", in, v, err)
		}
	}
	if v, err := s.ResolveCore(ctx, "2.555.1"); err != nil || v != "2.555.1" {
		t.Errorf("a version is returned as given: %q, %v", v, err)
	}
	for _, in := range []string{"weekly", "latest"} {
		if _, err := s.ResolveCore(ctx, in); err == nil || err.Error() != "current/latestCore.txt is empty" {
			t.Errorf("ResolveCore(%s): %v", in, err)
		}
	}

	truncated := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "10")
		io.WriteString(w, "2.5")
	}))
	defer truncated.Close()
	if _, err := (&Source{Base: truncated.URL}).ResolveCore(ctx, "lts"); err == nil || !strings.Contains(err.Error(), "unexpected EOF") {
		t.Errorf("a cut-off answer: %v", err)
	}
	gone := httptest.NewServer(http.NotFoundHandler())
	defer gone.Close()
	if _, err := (&Source{Base: gone.URL}).ResolveCore(ctx, "lts"); err == nil || !strings.Contains(err.Error(), "404 Not Found") {
		t.Errorf("a missing file: %v", err)
	}
}
