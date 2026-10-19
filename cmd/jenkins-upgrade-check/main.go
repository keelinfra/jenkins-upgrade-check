// Command jenkins-upgrade-check plans and rehearses a Jenkins upgrade.
package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/keelinfra/jenkins-upgrade-check/internal/dryrun"
	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
	"github.com/keelinfra/jenkins-upgrade-check/internal/report"
	"github.com/keelinfra/jenkins-upgrade-check/internal/snapshot"
	"github.com/keelinfra/jenkins-upgrade-check/internal/uc"
)

const tool = "jenkins-upgrade-check"

// Version is set at build time.
var Version = "dev"

const usage = `jenkins-upgrade-check plans and rehearses a Jenkins upgrade.

Usage:
  jenkins-upgrade-check collect  --url URL --user USER --token TOKEN [-o snapshot.json]
  jenkins-upgrade-check plan     (--snapshot snapshot.json | --plugins-file FILE --core VERSION) [--to lts]
  jenkins-upgrade-check dry-run  (--plan plan.json | --plugins-file FILE --core VERSION [--resolve])
  jenkins-upgrade-check resolve  --core VERSION (--plugins-file FILE | --names a,b | --top-plugins CSV --top N)
  jenkins-upgrade-check version

collect   reads a controller with GET requests only: core version, plugins, agents.
plan      works out the path, the plugin versions for each step, and what the
          update center already knows is wrong (security warnings, deprecations,
          plugins no longer distributed, agents that cannot follow).
dry-run   starts the target Jenkins in Docker with your plugin set and reports
          which plugins fail to load. Only the plugin list goes into the container.
resolve   prints a consistent plugin set for a core version (plugins.txt format).

Run "jenkins-upgrade-check COMMAND -h" for the flags of one command.
`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	var err error
	code := 0
	switch os.Args[1] {
	case "collect":
		err = cmdCollect(ctx, os.Args[2:])
	case "plan":
		err = cmdPlan(ctx, os.Args[2:])
	case "dry-run", "dryrun":
		code, err = cmdDryRun(ctx, os.Args[2:])
	case "resolve":
		err = cmdResolve(ctx, os.Args[2:])
	case "version", "--version":
		fmt.Println(tool, Version)
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
	if err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
		os.Exit(2)
	}
	os.Exit(code)
}

type sourceFlags struct {
	dir    string
	base   string
	cache  string
	maxAge time.Duration
}

func (s *sourceFlags) register(fs *flag.FlagSet) {
	fs.StringVar(&s.dir, "update-center-dir", "", "read update-center.actual.json and plugin-versions.json from this directory instead of the network")
	fs.StringVar(&s.base, "update-center", uc.DefaultBase, "update center root URL")
	def := ""
	if d, err := os.UserCacheDir(); err == nil {
		def = filepath.Join(d, tool)
	}
	fs.StringVar(&s.cache, "cache-dir", def, "where downloads are kept")
	fs.DurationVar(&s.maxAge, "cache-max-age", 12*time.Hour, "re-download update center data older than this")
}

func (s *sourceFlags) source() *uc.Source {
	return &uc.Source{Base: s.base, Dir: s.dir, CacheDir: s.cache, MaxAge: s.maxAge, UserAgent: tool + "/" + Version}
}

func load(ctx context.Context, src *uc.Source) (*uc.Releases, *uc.UpdateCenter, error) {
	fmt.Fprintln(os.Stderr, "reading update center data")
	rel, err := src.Releases(ctx)
	if err != nil {
		return nil, nil, err
	}
	center, err := src.UpdateCenter(ctx)
	if err != nil {
		return nil, nil, err
	}
	return rel, center, nil
}

func cmdCollect(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("collect", flag.ContinueOnError)
	c := &snapshot.Collector{Tool: tool, ToolVersion: Version}
	fs.StringVar(&c.URL, "url", os.Getenv("JENKINS_URL"), "controller URL (or JENKINS_URL)")
	fs.StringVar(&c.User, "user", os.Getenv("JENKINS_USER"), "user name (or JENKINS_USER)")
	fs.StringVar(&c.Token, "token", "", "API token (better: JENKINS_TOKEN, which stays out of shell history)")
	fs.StringVar(&c.CAFile, "ca-file", "", "extra CA certificates (PEM) to trust, for an internal CA")
	fs.BoolVar(&c.AllowHTTP, "allow-http", false, "allow a plain http URL that is not localhost")
	out := fs.String("o", "snapshot.json", "where to write the snapshot")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if c.Token == "" {
		c.Token = os.Getenv("JENKINS_TOKEN")
	}
	if c.URL == "" {
		return errors.New("--url is required")
	}
	s, err := c.Collect(ctx)
	if err != nil {
		return err
	}
	if err := s.Write(*out); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Jenkins %s: %d plugins, %d agents → %s (%d GET requests, nothing written to the controller)\n", s.Core, len(s.Plugins), len(s.Agents), *out, c.Requests())
	for _, n := range s.Notes {
		fmt.Fprintln(os.Stderr, "note:", n)
	}
	return nil
}

