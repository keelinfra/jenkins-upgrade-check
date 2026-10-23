package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/keelinfra/jenkins-upgrade-check/internal/uc"
)

// fakeWeb sends every request of the package's client to h, whatever the
// host, and returns the update center h also serves.
func fakeWeb(t *testing.T, h http.HandlerFunc) {
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	target, _ := url.Parse(srv.URL)
	oldClient, oldUC, oldNow := client, updateCenter, now
	t.Cleanup(func() { client, updateCenter, now = oldClient, oldUC, oldNow })
	client = &http.Client{Transport: rewrite{target}}
	updateCenter = &uc.Source{Base: srv.URL + "/uc"}
	now = func() time.Time { return time.Date(2026, 3, 9, 10, 0, 0, 0, time.FixedZone("CET", 3600)) }
}

// rewrite sends a request to target instead, keeping its path.
type rewrite struct{ target *url.URL }

func (r rewrite) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme, req.URL.Host = r.target.Scheme, r.target.Host
	return http.DefaultTransport.RoundTrip(req)
}

func matrix(args ...string) (int, string, string) {
	var out, errOut bytes.Buffer
	code := run(args, &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestRunUsage(t *testing.T) {
	if code, _, stderr := matrix(); code != 2 || stderr != "usage: matrix month | versions -to VERSION | table -results DIR\n" {
		t.Errorf("no command: exit %d, stderr %q", code, stderr)
	}
	if code, _, stderr := matrix("tabel"); code != 1 || stderr != "matrix: unknown command \"tabel\"\n" {
		t.Errorf("unknown command: exit %d, stderr %q", code, stderr)
	}
	for _, cmd := range []string{"versions", "table"} {
		if code, _, stderr := matrix(cmd, "-bogus"); code != 1 || !strings.HasSuffix(stderr, "matrix: flag provided but not defined: -bogus\n") {
			t.Errorf("%s -bogus: exit %d, stderr %q", cmd, code, stderr)
		}
	}
}

const statsPath = "/jenkins-infra/infra-statistics/gh-pages/jenkins-stats/svg/"

// month is the newest month with a ranking, counting back from the
// clock's month in UTC, across a year.
func TestMonth(t *testing.T) {
	var asked []string
	fakeWeb(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Errorf("%s %s", r.Method, r.URL.Path)
		}
		asked = append(asked, strings.TrimPrefix(r.URL.Path, statsPath))
		if r.URL.Path != statsPath+"202512-top-plugins1000.csv" {
			http.NotFound(w, r)
		}
	})
	code, stdout, stderr := matrix("month")
	if code != 0 || stdout != "202512\n" || stderr != "" {
		t.Errorf("exit %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	if want := []string{"202603-top-plugins1000.csv", "202602-top-plugins1000.csv", "202601-top-plugins1000.csv", "202512-top-plugins1000.csv"}; strings.Join(asked, " ") != strings.Join(want, " ") {
		t.Errorf("asked %q", asked)
	}
}

func TestMonthNone(t *testing.T) {
	n := 0
	fakeWeb(t, func(w http.ResponseWriter, r *http.Request) {
		n++
		http.NotFound(w, r)
	})
	if code, _, stderr := matrix("month"); code != 1 || stderr != "matrix: no usage statistics in the last 24 months\n" || n != 24 {
		t.Errorf("exit %d, stderr %q after %d requests", code, stderr, n)
	}
}

func TestMonthUnreachable(t *testing.T) {
	fakeWeb(t, http.NotFound)
	client = &http.Client{Transport: rewrite{&url.URL{Scheme: "http", Host: closedAddr(t)}}}
	if code, _, stderr := matrix("month"); code != 1 || !strings.Contains(stderr, "connection refused") {
		t.Errorf("exit %d, stderr %q", code, stderr)
	}
}

// closedAddr is an address nothing listens on.
func closedAddr(t *testing.T) string {
	srv := httptest.NewServer(http.NotFoundHandler())
	addr := srv.Listener.Addr().String()
	srv.Close()
	return addr
}

const metadata = `<metadata><versioning><versions>
<version>2.60.1</version><version>2.60.3</version><version>2.541.2</version><version>2.541.3</version><version>2.580.1</version>
</versions></versioning></metadata>`

func TestVersions(t *testing.T) {
	fakeWeb(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/releases/org/jenkins-ci/main/jenkins-war/maven-metadata.xml":
			w.Write([]byte(metadata))
		case "/uc/stable/latestCore.txt":
			w.Write([]byte("2.580.1\n"))
		default:
			http.NotFound(w, r)
		}
	})
	for _, args := range [][]string{{"versions", "-to", "2.580.1"}, {"versions"}} {
		code, stdout, stderr := matrix(args...)
		if code != 0 || stdout != "[\"2.60.3\",\"2.541.3\"]\n" || stderr != "" {
			t.Errorf("%q: exit %d, stdout %q, stderr %q", args, code, stdout, stderr)
		}
	}
	if code, stdout, _ := matrix("versions", "-to", "2.580.1", "-min-line", "100"); code != 0 || stdout != "[\"2.541.3\"]\n" {
		t.Errorf("-min-line 100: exit %d, stdout %q", code, stdout)
	}
	if code, _, stderr := matrix("versions", "-to", "2.60.3"); code != 1 || stderr != "matrix: no LTS line between 2.60 and 2.60.3\n" {
		t.Errorf("nothing older: exit %d, stderr %q", code, stderr)
	}
}

