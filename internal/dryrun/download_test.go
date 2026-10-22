package dryrun

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
)

// Download fetches what it can, checks every checksum, retries what may be
// passing trouble, and says why each plugin it could not get is missing.
func TestDownload(t *testing.T) {
	fastWaits(t)
	ok := []byte("ok plugin")
	hexed := []byte("hex plugin")
	var mu sync.Mutex
	hits := map[string]int{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		hits[r.URL.Path]++
		n := hits[r.URL.Path]
		mu.Unlock()
		switch r.URL.Path {
		case "/ok.hpi":
			w.Write(ok)
		case "/hex.hpi":
			w.Write(hexed)
		case "/flaky.hpi":
			if n == 1 {
				w.WriteHeader(http.StatusBadGateway)
				return
			}
			w.Write([]byte("flaky plugin"))
		case "/down.hpi":
			w.WriteHeader(http.StatusServiceUnavailable)
		case "/cut.hpi":
			w.Header().Set("Content-Length", "100")
			w.Write([]byte("partial"))
		case "/bad.hpi":
			w.Write([]byte("tampered"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	count := func(path string) int {
		mu.Lock()
		defer mu.Unlock()
		return hits[path]
	}
	gone := httptest.NewServer(http.NotFoundHandler())
	goneURL := gone.URL
	gone.Close()

	hexSum := sha256Hex(hexed)
	picks := []plan.Pick{
		{Name: "ok", Version: "1.0", URL: srv.URL + "/ok.hpi", SHA256: sum(ok)},
		{Name: "hex", Version: "1.0", URL: srv.URL + "/hex.hpi", SHA256: hexSum},
		{Name: "flaky", Version: "1.0", URL: srv.URL + "/flaky.hpi"},
		{Name: "down", Version: "1.0", URL: srv.URL + "/down.hpi"},
		{Name: "cut", Version: "1.0", URL: srv.URL + "/cut.hpi"},
		{Name: "bad", Version: "1.0", URL: srv.URL + "/bad.hpi", SHA256: sum(ok)},
		{Name: "missing", Version: "1.0", URL: srv.URL + "/missing.hpi"},
		{Name: "unreachable", Version: "1.0", URL: goneURL + "/x.hpi"},
		{Name: "malformed", Version: "1.0", URL: "http://[::1/x.hpi"},
		{Name: "private", Version: "0.1"},
	}
	dir, cache := t.TempDir(), t.TempDir()
	got, skipped := Download(context.Background(), picks, dir, cache)

	var names []string
	for _, p := range got {
		names = append(names, p.Name)
	}
	if !reflect.DeepEqual(names, []string{"flaky", "hex", "ok"}) {
		t.Errorf("got %q", names)
	}
	for name, want := range map[string]string{"ok": "ok plugin", "hex": "hex plugin", "flaky": "flaky plugin"} {
		if b, err := os.ReadFile(filepath.Join(dir, name+".jpi")); err != nil || string(b) != want {
			t.Errorf("%s.jpi = %q, %v", name, b, err)
		}
	}
	want := map[string]string{
		"private":     "private 0.1: no download known",
		"bad":         "bad 1.0: checksum mismatch for bad.hpi",
		"cut":         "cut 1.0: unexpected EOF",
		"down":        "down 1.0: GET " + srv.URL + "/down.hpi: 503 Service Unavailable",
		"missing":     "missing 1.0: GET " + srv.URL + "/missing.hpi: 404 Not Found",
		"unreachable": "unreachable 1.0: Get \"" + goneURL + "/x.hpi\": dial tcp",
		"malformed":   "malformed 1.0: parse \"http://[::1/x.hpi\": missing ']' in host",
	}
	if len(skipped) != len(want) || skipped[0] != want["private"] {
		t.Errorf("skipped = %q", skipped)
	}
	for _, s := range skipped {
		name, _, _ := strings.Cut(s, " ")
		if !strings.HasPrefix(s, want[name]) {
			t.Errorf("skipped %q, want it to start with %q", s, want[name])
		}
	}
	// A 404 is final; other failures are tried three times.
	for path, n := range map[string]int{"/ok.hpi": 1, "/flaky.hpi": 2, "/down.hpi": 3, "/cut.hpi": 3, "/missing.hpi": 1} {
		if got := count(path); got != n {
			t.Errorf("%s fetched %d times, want %d", path, got, n)
		}
	}
	// A download that fails its checksum is not kept.
	if _, err := os.Stat(filepath.Join(cache, "bad", "1.0", "bad.hpi")); !os.IsNotExist(err) {
		t.Errorf("a bad download stays in the cache: %v", err)
	}
	if left, _ := filepath.Glob(filepath.Join(cache, "cut", "1.0", "*")); len(left) != 0 {
		t.Errorf("a cut-off download leaves %q", left)
	}

	// The second time, what is cached is not fetched again.
	got, _ = Download(context.Background(), picks[:1], t.TempDir(), cache)
	if len(got) != 1 || count("/ok.hpi") != 1 {
		t.Errorf("got %v after %d fetches", got, count("/ok.hpi"))
	}
}

func TestDownloadUnwritable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("x")) }))
	defer srv.Close()
	file := filepath.Join(t.TempDir(), "file")
	os.WriteFile(file, nil, 0o644)
	picks := []plan.Pick{{Name: "p", Version: "1", URL: srv.URL + "/p.hpi"}}

	// The cache cannot be made.
	_, skipped := Download(context.Background(), picks, t.TempDir(), filepath.Join(file, "cache"))
	if len(skipped) != 1 || !strings.HasPrefix(skipped[0], "p 1: mkdir ") || !strings.HasSuffix(skipped[0], "not a directory") {
		t.Errorf("skipped = %q", skipped)
	}
	// The plugins directory cannot be made.
	_, skipped = Download(context.Background(), picks, filepath.Join(file, "plugins"), t.TempDir())
	if len(skipped) != 1 || !strings.HasSuffix(skipped[0], "not a directory") {
		t.Errorf("skipped = %q", skipped)
	}
}

// Without a cache directory, downloads are kept in the system's temporary
// directory.
func TestDownloadDefaultCache(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("x")) }))
	defer srv.Close()
	got, skipped := Download(context.Background(), []plan.Pick{{Name: "p", Version: "1", URL: srv.URL + "/p.hpi"}}, t.TempDir(), "")
	if len(got) != 1 || len(skipped) != 0 {
		t.Fatalf("got %v, skipped %q", got, skipped)
	}
	if _, err := os.Stat(filepath.Join(tmp, "jenkins-upgrade-check-plugins", "p", "1", "p.hpi")); err != nil {
		t.Errorf("not cached in the temporary directory: %v", err)
	}
}

func sha256Hex(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}