type inputFlags struct {
	snapshot    string
	pluginsFile string
	core        string
}

func (in *inputFlags) register(fs *flag.FlagSet) {
	fs.StringVar(&in.snapshot, "snapshot", "", "snapshot written by collect")
	fs.StringVar(&in.pluginsFile, "plugins-file", "", "plugins.txt (name:version per line) or saved /pluginManager/api/json")
	fs.StringVar(&in.core, "core", "", "Jenkins version the plugin list is from (with --plugins-file)")
}

func (in *inputFlags) load() (*snapshot.Snapshot, error) {
	switch {
	case in.snapshot != "" && in.pluginsFile != "":
		return nil, errors.New("give --snapshot or --plugins-file, not both")
	case in.snapshot != "":
		s, err := snapshot.Read(in.snapshot)
		if err != nil {
			return nil, err
		}
		if in.core != "" {
			s.Core = in.core
		}
		return s, nil
	case in.pluginsFile != "":
		if in.core == "" {
			return nil, errors.New("--plugins-file needs --core, the Jenkins version those plugins run on")
		}
		return snapshot.FromPluginList(in.pluginsFile, in.core)
	}
	return nil, errors.New("give --snapshot (from collect) or --plugins-file with --core")
}

func cmdPlan(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("plan", flag.ContinueOnError)
	var in inputFlags
	in.register(fs)
	var sf sourceFlags
	sf.register(fs)
	to := fs.String("to", "lts", `target: a version, "lts" or "weekly"`)
	out := fs.String("o", "", "write the Markdown report here (default: stdout)")
	jsonOut := fs.String("json", "", "also write the plan as JSON (dry-run --plan reads it)")
	pluginsOut := fs.String("plugins-out", "", "also write the target plugin set as plugins.txt")
	if err := fs.Parse(args); err != nil {
		return err
	}
	s, err := in.load()
	if err != nil {
		return err
	}
	src := sf.source()
	target, err := src.ResolveCore(ctx, *to)
	if err != nil {
		return err
	}
	rel, center, err := load(ctx, src)
	if err != nil {
		return err
	}
	p, err := plan.Build(s, target, rel, center)
	if err != nil {
		return err
	}
	p.Tool, p.ToolVersion = tool, Version
	if *jsonOut != "" {
		if err := writeJSON(*jsonOut, p); err != nil {
			return err
		}
	}
	if *pluginsOut != "" {
		if err := writePicks(*pluginsOut, p.Target.Plugins); err != nil {
			return err
		}
	}
	return withOutput(*out, func(w io.Writer) error { return report.Plan(w, p) })
}

func cmdDryRun(ctx context.Context, args []string) (int, error) {
	fs := flag.NewFlagSet("dry-run", flag.ContinueOnError)
	var in inputFlags
	in.register(fs)
	var sf sourceFlags
	sf.register(fs)
	planFile := fs.String("plan", "", "plan JSON written by plan --json: rehearse its target")
	resolve := fs.Bool("resolve", false, "with --plugins-file: move each plugin to the newest release --core can load first (default: use the versions as given)")
	java := fs.Int("java", 0, "Java version of the image (default: the oldest the core supports)")
	image := fs.String("image", "", "base image instead of jenkins/jenkins:<core>-jdk<N>")
	outDir := fs.String("out-dir", "", "keep the Jenkins log here")
	timeout := fs.Duration("timeout", 10*time.Minute, "how long Jenkins may take to start")
	keep := fs.Bool("keep", false, "keep the container, image and work directory")
	out := fs.String("o", "", "write the Markdown report here (default: stdout)")
	jsonOut := fs.String("json", "", "also write the result as JSON")
	if err := fs.Parse(args); err != nil {
		return 0, err
	}
	o := dryrun.Options{Java: *java, Image: *image, OutDir: *outDir, Timeout: *timeout, Keep: *keep, Progress: os.Stderr}
	if sf.cache != "" {
		o.CacheDir = filepath.Join(sf.cache, "plugins")
	}
	switch {
	case *planFile != "":
		var p plan.Plan
		b, err := os.ReadFile(*planFile)
		if err != nil {
			return 0, err
		}
		if err := json.Unmarshal(b, &p); err != nil {
			return 0, fmt.Errorf("%s: %w", *planFile, err)
		}
		o.Core = p.Target.Core
		o.Plugins = p.Target.Plugins
		if len(p.Steps) > 0 && o.Java == 0 {
			o.Java = p.Steps[len(p.Steps)-1].Java
		}
	default:
		s, err := in.load()
		if err != nil {
			return 0, err
		}
		rel, err := sf.source().Releases(ctx)
		if err != nil {
			return 0, err
		}
		picks := make([]plan.Pick, 0, len(s.Plugins))
		for _, pl := range s.Plugins {
			picks = append(picks, plan.Pick{Name: pl.Name, Version: pl.Version})
		}
		if *resolve {
			r := plan.Resolve(rel, picks, s.Core)
			for _, pr := range r.Problems {
				fmt.Fprintln(os.Stderr, "warning:", pr)
			}
			picks = r.Plugins
		} else {
			picks = plan.Locate(rel, picks)
		}
		o.Core = s.Core
		o.Plugins = picks
	}
	res, err := dryrun.Run(ctx, o)
	if err != nil {
		return 0, err
	}
	if *jsonOut != "" {
		if err := writeJSON(*jsonOut, res); err != nil {
			return 0, err
		}
	}
	if err := withOutput(*out, func(w io.Writer) error { return report.Rehearsal(w, res) }); err != nil {
		return 0, err
	}
	if !res.OK() {
		return 1, nil
	}
	return 0, nil
}

