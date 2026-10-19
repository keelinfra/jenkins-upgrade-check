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
	"github.com/keelinfra/jenkins-upgrade-check/internal/runbook"
	"github.com/keelinfra/jenkins-upgrade-check/internal/snapshot"
	"github.com/keelinfra/jenkins-upgrade-check/internal/uc"
	"github.com/keelinfra/jenkins-upgrade-check/internal/verify"
)

const tool = "jenkins-upgrade-check"

// Version is set at build time.
var Version = "dev"

const usage = `jenkins-upgrade-check plans, rehearses and verifies a Jenkins upgrade.

Usage:
  jenkins-upgrade-check collect  --url URL --user USER [-o snapshot.json]
  jenkins-upgrade-check plan     (--snapshot snapshot.json | --plugins-file FILE --core VERSION) [--to lts]
                                 [--runbook runbook.md --install deb|rpm|war|docker]
  jenkins-upgrade-check dry-run  (--plan plan.json | --plugins-file FILE --core VERSION [--resolve]) [--home JENKINS_HOME]
  jenkins-upgrade-check download --plan plan.json [--step N] --dir DIR
  jenkins-upgrade-check probe    [--quiet]
  jenkins-upgrade-check verify   --before before.json --url URL --user USER [--plan plan.json] [--core VERSION]
                                 [--home JENKINS_HOME] [--log FILE]
  jenkins-upgrade-check resolve  --core VERSION (--plugins-file FILE | --names a,b | --top-plugins CSV --top N)
  jenkins-upgrade-check version

collect   reads a controller with GET requests only: core version, plugins, agents, jobs.
plan      works out the path, the plugin versions for each step, and what the
          update center already knows is wrong (security warnings, deprecations,
          plugins no longer distributed, agents that cannot follow). With
          --runbook it also writes the commands for the upgrade and the rollback.
dry-run   starts the target Jenkins in Docker, without network, with your plugin
          set and optionally a copy of your configuration (--home), and says what
          failed, why, and what to do about it.
download  fetches the plugins of one plan step, checking checksums, before the window.
probe     prints the Groovy init script that records what Jenkins loaded as it
          starts; with --quiet it also starts Jenkins in quiet mode (no builds).
verify    compares the upgraded controller with the snapshot from before and says
          GO or ROLL BACK, and why (exit 0 or 1).
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
	case "download":
		err = cmdDownload(ctx, os.Args[2:])
	case "probe":
		err = cmdProbe(os.Args[2:])
	case "verify":
		code, err = cmdVerify(ctx, os.Args[2:])
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
	rb := fs.String("runbook", "", "also write the upgrade and rollback commands here (Markdown)")
	install := fs.String("install", "", "with --runbook: how Jenkins is installed: "+strings.Join(runbook.Installs, ", "))
	home := fs.String("home", "", "with --runbook: JENKINS_HOME on the controller host (default: the usual place for --install)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *rb != "" && *install == "" {
		return fmt.Errorf("--runbook needs --install, one of %s", strings.Join(runbook.Installs, ", "))
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
	if *rb != "" {
		url := ""
		if strings.HasPrefix(s.Source, "http") {
			url = s.Source
		}
		planFile := *jsonOut
		if planFile == "" {
			planFile = "plan.json"
		}
		o := runbook.Options{Install: *install, Home: *home, URL: url, PlanFile: planFile}
		if err := withOutput(*rb, func(w io.Writer) error { return runbook.Write(w, p, o) }); err != nil {
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
	home := fs.String("home", "", "JENKINS_HOME (or a copy of it) whose configuration to rehearse: jobs, folders, nodes, credentials and settings, without build records")
	casc := fs.String("casc", "", "Configuration as Code file or directory to load, if the controller reads it from outside JENKINS_HOME")
	java := fs.Int("java", 0, "Java version of the image (default: the oldest the core and the plugins support)")
	image := fs.String("image", "", "base image instead of jenkins/jenkins:<core>-jdk<N>")
	outDir := fs.String("out-dir", "", "keep the Jenkins log here")
	timeout := fs.Duration("timeout", 10*time.Minute, "how long Jenkins may take to start")
	keep := fs.Bool("keep", false, "keep the container and the work directory")
	out := fs.String("o", "", "write the Markdown report here (default: stdout)")
	jsonOut := fs.String("json", "", "also write the result as JSON")
	if err := fs.Parse(args); err != nil {
		return 0, err
	}
	o := dryrun.Options{Home: *home, CasC: *casc, Java: *java, Image: *image, OutDir: *outDir, Timeout: *timeout, Keep: *keep, Progress: os.Stderr}
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

func readPlan(path string) (*plan.Plan, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p plan.Plan
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &p, nil
}

func cmdDownload(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("download", flag.ContinueOnError)
	var sf sourceFlags
	sf.register(fs)
	planFile := fs.String("plan", "", "plan JSON written by plan --json")
	step := fs.Int("step", -1, "which step's plugin set (0 is the plugin update on the current core; default: the last)")
	dir := fs.String("dir", "", "where to put the .jpi files")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *planFile == "" || *dir == "" {
		return errors.New("--plan and --dir are required")
	}
	p, err := readPlan(*planFile)
	if err != nil {
		return err
	}
	picks := p.Target.Plugins
	core := p.Target.Core
	if *step >= 0 {
		if *step >= len(p.Steps) {
			return fmt.Errorf("the plan has steps 0 to %d", len(p.Steps)-1)
		}
		picks, core = p.Steps[*step].Plugins, p.Steps[*step].Core
		if len(picks) == 0 {
			return fmt.Errorf("step %d has no plugin set; make the plan again with this version", *step)
		}
	}
	if err := os.MkdirAll(*dir, 0o755); err != nil {
		return err
	}
	cache := ""
	if sf.cache != "" {
		cache = filepath.Join(sf.cache, "plugins")
	}
	fmt.Fprintf(os.Stderr, "downloading %d plugins for Jenkins %s into %s\n", len(picks), core, *dir)
	got, skipped := dryrun.Download(ctx, picks, *dir, cache)
	for _, sk := range skipped {
		fmt.Fprintln(os.Stderr, "not downloaded:", sk)
	}
	fmt.Fprintf(os.Stderr, "%d of %d plugins downloaded\n", len(got), len(picks))
	if len(skipped) > 0 {
		return fmt.Errorf("%d plugin(s) could not be downloaded; the plugins directory would be incomplete", len(skipped))
	}
	return nil
}

func cmdProbe(args []string) error {
	fs := flag.NewFlagSet("probe", flag.ContinueOnError)
	quiet := fs.Bool("quiet", false, "also put Jenkins into quiet mode as it starts, so no build runs until it is verified")
	if err := fs.Parse(args); err != nil {
		return err
	}
	_, err := os.Stdout.Write(dryrun.ProbeScript(*quiet))
	return err
}

func cmdVerify(ctx context.Context, args []string) (int, error) {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	c := &snapshot.Collector{Tool: tool, ToolVersion: Version}
	fs.StringVar(&c.URL, "url", os.Getenv("JENKINS_URL"), "controller URL (or JENKINS_URL)")
	fs.StringVar(&c.User, "user", os.Getenv("JENKINS_USER"), "user name (or JENKINS_USER); the token comes from JENKINS_TOKEN")
	fs.StringVar(&c.CAFile, "ca-file", "", "extra CA certificates (PEM) to trust")
	fs.BoolVar(&c.AllowHTTP, "allow-http", false, "allow a plain http URL that is not localhost")
	before := fs.String("before", "", "snapshot collected before the upgrade")
	planFile := fs.String("plan", "", "plan JSON: the plugin versions the upgrade was meant to install")
	core := fs.String("core", "", "the Jenkins version the upgrade was meant to reach (default: the plan's target)")
	home := fs.String("home", "", "JENKINS_HOME, to read what the probe script recorded and name the plugin a broken job needs")
	logFile := fs.String("log", "", "the Jenkins log of this start, to say why plugins did not load")
	out := fs.String("o", "", "write the Markdown report here (default: stdout)")
	jsonOut := fs.String("json", "", "also write the result as JSON")
	if err := fs.Parse(args); err != nil {
		return 0, err
	}
	c.Token = os.Getenv("JENKINS_TOKEN")
	if *before == "" || c.URL == "" {
		return 0, errors.New("--before and --url are required")
	}
	b, err := snapshot.Read(*before)
	if err != nil {
		return 0, err
	}
	in := verify.Input{Before: b, Core: *core, Home: *home}
	if *planFile != "" {
		p, err := readPlan(*planFile)
		if err != nil {
			return 0, err
		}
		if in.Core == "" {
			in.Core = p.To
		}
		in.Plugins = p.Target.Plugins
		for i := len(p.Steps) - 1; i >= 0; i-- {
			if p.Steps[i].Core == in.Core && len(p.Steps[i].Plugins) > 0 {
				in.Plugins = p.Steps[i].Plugins
				break
			}
		}
	}
	if *home != "" {
		pr, err := dryrun.ReadProbe(filepath.Join(*home, dryrun.ProbeFile))
		switch {
		case err != nil:
			fmt.Fprintf(os.Stderr, "note: no probe record in %s (%v); install it with: jenkins-upgrade-check probe --quiet > %s\n", *home, err, filepath.Join(*home, "init.groovy.d", "zz-jenkins-upgrade-check.groovy"))
		case in.Core != "" && pr.Version != in.Core:
			fmt.Fprintf(os.Stderr, "note: the probe record is from Jenkins %s, not %s; ignoring it\n", pr.Version, in.Core)
		default:
			in.Probe = pr
		}
	}
	if *logFile != "" {
		lb, err := os.ReadFile(*logFile)
		if err != nil {
			return 0, err
		}
		in.Log = string(lb)
	}
	after, err := c.Collect(ctx)
	switch {
	case errors.Is(err, snapshot.ErrRefused):
		return 0, err
	case err != nil:
		// A controller that does not answer after an upgrade is a verdict,
		// not a usage error: it did not start.
		in.Down = err.Error()
	}
	in.After = after
	if after != nil && in.Before != nil {
		// Why an agent did not come back is in its connection log.
		online := map[string]bool{}
		for _, a := range in.Before.Agents {
			online[a.Name] = !a.Offline
		}
		for _, a := range after.Agents {
			if online[a.Name] && a.Offline {
				if l, err := c.AgentLog(ctx, a.Name); err == nil {
					if in.AgentLogs == nil {
						in.AgentLogs = map[string]string{}
					}
					in.AgentLogs[a.Name] = l
				}
			}
		}
	}
	res := verify.Run(in)
	if *jsonOut != "" {
		if err := writeJSON(*jsonOut, res); err != nil {
			return 0, err
		}
	}
	if err := withOutput(*out, func(w io.Writer) error { return report.Verify(w, res) }); err != nil {
		return 0, err
	}
	fmt.Fprintf(os.Stderr, "%s\n", res.Verdict)
	if res.Verdict != verify.Go {
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
