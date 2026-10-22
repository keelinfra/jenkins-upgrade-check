// Command matrix does the bookkeeping around the upgrade matrix: which
// starting versions to rehearse, which month of usage statistics to use, and
// the table of every result. hack/experiment.sh does the rehearsals.
//
//	go run ./hack/matrix month
//	go run ./hack/matrix versions -to 2.580.1
//	go run ./hack/matrix table -results DIR -o matrix.md -json matrix.json
package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/keelinfra/jenkins-upgrade-check/internal/uc"
	"github.com/keelinfra/jenkins-upgrade-check/internal/version"
)

const (
	statsBase = "https://raw.githubusercontent.com/jenkins-infra/infra-statistics/gh-pages/jenkins-stats/svg"
	warIndex  = "https://repo.jenkins-ci.org/releases/org/jenkins-ci/main/jenkins-war/maven-metadata.xml"
	// firstLine is the oldest LTS line with a jenkins/jenkins image (2.60.3).
	firstLine = 60
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run is main without the process: it returns the exit code.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "usage: matrix month | versions -to VERSION | table -results DIR")
		return 2
	}
	ctx := context.Background()
	var err error
	switch args[0] {
	case "month":
		err = cmdMonth(ctx, stdout)
	case "versions":
		err = cmdVersions(ctx, args[1:], stdout)
	case "table":
		err = cmdTable(ctx, args[1:], stdout)
	default:
		err = fmt.Errorf("unknown command %q", args[0])
	}
	if err != nil {
		fmt.Fprintln(stderr, "matrix:", err)
		return 1
	}
	return 0
}

// client, updateCenter and now are variables so that tests can replace
// them.
var (
	client       = &http.Client{Timeout: 2 * time.Minute}
	updateCenter = &uc.Source{}
	now          = time.Now
)

func get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

// cmdMonth prints the newest month the usage statistics have a plugin
// ranking for. They are published with a delay, and sometimes not at all.
func cmdMonth(ctx context.Context, stdout io.Writer) error {
	now := now().UTC()
	for i := 0; i < 24; i++ {
		m := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, time.UTC).Format("200601")
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, statsBase+"/"+m+"-top-plugins1000.csv", nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			fmt.Fprintln(stdout, m)
			return nil
		}
	}
	return errors.New("no usage statistics in the last 24 months")
}

// cmdVersions prints, as a JSON array, the last release of every LTS line
// older than the target's line: where controllers that stopped updating sit.
func cmdVersions(ctx context.Context, args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("versions", flag.ContinueOnError)
	to := fs.String("to", "lts", "target version, or lts")
	min := fs.Int("min-line", firstLine, "oldest LTS line to include")
	if err := fs.Parse(args); err != nil {
		return err
	}
	target, err := updateCenter.ResolveCore(ctx, *to)
	if err != nil {
		return err
	}
	b, err := get(ctx, warIndex)
	if err != nil {
		return err
	}
	out := LastPatches(string(b), target, *min)
	if len(out) == 0 {
		return fmt.Errorf("no LTS line between 2.%d and %s", *min, target)
	}
	return json.NewEncoder(stdout).Encode(out)
}

var ltsVersion = regexp.MustCompile(`<version>2\.(\d+)\.(\d+)</version>`)

// LastPatches picks the newest 2.N.M release of every line N from min up to,
// not including, the target's line.
func LastPatches(metadata, target string, min int) []string {
	tl, _ := strconv.Atoi(strings.Split(version.Base(target), ".")[1])
	last := map[int]int{}
	for _, m := range ltsVersion.FindAllStringSubmatch(metadata, -1) {
		line, _ := strconv.Atoi(m[1])
		patch, _ := strconv.Atoi(m[2])
		if line < min || line >= tl {
			continue
		}
		if p, ok := last[line]; !ok || patch > p {
			last[line] = patch
		}
	}
	lines := make([]int, 0, len(last))
	for l := range last {
		lines = append(lines, l)
	}
	sort.Ints(lines)
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = fmt.Sprintf("2.%d.%d", l, last[l])
	}
	return out
}