func cmdResolve(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("resolve", flag.ContinueOnError)
	var sf sourceFlags
	sf.register(fs)
	core := fs.String("core", "", "Jenkins version to resolve for")
	pluginsFile := fs.String("plugins-file", "", "plugins.txt; versions are a starting point and may be raised")
	names := fs.String("names", "", "comma-separated plugin names")
	topFile := fs.String("top-plugins", "", `CSV of "name","installations" (Jenkins usage statistics), most installed first or in any order`)
	top := fs.Int("top", 100, "with --top-plugins: how many of the most installed to take")
	out := fs.String("o", "", "write plugins.txt here (default: stdout)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *core == "" {
		return errors.New("--core is required")
	}
	src := sf.source()
	target, err := src.ResolveCore(ctx, *core)
	if err != nil {
		return err
	}
	rel, err := src.Releases(ctx)
	if err != nil {
		return err
	}
	var picks []plan.Pick
	switch {
	case *pluginsFile != "":
		s, err := snapshot.FromPluginList(*pluginsFile, target)
		if err != nil {
			return err
		}
		for _, p := range s.Plugins {
			picks = append(picks, plan.Pick{Name: p.Name, Version: p.Version})
		}
	case *names != "":
		for _, n := range strings.Split(*names, ",") {
			if n = strings.TrimSpace(n); n != "" {
				picks = append(picks, plan.Pick{Name: n, Version: "0"})
			}
		}
	case *topFile != "":
		list, err := topPlugins(*topFile, *top, rel, target)
		if err != nil {
			return err
		}
		for _, n := range list {
			picks = append(picks, plan.Pick{Name: n, Version: "0"})
		}
	default:
		return errors.New("give --plugins-file, --names or --top-plugins")
	}
	r := plan.Resolve(rel, picks, target)
	for _, pr := range r.Problems {
		fmt.Fprintln(os.Stderr, "warning:", pr)
	}
	fmt.Fprintf(os.Stderr, "%d plugins for Jenkins %s\n", len(r.Plugins), target)
	return withOutput(*out, func(w io.Writer) error { return picksTxt(w, r.Plugins) })
}

// topPlugins reads a usage statistics CSV and returns the most installed
// plugins that have a release for the core.
func topPlugins(path string, n int, rel *uc.Releases, core string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows, err := csv.NewReader(bufio.NewReader(f)).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	type row struct {
		name  string
		count int
	}
	var all []row
	for _, r := range rows {
		if len(r) < 2 {
			continue
		}
		c, err := strconv.Atoi(strings.TrimSpace(r[1]))
		if err != nil {
			continue
		}
		all = append(all, row{strings.TrimSpace(r[0]), c})
	}
	// Most installed first.
	for i := 1; i < len(all); i++ {
		for j := i; j > 0 && all[j].count > all[j-1].count; j-- {
			all[j], all[j-1] = all[j-1], all[j]
		}
	}
	var out []string
	for _, r := range all {
		if len(out) >= n {
			break
		}
		if rel.LatestFor(r.name, core) == nil {
			continue
		}
		out = append(out, r.name)
	}
	return out, nil
}

func picksTxt(w io.Writer, picks []plan.Pick) error {
	for _, p := range picks {
		if _, err := fmt.Fprintf(w, "%s:%s\n", p.Name, p.Version); err != nil {
			return err
		}
	}
	return nil
}

func writePicks(path string, picks []plan.Pick) error {
	return withOutput(path, func(w io.Writer) error { return picksTxt(w, picks) })
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func withOutput(path string, fn func(io.Writer) error) error {
	if path == "" || path == "-" {
		return fn(os.Stdout)
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := fn(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