func TestVersionsErrors(t *testing.T) {
	fakeWeb(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/uc/stable/latestCore.txt" {
			return // empty
		}
		http.NotFound(w, r)
	})
	if code, _, stderr := matrix("versions", "-to", "lts"); code != 1 || stderr != "matrix: stable/latestCore.txt is empty\n" {
		t.Errorf("empty latestCore: exit %d, stderr %q", code, stderr)
	}
	if code, _, stderr := matrix("versions", "-to", "2.580.1"); code != 1 || stderr != "matrix: GET "+warIndex+": 404 Not Found\n" {
		t.Errorf("no war index: exit %d, stderr %q", code, stderr)
	}
	client = &http.Client{Transport: rewrite{&url.URL{Scheme: "http", Host: closedAddr(t)}}}
	if code, _, stderr := matrix("versions", "-to", "2.580.1"); code != 1 || !strings.Contains(stderr, "connection refused") {
		t.Errorf("unreachable: exit %d, stderr %q", code, stderr)
	}
}

func TestGetBadURL(t *testing.T) {
	if _, err := get(t.Context(), "http://[::1/x"); err == nil || !strings.Contains(err.Error(), "missing ']' in host") {
		t.Errorf("err = %v", err)
	}
}

// table downloads the statistics for the rows' month unless given a file,
// and writes the table to stdout or -o, and the matrix to -json.
func TestTableCommand(t *testing.T) {
	stats, err := os.ReadFile("testdata/jenkins.csv")
	if err != nil {
		t.Fatal(err)
	}
	fakeWeb(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != statsPath+"202603-jenkins.csv" {
			http.NotFound(w, r)
			return
		}
		w.Write(stats)
	})
	code, stdout, stderr := matrix("table", "-results", "testdata/results")
	if code != 0 || stderr != "" || !strings.HasPrefix(stdout, "# Upgrading every Jenkins LTS line to 2.580.1\n\nRehearsed 2026-03-09 with") ||
		!strings.Contains(stdout, "(Jenkins usage statistics, March 2026)") || !strings.Contains(stdout, "share of the 1000 controllers") {
		t.Errorf("exit %d, stderr %q, stdout:\n%s", code, stderr, stdout)
	}

	dir := t.TempDir()
	md, js := filepath.Join(dir, "matrix.md"), filepath.Join(dir, "matrix.json")
	code, stdout, stderr = matrix("table", "-results", "testdata/results", "-stats", "testdata/jenkins.csv", "-o", md, "-json", js)
	if code != 0 || stdout != "" || stderr != "" {
		t.Fatalf("exit %d, stdout %q, stderr %q", code, stdout, stderr)
	}
	b, err := os.ReadFile(md)
	if err != nil || !strings.Contains(string(b), "| [2.346.3](2.346.3-top300/summary.md) | 30.0% |") {
		t.Errorf("matrix.md = %q, %v", b, err)
	}
	var m Matrix
	if b, err = os.ReadFile(js); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if m.To != "2.580.1" || m.Month != "202603" || m.Controllers != 1000 || len(m.Rows) != 2 || !m.Generated.Equal(time.Date(2026, 3, 9, 9, 0, 0, 0, time.UTC)) {
		t.Errorf("matrix.json = %+v", m)
	}
}