// Rehearsal is one start of Jenkins, as hack/experiment.sh records it.
type Rehearsal struct {
	Core string `json:"core"`
	Java int    `json:"java,omitempty"`
	// PlanJava is the Java the plan switches to for this stop; the
	// rehearsal runs a newer one when a plugin needs it, and JavaReason
	// says which.
	PlanJava   int    `json:"planJava,omitempty"`
	JavaReason string `json:"javaReason,omitempty"`
	Image      string `json:"image,omitempty"`
	Started    bool   `json:"started"`
	Reason     string `json:"reason,omitempty"`
	Plugins    int    `json:"plugins"`
	Failed     int    `json:"failed"`
	// Broken plugins loaded, but some of their extensions did not.
	Broken      int      `json:"broken"`
	BrokenNames []string `json:"brokenNames,omitempty"`
	Skipped     int      `json:"skipped"`
	Roots       []Root   `json:"roots"`
}

// Root is one cause of failed plugins.
type Root struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Kind    string `json:"kind,omitempty"`
	Why     string `json:"why,omitempty"`
	Cascade int    `json:"cascade,omitempty"`
}

// Clean reports whether Jenkins started and loaded everything it was given.
func (r Rehearsal) Clean() bool { return r.Started && r.Failed == 0 && r.Broken == 0 }

// Row is one starting version and plugin count.
type Row struct {
	From         string      `json:"from"`
	To           string      `json:"to"`
	Top          int         `json:"top"`
	Month        string      `json:"month"`
	Baseline     Rehearsal   `json:"baseline"`
	AsIs         Rehearsal   `json:"asIs"`
	Stops        []Rehearsal `json:"stops"`
	PlanProblems int         `json:"planProblems"`
	// Share is the fraction of reporting controllers on the From line;
	// Dir is where the row's details are, relative to the table.
	Share float64 `json:"share"`
	Dir   string  `json:"dir"`
}

// PlannedClean reports whether every stop of the plan loaded everything.
func (r Row) PlannedClean() bool {
	if len(r.Stops) == 0 {
		return false
	}
	for _, s := range r.Stops {
		if !s.Clean() {
			return false
		}
	}
	return true
}

// Matrix is every row for one target.
type Matrix struct {
	To          string    `json:"to"`
	Month       string    `json:"month"`
	Generated   time.Time `json:"generated"`
	Controllers int       `json:"controllers"`
	Rows        []Row     `json:"rows"`
}

func cmdTable(ctx context.Context, args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("table", flag.ContinueOnError)
	dir := fs.String("results", "results", "directory with the row.json files")
	stats := fs.String("stats", "", "usage statistics CSV of Jenkins versions (default: download it for the rows' month)")
	out := fs.String("o", "", "write the Markdown table here (default: stdout)")
	jsonOut := fs.String("json", "", "also write the matrix as JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}
	m, err := Load(*dir)
	if err != nil {
		return err
	}
	var csvData []byte
	if *stats != "" {
		csvData, err = os.ReadFile(*stats)
	} else {
		csvData, err = get(ctx, statsBase+"/"+m.Month+"-jenkins.csv")
	}
	if err != nil {
		return err
	}
	if err := m.AddShares(string(csvData)); err != nil {
		return err
	}
	m.Generated = now().UTC()
	if *jsonOut != "" {
		b, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(*jsonOut, append(b, '\n'), 0o644); err != nil {
			return err
		}
	}
	w := stdout
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}
	_, err = io.WriteString(w, m.Markdown())
	return err
}