func TestTableCommandErrors(t *testing.T) {
	fakeWeb(t, http.NotFound)
	dir := t.TempDir()
	file := filepath.Join(dir, "file")
	if err := os.WriteFile(file, []byte("\"2.346.3\",\"0\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"-results", dir}, "matrix: no row.json under " + dir + "\n"},
		{[]string{"-results", filepath.Join(dir, "none")}, "matrix: lstat " + filepath.Join(dir, "none") + ": no such file or directory\n"},
		{[]string{"-results", "testdata/results"}, "matrix: GET " + statsBase + "/202603-jenkins.csv: 404 Not Found\n"},
		{[]string{"-results", "testdata/results", "-stats", filepath.Join(dir, "none.csv")}, "matrix: open " + filepath.Join(dir, "none.csv") + ": no such file or directory\n"},
		{[]string{"-results", "testdata/results", "-stats", file}, "matrix: usage statistics: no installations\n"},
		{[]string{"-results", "testdata/results", "-stats", "testdata/jenkins.csv", "-json", filepath.Join(file, "m.json")}, "matrix: open " + filepath.Join(file, "m.json") + ": not a directory\n"},
		{[]string{"-results", "testdata/results", "-stats", "testdata/jenkins.csv", "-o", filepath.Join(file, "m.md")}, "matrix: open " + filepath.Join(file, "m.md") + ": not a directory\n"},
	} {
		if code, _, stderr := matrix(append([]string{"table"}, c.args...)...); code != 1 || stderr != c.want {
			t.Errorf("%q: exit %d, stderr %q, want %q", c.args, code, stderr, c.want)
		}
	}
}

func TestLoadErrors(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "x", "row.json")
	os.MkdirAll(filepath.Dir(bad), 0o755)
	os.WriteFile(bad, []byte("{"), 0o644)
	if _, err := Load(dir); err == nil || err.Error() != bad+": unexpected end of JSON input" {
		t.Errorf("bad JSON: %v", err)
	}
	// A row.json that leads nowhere cannot be read.
	dir = t.TempDir()
	if err := os.Symlink(filepath.Join(dir, "gone"), filepath.Join(dir, "row.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(dir); err == nil || err.Error() != "open "+filepath.Join(dir, "row.json")+": no such file or directory" {
		t.Errorf("unreadable row.json: %v", err)
	}
}

// Rows with more plugins come first; then oldest first.
func TestLoadOrder(t *testing.T) {
	dir := t.TempDir()
	for i, r := range []Row{{From: "2.541.3", Top: 100}, {From: "2.60.3", Top: 100}, {From: "2.541.3", Top: 300, To: "2.580.1", Month: "202603"}} {
		b, _ := json.Marshal(r)
		p := filepath.Join(dir, string(rune('a'+i)), "row.json")
		os.MkdirAll(filepath.Dir(p), 0o755)
		os.WriteFile(p, b, 0o644)
	}
	m, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range m.Rows {
		got = append(got, r.Dir)
	}
	if strings.Join(got, " ") != "2.541.3-top300 2.60.3-top100 2.541.3-top100" || m.To != "2.580.1" || m.Month != "202603" {
		t.Errorf("rows %q, to %q, month %q", got, m.To, m.Month)
	}
}

func TestAddSharesErrors(t *testing.T) {
	m := &Matrix{}
	if err := m.AddShares("\"a\"b,1\n"); err == nil || !strings.HasPrefix(err.Error(), "usage statistics: parse error on line 1") {
		t.Errorf("bad CSV: %v", err)
	}
	// Rows of one field, and counts that are no numbers, are not counted.
	if err := m.AddShares("2.346.3\n2.541.3\n"); err == nil || err.Error() != "usage statistics: no installations" {
		t.Errorf("one field: %v", err)
	}
	if err := m.AddShares("\"version\",\"installations\"\n"); err == nil || err.Error() != "usage statistics: no installations" {
		t.Errorf("header only: %v", err)
	}
}

func TestCells(t *testing.T) {
	long := strings.Repeat("é", 130)
	for _, c := range []struct {
		r    Rehearsal
		want string
	}{
		{Rehearsal{Reason: "no image for Java 8"}, "— not rehearsed: no image for Java 8"},
		{Rehearsal{Image: "i", Plugins: 3}, "❌ Jenkins does not start"},
		{Rehearsal{Image: "i", Reason: "a | b"}, "❌ Jenkins does not start (a / b)"},
		{Rehearsal{Image: "i", Reason: long}, "❌ Jenkins does not start (" + strings.Repeat("é", 120) + "…)"},
		{Rehearsal{Started: true, Plugins: 10, Failed: 5, Skipped: 2, Roots: []Root{{Name: "a"}, {Name: "b"}, {Name: "c"}, {Name: "d"}, {Name: "e"}}},
			"⚠️ 5 of 10 fail: a, b, c, and 2 more (2 not downloadable)"},
		{Rehearsal{Started: true, Plugins: 10, Broken: 1}, "⚠️ 1 of 10 load in part"},
		{Rehearsal{Started: true, Plugins: 10}, "✅ all 10 load"},
	} {
		if got := cell(c.r); got != c.want {
			t.Errorf("cell(%+v) = %q, want %q", c.r, got, c.want)
		}
	}
}

func TestPlannedAndPath(t *testing.T) {
	clean := Rehearsal{Core: "2.541.3", Started: true, Plugins: 5}
	bad := Rehearsal{Core: "2.580.1", Started: true, Plugins: 5, Failed: 1, Roots: []Root{{Name: "x"}}}
	for _, c := range []struct {
		r    Row
		want string
	}{
		{Row{}, "— no plan"},
		{Row{Stops: []Rehearsal{bad}}, "⚠️ 1 of 5 fail: x"},
		{Row{Stops: []Rehearsal{clean, bad}}, "⚠️ 1 of 5 fail: x (at 2.580.1)"},
		{Row{Stops: []Rehearsal{clean}}, "✅ all 5 load"},
	} {
		if got := planned(c.r); got != c.want {
			t.Errorf("planned = %q, want %q", got, c.want)
		}
	}
	if (Row{}).PlannedClean() || (Row{Stops: []Rehearsal{clean, bad}}).PlannedClean() || !(Row{Stops: []Rehearsal{clean}}).PlannedClean() {
		t.Error("PlannedClean")
	}
	// Without a Java of its own, a stop shows the rehearsal's.
	if got := path(Row{Stops: []Rehearsal{{Core: "2.541.3", Java: 17}, {Core: "2.580.1"}}}); got != "2.541.3 (Java 17) → 2.580.1" {
		t.Errorf("path = %q", got)
	}
}

func TestShare(t *testing.T) {
	for f, want := range map[float64]string{0: "–", 0.0004: "<0.1%", 0.001: "0.1%", 0.256: "25.6%"} {
		if got := share(f); got != want {
			t.Errorf("share(%v) = %q, want %q", f, got, want)
		}
	}
}

// A month that is no yyyymm is shown as it is, and a matrix where every
// baseline is clean says nothing about starting points.
func TestMarkdownOddMonth(t *testing.T) {
	m := &Matrix{To: "2.580.1", Month: "latest", Controllers: 5, Rows: []Row{{From: "2.541.3", Top: 10, Dir: "2.541.3-top10",
		Baseline: Rehearsal{Started: true, Plugins: 10}, AsIs: Rehearsal{Started: true, Plugins: 10}}}}
	md := m.Markdown()
	for _, want := range []string{"(Jenkins usage statistics, latest)", "Replacing only the core: Jenkins does not start from 0 of them, and starts with plugins that fail to load from 0 more; – of all",
		"every stop loads cleanly from 0 of 1.\n"} {
		if !strings.Contains(md, want) {
			t.Errorf("table lacks %q:\n%s", want, md)
		}
	}
	if strings.Contains(md, "starting points already have") {
		t.Error("clean baselines are reported")
	}
}