// Load reads every row.json under dir.
func Load(dir string) (*Matrix, error) {
	m := &Matrix{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != "row.json" {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		var r Row
		if err := json.Unmarshal(b, &r); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		r.Dir = fmt.Sprintf("%s-top%d", r.From, r.Top)
		m.Rows = append(m.Rows, r)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(m.Rows) == 0 {
		return nil, fmt.Errorf("no row.json under %s", dir)
	}
	sort.Slice(m.Rows, func(i, j int) bool {
		a, b := m.Rows[i], m.Rows[j]
		if a.Top != b.Top {
			return a.Top > b.Top
		}
		return version.Less(a.From, b.From)
	})
	m.To, m.Month = m.Rows[0].To, m.Rows[0].Month
	return m, nil
}

// AddShares sets each row's share of controllers from the usage statistics
// CSV of Jenkins versions ("version","installations").
func (m *Matrix) AddShares(data string) error {
	recs, err := csv.NewReader(strings.NewReader(data)).ReadAll()
	if err != nil {
		return fmt.Errorf("usage statistics: %w", err)
	}
	byLine := map[string]int{}
	total := 0
	for _, r := range recs {
		if len(r) < 2 {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(r[1]))
		if err != nil {
			continue
		}
		total += n
		if v := strings.TrimSpace(r[0]); strings.Count(v, ".") == 2 {
			byLine[version.Base(v)] += n
		}
	}
	if total == 0 {
		return errors.New("usage statistics: no installations")
	}
	m.Controllers = total
	for i := range m.Rows {
		m.Rows[i].Share = float64(byLine[version.Base(m.Rows[i].From)]) / float64(total)
	}
	return nil
}

// cell describes one rehearsal in a few words.
func cell(r Rehearsal) string {
	var s string
	switch {
	case !r.Started && r.Image == "" && r.Plugins == 0:
		return "— not rehearsed: " + r.Reason
	case !r.Started:
		s = "❌ Jenkins does not start"
		if names := rootNames(r.Roots); names != "" {
			s += ": " + names
		} else if why := shortReason(r.Reason); why != "" {
			s += " (" + why + ")"
		}
	case r.Failed > 0:
		s = fmt.Sprintf("⚠️ %d of %d fail: %s", r.Failed, r.Plugins, rootNames(r.Roots))
	case r.Broken > 0:
		s = fmt.Sprintf("⚠️ %d of %d load in part", r.Broken, r.Plugins)
		if names := firstNames(r.BrokenNames); names != "" {
			s += ": " + names
		}
	default:
		s = fmt.Sprintf("✅ all %d load", r.Plugins)
	}
	if r.Skipped > 0 {
		s += fmt.Sprintf(" (%d not downloadable)", r.Skipped)
	}
	return s
}

// shortReason keeps the part of a reason that says why Jenkins stopped.
func shortReason(reason string) string {
	if _, after, ok := strings.Cut(reason, "; the log ends: "); ok {
		reason = after
	}
	reason = strings.ReplaceAll(reason, "|", "/")
	if r := []rune(reason); len(r) > 120 {
		reason = string(r[:120]) + "…"
	}
	return reason
}

func rootNames(roots []Root) string {
	names := make([]string, len(roots))
	for i, r := range roots {
		names[i] = r.Name
	}
	return firstNames(names)
}

// firstNames lists three names and counts the rest.
func firstNames(names []string) string {
	if len(names) > 3 {
		return strings.Join(names[:3], ", ") + fmt.Sprintf(", and %d more", len(names)-3)
	}
	return strings.Join(names, ", ")
}

func path(r Row) string {
	var hops []string
	for _, s := range r.Stops {
		h := s.Core
		switch {
		case s.JavaReason != "" && s.Java > 0:
			// The plan does not know every plugin's Java; the rehearsal does.
			h += fmt.Sprintf(" (Java %d: %s)", s.Java, s.JavaReason)
		case s.PlanJava > 0:
			h += fmt.Sprintf(" (Java %d)", s.PlanJava)
		case s.Java > 0:
			h += fmt.Sprintf(" (Java %d)", s.Java)
		}
		hops = append(hops, h)
	}
	return strings.Join(hops, " → ")
}

func planned(r Row) string {
	if len(r.Stops) == 0 {
		return "— no plan"
	}
	for _, s := range r.Stops {
		if !s.Clean() {
			c := cell(s)
			if len(r.Stops) > 1 {
				c += " (at " + s.Core + ")"
			}
			return c
		}
	}
	return cell(r.Stops[len(r.Stops)-1])
}

func share(f float64) string {
	if f == 0 {
		return "–"
	}
	if f < 0.001 {
		return "<0.1%"
	}
	return fmt.Sprintf("%.1f%%", f*100)
}

// Markdown renders the matrix with a summary for each plugin count.
func (m *Matrix) Markdown() string {
	b := &strings.Builder{}
	month := m.Month
	if t, err := time.Parse("200601", m.Month); err == nil {
		month = t.Format("January 2006")
	}
	fmt.Fprintf(b, "# Upgrading every Jenkins LTS line to %s\n\n", m.To)
	fmt.Fprintf(b, "Rehearsed %s with [jenkins-upgrade-check](https://github.com/keelinfra/jenkins-upgrade-check). "+
		"Each row is a controller on the last release of an LTS line with the most installed plugins "+
		"(Jenkins usage statistics, %s), each at the newest release that line can load. "+
		"Every rehearsal starts Jenkins in a container with no network and records which plugins fail to load.\n\n",
		m.Generated.Format("2006-01-02"), month)
	b.WriteString("- **As it is**: the starting controller itself.\n")
	fmt.Fprintf(b, "- **Core only**: the same plugin files on %s, as when only the war, the package or the image is replaced.\n", m.To)
	b.WriteString("- **Planned**: the path `jenkins-upgrade-check plan` gives, which switches Java where no Java runs both versions, " +
		"with the plugin versions it computes for each stop. Every stop is rehearsed; the cell shows the first that does not load cleanly.\n")
	fmt.Fprintf(b, "- **Controllers**: share of the %d controllers that reported usage in %s running a release of that line.\n\n", m.Controllers, month)

	tops := []int{}
	byTop := map[int][]Row{}
	for _, r := range m.Rows {
		if _, ok := byTop[r.Top]; !ok {
			tops = append(tops, r.Top)
		}
		byTop[r.Top] = append(byTop[r.Top], r)
	}

	b.WriteString("## Summary\n\n")
	for _, top := range tops {
		rows := byTop[top]
		var noStart, fails, plannedOK, baseBad int
		var breakShare float64
		for _, r := range rows {
			switch {
			case !r.AsIs.Started:
				noStart++
				breakShare += r.Share
			case !r.AsIs.Clean():
				fails++
				breakShare += r.Share
			}
			if r.PlannedClean() {
				plannedOK++
			}
			if !r.Baseline.Clean() {
				baseBad++
			}
		}
		fmt.Fprintf(b, "**%d most installed plugins, %d starting versions.** ", top, len(rows))
		fmt.Fprintf(b, "Replacing only the core: Jenkins does not start from %d of them, and starts with plugins that fail to load from %d more; "+
			"%s of all reporting controllers are on those lines. ", noStart, fails, share(breakShare))
		fmt.Fprintf(b, "Following the plan: every stop loads cleanly from %d of %d.", plannedOK, len(rows))
		if baseBad > 0 {
			fmt.Fprintf(b, " %d starting points already have plugins that do not load before any upgrade; see their rows.", baseBad)
		}
		b.WriteString("\n\n")
	}

	for _, top := range tops {
		fmt.Fprintf(b, "## %d most installed plugins\n\n", top)
		fmt.Fprintf(b, "| From | Controllers | As it is | Core only → %s | Planned path | Planned |\n", m.To)
		b.WriteString("|---|---|---|---|---|---|\n")
		for _, r := range byTop[top] {
			fmt.Fprintf(b, "| [%s](%s/summary.md) | %s | %s | %s | %s | %s |\n",
				r.From, r.Dir, share(r.Share), cell(r.Baseline), cell(r.AsIs), path(r), planned(r))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Limits\n\n")
	b.WriteString("- Plugins are at the newest release each line can load. A controller whose plugins are older than that has more to change, not less.\n")
	b.WriteString("- There is no real configuration and no network: what plugins do with jobs, credentials and outside systems is not exercised. " +
		"`jenkins-upgrade-check dry-run --home` rehearses your own configuration.\n")
	b.WriteString("- Core only is rehearsed as an upgrade with no update center, as on a controller in a closed network. " +
		"Jenkins 2.574 and 2.577 stopped carrying the plugins split out of core (sshd, javax-mail-api, instance-identity and others) in the war, " +
		"so from 2.580.1 on such a controller no longer gets them on upgrade; one that reaches the update center may fetch some of them itself. " +
		"The planned path adds them where its plugins still expect them.\n")
	b.WriteString("- A plugin that loads can still fail at run time, in a pipeline step or a page. The rehearsal shows what fails at startup.\n")
	b.WriteString("- The usage statistics count controllers that send them; controllers inside closed networks usually do not.\n")
	return b.String()
}
