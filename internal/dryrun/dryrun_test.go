package dryrun

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/keelinfra/jenkins-upgrade-check/internal/plan"
)

func TestParseLog(t *testing.T) {
	b, err := os.ReadFile("testdata/failed.log")
	if err != nil {
		t.Fatal(err)
	}
	failed, errs := ParseLog(string(b))
	want := []Failure{
		{Name: "workflow-cps", Version: "2.90", Reasons: []string{"Update required: Script Security Plugin (script-security 1.75) to be updated to 1336.vf33a_a_9863911 or higher"}},
		{Name: "workflow-aggregator", Version: "2.6", Reasons: []string{"Failed to load: Pipeline: Groovy (workflow-cps 2.90)"}},
	}
	if !reflect.DeepEqual(failed, want) {
		t.Errorf("failed =\n%+v\nwant\n%+v", failed, want)
	}
	if len(errs) != 1 || !strings.Contains(errs[0], "NoClassDefFoundError: javax/servlet/ServletException") {
		t.Errorf("errors = %q, want the NoClassDefFoundError and not the optional component", errs)
	}
}

func TestVerify(t *testing.T) {
	f := t.TempDir() + "/x"
	if err := os.WriteFile(f, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	// sha256("hello"), as the update center (base64) and as hex.
	for _, sum := range []string{"LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ=", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", ""} {
		if err := verify(f, sum); err != nil {
			t.Errorf("verify(%q): %v", sum, err)
		}
	}
	if err := verify(f, "AAAA"); err == nil {
		t.Error("a wrong checksum must fail")
	}
}

// A real 2.580.1 log: newer Jenkins prefixes continuation lines with [LF]>.
func TestParseLogNewFormat(t *testing.T) {
	b, err := os.ReadFile("testdata/failed-2.580.1.log")
	if err != nil {
		t.Fatal(err)
	}
	failed, errs := ParseLog(string(b))
	want := []Failure{
		{Name: "fake-bad-class", Version: "1.0", Reasons: []string{"java.lang.ClassNotFoundException: org.example.Missing", "java.io.IOException: Unable to load org.example.Missing from fake-bad-class"}},
		{Name: "fake-missing-dep", Version: "1.0", Reasons: []string{"Plugin is missing: does-not-exist (1.0)"}},
		{Name: "fake-needs-broken", Version: "1.0", Reasons: []string{"Failed to load: Fake fake-missing-dep (fake-missing-dep 1.0)"}},
		{Name: "fake-new-core", Version: "1.0", Reasons: []string{"Jenkins (9.999) or higher required"}},
	}
	if !reflect.DeepEqual(failed, want) {
		t.Errorf("failed =\n%+v\nwant\n%+v", failed, want)
	}
	rest, items := ItemErrors(errs)
	if items["pipeline-job"] != "com.thoughtworks.xstream.mapper.CannotResolveClassException: flow-definition" {
		t.Errorf("item errors = %q", items)
	}
	if len(rest) != 0 {
		t.Errorf("other errors = %q", rest)
	}
}

func TestRoots(t *testing.T) {
	failed := []Failure{
		{Name: "a", Reasons: []string{"Plugin is missing: gone (1.0)"}},
		{Name: "b", Reasons: []string{"Failed to load: Plugin A (a 1.0)"}},
		{Name: "c", Reasons: []string{"Failed to load: Plugin B (b 2.0)"}},
		{Name: "d", Reasons: []string{"java.lang.NoClassDefFoundError: javax/servlet/http/HttpServletRequest"}},
		{Name: "e", Reasons: []string{"Update required: Script Security (script-security 1.75) to be updated to 1336.v1 or higher"}},
		{Name: "f", Reasons: []string{"java.lang.UnsupportedClassVersionError: x has been compiled by a more recent version of the Java Runtime (class file version 61.0)"}},
		{Name: "g", Reasons: []string{"Java (17) or higher required", "Failed to load: Plugin A (a 1.0)"}},
	}
	roots := Roots(failed)
	got := map[string]Root{}
	for _, r := range roots {
		got[r.Name] = r
	}
	if roots[0].Name != "a" || !reflect.DeepEqual(roots[0].Cascade, []string{"b", "c"}) {
		t.Errorf("first root = %+v, want a with cascade b, c", roots[0])
	}
	for name, kind := range map[string]string{"a": CauseMissingDependency, "d": CauseJakartaEE9, "e": CauseOldDependency, "f": CauseNewerJava, "g": CauseJavaTooOld} {
		if got[name].Diagnosis.Kind != kind {
			t.Errorf("%s: kind %q, want %q (%s)", name, got[name].Diagnosis.Kind, kind, got[name].Diagnosis.Why)
		}
	}
	if !strings.Contains(got["f"].Diagnosis.Why, "Java 17") {
		t.Errorf("class file 61 is Java 17: %s", got["f"].Diagnosis.Why)
	}
	if _, ok := got["b"]; ok {
		t.Error("b failed only because a did; it is not a root")
	}
}

func TestProbeReasons(t *testing.T) {
	got := probeReasons([]string{"java.io.IOException: Failed to load: X (x 1)\n - Plugin is missing: y (2)\n - Jenkins (2.500) or higher required"})
	if !reflect.DeepEqual(got, []string{"Plugin is missing: y (2)", "Jenkins (2.500) or higher required"}) {
		t.Errorf("dependency reasons = %q", got)
	}
	got = probeReasons([]string{"java.io.IOException: Failed to load: X", "java.lang.ClassNotFoundException: org.Foo"})
	if !reflect.DeepEqual(got, []string{"java.lang.ClassNotFoundException: org.Foo"}) {
		t.Errorf("exception reasons = %q", got)
	}
}

func TestParseManifest(t *testing.T) {
	m := parseManifest(strings.NewReader("Manifest-Version: 1.0\r\nShort-Name: gradle\r\nPlugin-Dependencies: workflow-api:1.0,structs:1.2\r\n 0,workflow-cps:2.0\r\nMinimum-Java-Version: 11\r\n"))
	if m["Plugin-Dependencies"] != "workflow-api:1.0,structs:1.20,workflow-cps:2.0" || m["Minimum-Java-Version"] != "11" {
		t.Errorf("manifest = %q", m)
	}
}

func TestCopyHome(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	write := func(rel, body string) {
		p := filepath.Join(src, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("config.xml", "<hudson><numExecutors>4</numExecutors></hudson>")
	write("credentials.xml", "<c/>")
	write("secret.key", "k")
	write("secrets/master.key", "m")
	write("queue.xml", "<q/>")
	write("jobs/a/config.xml", "<project/>")
	write("jobs/a/builds/1/build.xml", "<build/>")
	write("jobs/f/config.xml", "<folder/>")
	write("jobs/f/jobs/b/config.xml", "<project/>")
	write("jobs/f/jobs/b/workspace/x/config.xml", "<no/>")
	write("nodes/n1/config.xml", "<slave/>")
	write("workspace/a/file", "x")
	write("plugins/git.jpi", "x")
	hc, err := copyHome(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(hc.Items, []string{"jobs/a", "jobs/f", "jobs/f/jobs/b"}) {
		t.Errorf("items = %q", hc.Items)
	}
	for _, rel := range []string{"credentials.xml", "secret.key", "secrets/master.key", "nodes/n1/config.xml", "jobs/f/jobs/b/config.xml"} {
		if _, err := os.Stat(filepath.Join(dst, rel)); err != nil {
			t.Errorf("%s not copied", rel)
		}
	}
	for _, rel := range []string{"queue.xml", "jobs/a/builds/1/build.xml", "workspace/a/file", "plugins/git.jpi", "jobs/f/jobs/b/workspace/x/config.xml"} {
		if _, err := os.Stat(filepath.Join(dst, rel)); err == nil {
			t.Errorf("%s copied; it is not configuration", rel)
		}
	}
	b, _ := os.ReadFile(filepath.Join(dst, "config.xml"))
	if !strings.Contains(string(b), "<numExecutors>0</numExecutors>") {
		t.Errorf("built-in node still has executors: %s", b)
	}
	if _, err := copyHome(t.TempDir(), dst); err == nil {
		t.Error("a directory without config.xml is not a JENKINS_HOME")
	}
}

func TestPluginOf(t *testing.T) {
	xml := []byte(`<?xml version='1.1' encoding='UTF-8'?>
<flow-definition plugin="workflow-job@1385.v1">
  <builders><hudson.plugins.gradle.Gradle plugin="gradle@2.12"><tasks>build</tasks></hudson.plugins.gradle.Gradle></builders>
  <definition class="org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition" plugin="workflow-cps@3894.v2"/>
  <x.Outer_-Inner plugin="inner@1"/>
</flow-definition>`)
	for class, want := range map[string]string{
		"":                             "workflow-job 1385.v1",
		"hudson.plugins.gradle.Gradle": "gradle 2.12",
		"org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition": "workflow-cps 3894.v2",
		"x.Outer$Inner": "inner 1",
		"not.There":     " ",
	} {
		n, v := pluginOf(xml, class)
		if n+" "+v != want {
			t.Errorf("pluginOf(%q) = %q %q, want %q", class, n, v, want)
		}
	}
}

func TestItemName(t *testing.T) {
	if got := itemName("jobs/f/jobs/b"); got != "f/b" {
		t.Errorf("itemName = %q", got)
	}
}

// An extension Jenkins could not load (old email-ext on a core without
// javax.mail) leaves its plugin working in part.
func TestScoutFailureBreaksItsPlugin(t *testing.T) {
	log := "2026-10-21 12:41:50.000+0000 [id=30]\tWARNING\thudson.ExtensionFinder$Sezpoz#scout: Failed to scout hudson.plugins.emailext.ExtendedEmailPublisher$DescriptorImpl\n" +
		"java.lang.ClassNotFoundException: javax.mail.internet.AddressException\n" +
		"\tat java.base/java.lang.Class.forName0(Native Method)\n" +
		"Caused: java.lang.NoClassDefFoundError: javax/mail/internet/AddressException\n" +
		"2026-10-21 12:41:51.000+0000 [id=30]\tINFO\tjenkins.InitReactorRunner$1#onAttained: Completed initialization\n"
	failed, errs := ParseLog(log)
	if len(failed) != 0 || len(errs) != 1 || !strings.Contains(errs[0], "Failed to scout hudson.plugins.emailext.ExtendedEmailPublisher$DescriptorImpl — java.lang.ClassNotFoundException: javax.mail.internet.AddressException") {
		t.Fatalf("failed = %v, errors = %q", failed, errs)
	}

	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "email-ext.jpi"))
	if err != nil {
		t.Fatal(err)
	}
	z := zip.NewWriter(f)
	if _, err := z.Create("WEB-INF/classes/hudson/plugins/emailext/ExtendedEmailPublisher$DescriptorImpl.class"); err != nil {
		t.Fatal(err)
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	f.Close()
	broken, rest := attribute(errs, dir, nil, map[string]string{"email-ext": "2.40"})
	if len(broken) != 1 || broken[0].Name != "email-ext" || broken[0].Version != "2.40" || len(rest) != 0 {
		t.Errorf("broken = %+v, rest = %q", broken, rest)
	}
}

// The plan's Java comes first unless a plugin needs more; an explicit Java
// always does.
func TestPickImageJava(t *testing.T) {
	docker := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(docker, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		java, prefer, min int
		want              string
	}{
		{0, 8, 0, "jenkins/jenkins:2.346.3-jdk8"},
		{0, 8, 11, "jenkins/jenkins:2.346.3-jdk11"},
		{0, 17, 11, "jenkins/jenkins:2.346.3-jdk17"},
		{8, 0, 11, "jenkins/jenkins:2.346.3-jdk8"},
		{0, 0, 0, "jenkins/jenkins:2.346.3-jdk8"},
		{0, 0, 21, "jenkins/jenkins:2.346.3-jdk8"},
	} {
		o := Options{Core: "2.346.3", Java: c.java, PreferJava: c.prefer, Docker: docker}
		got, _, err := o.pickImage(context.Background(), c.min)
		if err != nil || got != c.want {
			t.Errorf("java %d, prefer %d, plugins need %d: %s, %v; want %s", c.java, c.prefer, c.min, got, err, c.want)
		}
	}
}

func TestLastWords(t *testing.T) {
	log := "Running from: /usr/share/jenkins/jenkins.war\njava.lang.IllegalStateException: boom\n\tat a.b(C.java:1)\n\t... 3 more\n\n"
	if got := lastWords(log); got != "java.lang.IllegalStateException: boom" {
		t.Errorf("lastWords = %q", got)
	}
	log = "2026-10-21 13:19:20.000+0000 [id=26]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed to load the plugins\n" +
		"java.lang.LinkageError: boom\n\tat a.b(C.java:1)\n" +
		"2026-10-21 13:19:25.276+0000 [id=30]\tINFO\tjenkins.model.Jenkins$13#onAttained: Started termination\n" +
		"2026-10-21 13:19:25.279+0000 [id=30]\tSEVERE\tjenkins.model.Jenkins$13#onTaskFailed: Failed IOHubProvider.cleanUp\n" +
		"2026-10-21 13:19:25.491+0000 [id=31]\tINFO\twinstone.Logger#logInternal: Control thread shutdown successfully\n"
	if got, want := lastWords(log), "jenkins.InitReactorRunner$1#onTaskFailed: Failed to load the plugins — java.lang.LinkageError: boom"; got != want {
		t.Errorf("lastWords = %q, want %q", got, want)
	}
}

// A rehearsal of an upgrade tells Jenkins the version it came from, unless
// the copied home already says so; Jenkins then logs the plugins split out
// of core that it installs.
func TestUpgradeFrom(t *testing.T) {
	home := t.TempDir()
	if ok, err := markUpgrade(home, "", "2.580.1"); ok || err != nil {
		t.Errorf("no version: %v, %v", ok, err)
	}
	if ok, err := markUpgrade(home, "2.580.1", "2.580.1"); ok || err != nil {
		t.Errorf("same version: %v, %v", ok, err)
	}
	if ok, err := markUpgrade(home, "2.249.3", "2.580.1"); !ok || err != nil {
		t.Fatalf("upgrade: %v, %v", ok, err)
	}
	if b, _ := os.ReadFile(filepath.Join(home, lastExecVersion)); string(b) != "2.249.3" {
		t.Errorf("recorded %q", b)
	}
	if ok, err := markUpgrade(home, "2.346.3", "2.580.1"); ok || err != nil {
		t.Errorf("home with a version: %v, %v", ok, err)
	}

	log := "2026-10-21 13:36:01.000+0000 [id=30]\tINFO\thudson.PluginManager#loadDetachedPlugins: Upgraded Jenkins from version 2.249.3 to version 2.580.1. Loaded detached plugins (and dependencies): [sshd.hpi, javax-mail-api.hpi, instance-identity]\n" +
		"2026-10-21 13:36:02.000+0000 [id=30]\tINFO\thudson.PluginManager#loadDetachedPlugins: Upgraded detached plugins (and dependencies): [sshd]\n"
	if got := detachedPlugins(log); !slices.Equal(got, []string{"sshd", "javax-mail-api", "instance-identity"}) {
		t.Errorf("detached = %q", got)
	}
	if got := detachedPlugins("Jenkins is fully up and running"); got != nil {
		t.Errorf("no upgrade: %q", got)
	}
}

func TestResultOK(t *testing.T) {
	for _, c := range []struct {
		r    Result
		want bool
	}{
		{Result{Started: true}, true},
		{Result{Started: false}, false},
		{Result{Started: true, Failed: []Failure{{Name: "a"}}}, false},
		{Result{Started: true, Broken: []Failure{{Name: "a"}}}, false},
		{Result{Started: true, Home: &HomeResult{Items: 2, Loaded: 2}}, true},
		{Result{Started: true, Home: &HomeResult{Missing: []ConfigProblem{{Item: "a"}}}}, false},
		{Result{Started: true, Home: &HomeResult{Unreadable: []ConfigProblem{{Item: "a"}}}}, false},
	} {
		if got := c.r.OK(); got != c.want {
			t.Errorf("OK(%+v) = %v, want %v", c.r, got, c.want)
		}
	}
}

func TestProbeScript(t *testing.T) {
	plain := ProbeScript(false)
	if !bytes.HasPrefix(plain, []byte("// Written by jenkins-upgrade-check")) || !bytes.Contains(plain, []byte(ProbeFile)) {
		t.Errorf("probe script starts %q", plain[:40])
	}
	quiet := string(ProbeScript(true))
	if !strings.HasPrefix(quiet, "// Start quiet:") || !strings.Contains(quiet, "jenkins.model.Jenkins.get().doQuietDown()\n\n// Written by") {
		t.Errorf("quiet probe script starts %q", quiet[:200])
	}
	if !bytes.Equal(ProbeScript(false), plain) {
		t.Error("the quiet script must not change the plain one")
	}
}

func TestReadProbe(t *testing.T) {
	dir := t.TempDir()
	if _, err := ReadProbe(filepath.Join(dir, ProbeFile)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("missing: %v", err)
	}
	bad := filepath.Join(dir, "bad.json")
	os.WriteFile(bad, []byte(`{"completed":tru`), 0o644)
	if _, err := ReadProbe(bad); err == nil || !strings.HasPrefix(err.Error(), bad+": ") {
		t.Errorf("bad JSON: %v", err)
	}
	good := filepath.Join(dir, "good.json")
	os.WriteFile(good, []byte(`{"completed":true,"version":"2.580.1",
		"plugins":[{"name":"git","version":"5.8","active":true,"enabled":true},{"name":"off","version":"1","active":false,"enabled":false},{"name":"broken","version":"2","active":false,"enabled":true}],
		"failed":[{"name":"dep","causes":["java.io.IOException: Failed to load: Dep (dep 1)\n - Plugin is missing: x (1.0)"]}]}`), 0o644)
	p, err := ReadProbe(good)
	if err != nil {
		t.Fatal(err)
	}
	if !p.Completed || p.Version != "2.580.1" {
		t.Errorf("probe = %+v", p)
	}
	if got := p.Active(); !reflect.DeepEqual(got, map[string]bool{"git": true}) {
		t.Errorf("active = %v", got)
	}
	// Failed plugins with their reasons, and enabled plugins that are not
	// running; a disabled plugin is not a failure.
	want := []Failure{
		{Name: "broken", Version: "2", Reasons: []string{reasonInactive}},
		{Name: "dep", Reasons: []string{"Plugin is missing: x (1.0)"}},
	}
	if got := p.Failures(); !reflect.DeepEqual(got, want) {
		t.Errorf("failures = %+v", got)
	}
}

// What the log says comes first; the probe fills in a reason the log did
// not give.
func TestMergeProbeKeepsLogReasons(t *testing.T) {
	failed := []Failure{{Name: "a", Version: "1"}, {Name: "b", Version: "2", Reasons: []string{"from the log"}}}
	p := &Probe{}
	json.Unmarshal([]byte(`{"failed":[{"name":"a","causes":["java.lang.LinkageError: x"]},{"name":"b","causes":["java.lang.LinkageError: y"]}]}`), p)
	got := mergeProbe(failed, p, nil)
	want := []Failure{{Name: "a", Version: "1", Reasons: []string{"java.lang.LinkageError: x"}}, {Name: "b", Version: "2", Reasons: []string{"from the log"}}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("merged = %+v", got)
	}
}

func TestShortClassAndNames(t *testing.T) {
	if got := ShortClass("com.thoughtworks.xstream.mapper.CannotResolveClassException: x"); got != "CannotResolveClassException: x" {
		t.Errorf("ShortClass = %q", got)
	}
	if got := ShortClass("java.lang.Error: x"); got != "java.lang.Error: x" {
		t.Errorf("ShortClass = %q", got)
	}
	if got := itemName("orphan"); got != "orphan" {
		t.Errorf("a path with no item part is kept: %q", got)
	}
	if got := ItemPath("team/apps/deploy"); got != "jobs/team/jobs/apps/jobs/deploy" {
		t.Errorf("ItemPath = %q", got)
	}
	if !IsNoise("hudson.diagnosis.OldDataMonitor") || !IsNoise("OldData-123") || IsNoise("OldDataX") || IsNoise("jenkins.security.QueueItemAuthenticatorMonitor") {
		t.Error("IsNoise")
	}
}

// Errors that name a class go to the plugin that ships it, unless that
// plugin failed anyway; the rest stay.
func TestAttribute(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.jpi"), pluginFile(t, "", []string{"org/a/A.class"}, nil), 0o644)
	os.WriteFile(filepath.Join(dir, "b.jpi"), pluginFile(t, "", []string{"org/b/B.class"}, nil), 0o644)
	errs := []string{
		"x: Failed to load org.a.A — boom",
		"x: Failed to instantiate class org.a.A — again",
		"x: Failed to load org.b.B — b failed",
		"x: Failed to load org.nobody.C — no plugin",
		"x: something else",
	}
	for i := 0; i < 6; i++ {
		errs = append(errs, fmt.Sprintf("x: Failed to inject org.a.A — %d", i))
	}
	broken, rest := Attribute(errs, dir, map[string]bool{"b": true})
	if len(broken) != 1 || broken[0].Name != "a" || broken[0].Version != "" || len(broken[0].Reasons) != 5 {
		t.Errorf("broken = %+v (at most five reasons each)", broken)
	}
	if !reflect.DeepEqual(rest, []string{"x: Failed to load org.nobody.C — no plugin", "x: something else"}) {
		t.Errorf("rest = %q", rest)
	}
}

func TestWriteTar(t *testing.T) {
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{"ref/a.txt": "hello", "ref/sub/b.txt": "b"})
	if err := os.Symlink("a.txt", filepath.Join(dir, "ref", "link")); err != nil {
		t.Skip("no symbolic links here:", err)
	}
	var buf bytes.Buffer
	if err := writeTar(&buf, dir, "ref", 1000); err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(&buf)
	got := map[string]string{}
	for {
		h, err := tr.Next()
		if err != nil {
			break
		}
		c, _ := io.ReadAll(tr)
		got[h.Name] = fmt.Sprintf("%c %d %s %s", h.Typeflag, h.Uid, h.Linkname, c)
	}
	want := map[string]string{
		"ref/":          "5 1000  ",
		"ref/a.txt":     "0 1000  hello",
		"ref/link":      "2 1000 a.txt ",
		"ref/sub/":      "5 1000  ",
		"ref/sub/b.txt": "0 1000  b",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tar =\n%q\nwant\n%q", got, want)
	}
	if err := writeTar(io.Discard, dir, "missing", 1000); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a missing directory: %v", err)
	}
}

func TestCopyIn(t *testing.T) {
	d := newFakeDocker(t)
	dir := t.TempDir()
	writeTree(t, dir, map[string]string{"ref/a": "a"})
	o := &Options{Docker: d.path()}
	if err := o.copyIn(context.Background(), dir, "ref", "c:/usr/share/jenkins/"); err != nil {
		t.Fatal(err)
	}
	if got := d.refFile("ref/a"); got != "a" {
		t.Errorf("ref/a = %q", got)
	}
	// docker took the stream, but it was cut short.
	if err := o.copyIn(context.Background(), dir, "missing", "c:/x"); err == nil || !strings.HasPrefix(err.Error(), "docker cp: lstat ") {
		t.Errorf("a missing directory: %v", err)
	}
	o.Docker = filepath.Join(t.TempDir(), "docker")
	if err := o.copyIn(context.Background(), dir, "ref", "c:/x"); err == nil || !strings.HasPrefix(err.Error(), "docker cp: ") {
		t.Errorf("no docker: %v", err)
	}
	if (&Options{}).docker() != "docker" {
		t.Error("the default is docker on the PATH")
	}
}

func TestPickImage(t *testing.T) {
	d := newFakeDocker(t)
	ctx := context.Background()
	var progress bytes.Buffer

	// --image is pulled as given, with the Java given.
	d.set("images", "registry.example.com/jenkins:2.580.1\n")
	o := Options{Core: "2.580.1", Image: "registry.example.com/jenkins:2.580.1", Java: 21, Docker: d.path(), Progress: &progress}
	if img, java, err := o.pickImage(ctx, 0); err != nil || img != o.Image || java != 21 {
		t.Errorf("--image: %s %d %v", img, java, err)
	}
	if !strings.Contains(progress.String(), "pulling registry.example.com/jenkins:2.580.1") {
		t.Errorf("progress = %q", progress.String())
	}
	o.Image = "registry.example.com/jenkins:nope"
	if _, _, err := o.pickImage(ctx, 0); err == nil || err.Error() != "docker pull registry.example.com/jenkins:nope: Error response from daemon: manifest for registry.example.com/jenkins:nope not found: manifest unknown" {
		t.Errorf("--image missing: %v", err)
	}

	// Before 2.303 the Java 8 image has the plain tag.
	d.set("images", "jenkins/jenkins:2.249.3\n")
	o = Options{Core: "2.249.3", Docker: d.path()}
	if img, java, err := o.pickImage(ctx, 0); err != nil || img != "jenkins/jenkins:2.249.3" || java != 8 {
		t.Errorf("2.249.3: %s %d %v", img, java, err)
	}

	// A core with no -jdkN image: the plain tag, Java unknown.
	d.set("images", "jenkins/jenkins:2.580.1\n")
	o = Options{Core: "2.580.1", Docker: d.path()}
	if img, java, err := o.pickImage(ctx, 0); err != nil || img != "jenkins/jenkins:2.580.1" || java != 0 {
		t.Errorf("plain tag: %s %d %v", img, java, err)
	}
}

func TestCopyFileAndTree(t *testing.T) {
	src := t.TempDir()
	writeTree(t, src, map[string]string{"casc/a.yaml": "a", "casc/sub/b.yaml": "b"})
	os.Symlink("a.yaml", filepath.Join(src, "casc", "link.yaml"))

	dst := t.TempDir()
	if err := copyTree(filepath.Join(src, "casc"), dst); err != nil {
		t.Fatal(err)
	}
	for rel, want := range map[string]string{"a.yaml": "a", "sub/b.yaml": "b"} {
		if b, err := os.ReadFile(filepath.Join(dst, rel)); err != nil || string(b) != want {
			t.Errorf("%s = %q, %v", rel, b, err)
		}
	}
	if _, err := os.Lstat(filepath.Join(dst, "link.yaml")); err == nil {
		t.Error("only regular files are copied")
	}
	// A single file goes into the directory under its own name.
	one := t.TempDir()
	if err := copyTree(filepath.Join(src, "casc", "a.yaml"), one); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(one, "a.yaml")); string(b) != "a" {
		t.Errorf("a.yaml = %q", b)
	}
	if err := copyTree(filepath.Join(src, "missing"), one); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("missing: %v", err)
	}

	file := filepath.Join(src, "casc", "a.yaml")
	for name, c := range map[string]struct{ src, dst string }{
		"missing source":  {filepath.Join(src, "nope"), filepath.Join(dst, "x")},
		"parent is file":  {file, filepath.Join(file, "x")},
		"target is dir":   {file, filepath.Join(dst, "sub")},
		"source is a dir": {filepath.Join(src, "casc"), filepath.Join(dst, "y")},
	} {
		if err := copyFile(c.src, c.dst); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}

func TestLastWordsEdges(t *testing.T) {
	if got := lastWords("\tat a.b(C.java:1)\n\t... 3 more\n\n"); got != "" {
		t.Errorf("only a stack trace: %q", got)
	}
	long := strings.Repeat("é", 400)
	if got := lastWords(long); got != strings.Repeat("é", 300)+"…" {
		t.Errorf("a long line is clipped to 300 characters: %d", len([]rune(got)))
	}
}

// A warning that says a class is missing is an error worth showing; a
// plugin reported twice is listed once, with the reasons it was given.
func TestParseLogWarningsAndRepeats(t *testing.T) {
	log := "2026-10-21 12:41:50.000+0000 [id=30]\tWARNING\thudson.PluginManager#x: java.lang.NoClassDefFoundError: org/gone/Thing\n" +
		"java.lang.ClassNotFoundException: org.gone.Thing\n" +
		"\tat a.b(C.java:1)\n" +
		"2026-10-21 12:41:51.000+0000 [id=30]\tWARNING\thudson.Whatever#y: just a warning\n" +
		"not an error\n" +
		"2026-10-21 12:41:52.000+0000 [id=41]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed Loading plugin A v1 (a)\n" +
		"2026-10-21 12:41:53.000+0000 [id=41]\tSEVERE\tjenkins.InitReactorRunner$1#onTaskFailed: Failed Loading plugin A v1 (a)\n" +
		"java.io.IOException: Failed to load: A (a 1)\n" +
		" - Plugin is missing: b (2.0)\n"
	failed, errs := ParseLog(log)
	if want := []Failure{{Name: "a", Version: "1", Reasons: []string{"Plugin is missing: b (2.0)"}}}; !reflect.DeepEqual(failed, want) {
		t.Errorf("failed = %+v", failed)
	}
	if want := []string{"hudson.PluginManager#x: java.lang.NoClassDefFoundError: org/gone/Thing — java.lang.ClassNotFoundException: org.gone.Thing"}; !reflect.DeepEqual(errs, want) {
		t.Errorf("errors = %q", errs)
	}
}

func TestExplainKinds(t *testing.T) {
	for _, c := range []struct {
		reasons []string
		kind    string
		why     string
		fix     string
	}{
		{[]string{"Jenkins (2.479.1) or higher required"}, CauseCoreTooOld, "this release needs Jenkins 2.479.1 or newer", "use an older release of the plugin, or upgrade Jenkins first"},
		{[]string{"Plugin is disabled: Credentials (credentials)"}, CauseDisabledDep, "needs Credentials, which is disabled", "enable Credentials"},
		{[]string{"(trilead-api 1.0) failed to load. Fix this plugin first"}, CauseOther, "(trilead-api 1.0) failed to load. Fix this plugin first", "read the log for this plugin"},
		{[]string{"SSH Build Agents (ssh-slaves 1.31) failed to load. Fix this plugin first"}, CauseDependencyFailed, "needs ssh-slaves, which failed to load", "fix ssh-slaves first"},
		{[]string{"java.lang.NoClassDefFoundError: org/acegisecurity/Authentication"}, CauseRemovedAPI, "uses org.acegisecurity.Authentication, which this Jenkins and its plugins no longer provide", "update the plugin; if no release fixes it, the plugin is abandoned and needs a replacement"},
		{[]string{reasonInactive}, CauseInactive, "installed but not active after startup", "read the log for this plugin"},
		{[]string{"java.lang.IllegalStateException: boom", "second"}, CauseOther, "java.lang.IllegalStateException: boom", "read the log for this plugin"},
		{nil, CauseOther, "Jenkins did not say why", "read the log for this plugin"},
	} {
		d := Explain(Failure{Name: "p", Reasons: c.reasons})
		if d.Kind != c.kind || d.Why != c.why || d.Fix != c.fix {
			t.Errorf("Explain(%q) = %+v, want %s / %s / %s", c.reasons, d, c.kind, c.why, c.fix)
		}
	}
}

// A plugin whose dependency failed, where that dependency did not fail in
// this rehearsal (it is not in the list), is its own root; a cycle of
// dependency failures ends.
func TestRootsDependencyNotInList(t *testing.T) {
	roots := Roots([]Failure{
		{Name: "a", Reasons: []string{"Failed to load: Gone (gone 1.0)"}},
		{Name: "x", Reasons: []string{"Failed to load: Y (y 1)"}},
		{Name: "y", Reasons: []string{"Failed to load: X (x 1)"}},
	})
	var names []string
	for _, r := range roots {
		names = append(names, r.Name)
		if len(r.Cascade) != 0 {
			t.Errorf("%s has a cascade: %q", r.Name, r.Cascade)
		}
	}
	if !reflect.DeepEqual(names, []string{"a", "x", "y"}) {
		t.Errorf("roots = %+v", roots)
	}
}

func TestCopyHomeEdges(t *testing.T) {
	if _, err := copyHome(filepath.Join(t.TempDir(), "missing"), t.TempDir()); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("missing: %v", err)
	}
	file := filepath.Join(t.TempDir(), "config.xml")
	os.WriteFile(file, []byte("<hudson/>"), 0o644)
	if _, err := copyHome(file, t.TempDir()); err == nil || err.Error() != file+" is not a directory" {
		t.Errorf("a file: %v", err)
	}

	src := t.TempDir()
	writeTree(t, src, map[string]string{
		"config.xml":                      "<hudson/>",
		"casc.yml":                        "x",
		"init.groovy":                     "x",
		"jenkins.log":                     "not configuration",
		"secrets/master.key":              "m",
		"jobs/a/config.xml":               "<project/>",
		"jobs/a/builds/1/build.xml":       "<build/>",
		"users/admin_123/config.xml":      "<user/>",
		"init.groovy.d/set-url.groovy":    "x",
		"casc_configs/jenkins.yaml":       "x",
		"nodes/n1/config.xml":             "<slave/>",
		"nodes/n1/disk-usage.xml":         "<x/>",
		"jobs/a/lastStable-target/marker": "x",
	})
	// Jenkins keeps links to builds next to the job; they are not copied.
	os.Symlink("builds/1", filepath.Join(src, "jobs", "a", "lastStable"))

	// A target where secrets/ cannot be made: those files are reported,
	// the rest is copied.
	dst := t.TempDir()
	os.WriteFile(filepath.Join(dst, "secrets"), nil, 0o644)
	hc, err := copyHome(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if len(hc.Skipped) != 1 || !strings.HasPrefix(hc.Skipped[0], "secrets/master.key: mkdir ") {
		t.Errorf("skipped = %q", hc.Skipped)
	}
	for _, rel := range []string{"casc.yml", "init.groovy", "users/admin_123/config.xml", "init.groovy.d/set-url.groovy", "casc_configs/jenkins.yaml", "nodes/n1/config.xml"} {
		if _, err := os.Stat(filepath.Join(dst, rel)); err != nil {
			t.Errorf("%s not copied", rel)
		}
	}
	for _, rel := range []string{"jenkins.log", "nodes/n1/disk-usage.xml", "jobs/a/lastStable", "jobs/a/builds/1/build.xml"} {
		if _, err := os.Lstat(filepath.Join(dst, rel)); err == nil {
			t.Errorf("%s copied", rel)
		}
	}
	if hc.Files != 8 || !reflect.DeepEqual(hc.Items, []string{"jobs/a"}) {
		t.Errorf("files = %d, items = %q", hc.Files, hc.Items)
	}

	// Without a way to take the built-in node's executors away, there is
	// no rehearsal: it could run builds.
	dst = t.TempDir()
	os.Mkdir(filepath.Join(dst, "config.xml"), 0o755)
	if _, err := copyHome(src, dst); err == nil || !strings.Contains(err.Error(), "is a directory") {
		t.Errorf("err = %v", err)
	}
}

func TestInBuildRecords(t *testing.T) {
	for rel, want := range map[string]bool{
		"jobs/a/builds/1/build.xml":          true,
		"jobs/a/workspace@2/config.xml":      true,
		"jobs/a/config.xml":                  false,
		"jobs/builder/jobs/b/config.xml":     false,
		"jobs/f/jobs/b/workspace/config.xml": true,
	} {
		if got := inBuildRecords(rel); got != want {
			t.Errorf("inBuildRecords(%s) = %v", rel, got)
		}
	}
}

func TestBlame(t *testing.T) {
	ref := t.TempDir()
	writeTree(t, ref, map[string]string{
		"jobs/a/config.xml": `<flow-definition plugin="workflow-job@1385.v1"><definition class="org.x.Cps" plugin="workflow-cps"/></flow-definition>`,
		"jobs/b/config.xml": `<project><builders/></project>`,
	})
	failed := map[string]bool{"workflow-cps": true}
	active := map[string]bool{"workflow-job": true}
	for _, c := range []struct {
		p    ConfigProblem
		want ConfigProblem
	}{
		// The class names an element with a plugin of no version.
		{ConfigProblem{Path: "jobs/a", Problem: "CannotResolveClassException: org.x.Cps"},
			ConfigProblem{Path: "jobs/a", Problem: "CannotResolveClassException: org.x.Cps", Plugin: "workflow-cps", PluginState: PluginFailed}},
		// The class is the root element, written as a path.
		{ConfigProblem{Path: "jobs/a", Problem: "NoClassDefFoundError: flow/definition"},
			ConfigProblem{Path: "jobs/a", Problem: "NoClassDefFoundError: flow/definition", Plugin: "workflow-job 1385.v1", PluginState: PluginLoaded}},
		// The root element names no plugin.
		{ConfigProblem{Path: "jobs/b", Problem: "ClassNotFoundException: a.B"}, ConfigProblem{Path: "jobs/b", Problem: "ClassNotFoundException: a.B"}},
		// No item, or no config.xml to read.
		{ConfigProblem{Item: "(global configuration)", Problem: "x"}, ConfigProblem{Item: "(global configuration)", Problem: "x"}},
		{ConfigProblem{Path: "jobs/gone", Problem: "x"}, ConfigProblem{Path: "jobs/gone", Problem: "x"}},
	} {
		p := c.p
		Blame(&p, ref, failed, active)
		if !reflect.DeepEqual(p, c.want) {
			t.Errorf("Blame(%+v) = %+v, want %+v", c.p, p, c.want)
		}
	}
}

func TestMinimumJava(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "old.jpi"), pluginFile(t, "Manifest-Version: 1.0\nMinimum-Java-Version: 1.8\n", nil, nil), 0o644)
	os.WriteFile(filepath.Join(dir, "eleven.jpi"), pluginFile(t, "Minimum-Java-Version: 11.0.2\r\nLong-Name: a long\r\n value\r\n\r\n", nil, nil), 0o644)
	os.WriteFile(filepath.Join(dir, "none.jpi"), pluginFile(t, "Short-Name: none\n", nil, nil), 0o644)
	os.WriteFile(filepath.Join(dir, "nomanifest.jpi"), pluginFile(t, "", []string{"a/B.class"}, nil), 0o644)
	os.WriteFile(filepath.Join(dir, "junk.jpi"), []byte("not a zip"), 0o644)
	os.WriteFile(filepath.Join(dir, "odd.jpi"), pluginFile(t, "Minimum-Java-Version: eleven\n", nil, nil), 0o644)
	if n, who := minimumJava(dir); n != 11 || who != "eleven" {
		t.Errorf("minimumJava = %d, %q", n, who)
	}
	if n, who := minimumJava(filepath.Join(dir, "missing")); n != 0 || who != "" {
		t.Errorf("no plugins: %d, %q", n, who)
	}
	m, err := manifest(filepath.Join(dir, "eleven.jpi"))
	if err != nil || m["Long-Name"] != "a longvalue" {
		t.Errorf("manifest = %q, %v", m, err)
	}
	if m, err := manifest(filepath.Join(dir, "nomanifest.jpi")); err != nil || len(m) != 0 {
		t.Errorf("no manifest = %q, %v", m, err)
	}
	if _, err := manifest(filepath.Join(dir, "junk.jpi")); err == nil {
		t.Error("not a zip")
	}
	if m := parseManifest(strings.NewReader("A: 1\nno colon here\n continued\n")); m["A"] != "1" || len(m) != 1 {
		t.Errorf("a line without a colon ends the entry: %q", m)
	}
}

func TestClassIndex(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.jpi"), pluginFile(t, "", []string{"org/a/A.class", "org/a/A$Inner.class"},
		map[string][]string{"a-lib.jar": {"org/lib/Shared.class", "META-INF/MANIFEST.MF"}}), 0o644)
	os.WriteFile(filepath.Join(dir, "b.hpi"), pluginFile(t, "", nil,
		map[string][]string{"b.jar": {"org/b/B.class", "org/lib/Shared.class"}}), 0o644)
	// A jar that is not one is passed over.
	var buf bytes.Buffer
	z := zip.NewWriter(&buf)
	w, _ := z.Create("WEB-INF/lib/broken.jar")
	w.Write([]byte("not a jar"))
	z.Close()
	os.WriteFile(filepath.Join(dir, "c.jpi"), buf.Bytes(), 0o644)
	os.WriteFile(filepath.Join(dir, "junk.jpi"), []byte("not a zip"), 0o644)

	idx := buildClassIndex(dir)
	for class, want := range map[string]string{
		"org.a.A":           "a",
		"org.a.A$Inner":     "a",
		"org.b.B":           "b",
		"org.lib.Shared":    "a", // the first plugin that ships it
		"org.b.B.class":     "b",
		"org.nobody.Nobody": "",
	} {
		if got := idx.owner(class); got != want {
			t.Errorf("owner(%s) = %q, want %q", class, got, want)
		}
	}
	if len(idx) != 4 {
		t.Errorf("index = %v", idx)
	}
}

// TestDownloadWithoutCache: with no cache directory, plugins go straight into
// dir, checked against their checksum, and nothing is written to a fixed path
// in the shared temp directory.
func TestDownloadWithoutCache(t *testing.T) {
	body := []byte("plugin bytes")
	sum := sha256.Sum256(body)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/gone.hpi" {
			http.NotFound(w, r)
			return
		}
		w.Write(body)
	}))
	defer srv.Close()
	t.Setenv("TMPDIR", t.TempDir())
	dir := t.TempDir()
	picks := []plan.Pick{
		{Name: "good", Version: "1.0", URL: srv.URL + "/good.hpi", SHA256: base64.StdEncoding.EncodeToString(sum[:])},
		{Name: "bad", Version: "1.0", URL: srv.URL + "/bad.hpi", SHA256: base64.StdEncoding.EncodeToString(make([]byte, 32))},
		{Name: "gone", Version: "1.0", URL: srv.URL + "/gone.hpi"},
	}
	got, skipped := Download(context.Background(), picks, dir, "")
	if len(got) != 1 || got[0].Name != "good" || len(skipped) != 2 ||
		!strings.Contains(strings.Join(skipped, "\n"), "checksum mismatch") ||
		!strings.Contains(strings.Join(skipped, "\n"), "404 Not Found") {
		t.Fatalf("got %v, skipped %v; want good downloaded, bad refused and gone reported", got, skipped)
	}
	if b, err := os.ReadFile(filepath.Join(dir, "good.jpi")); err != nil || string(b) != string(body) {
		t.Errorf("good.jpi = %q, %v", b, err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Errorf("dir holds %d entries, want only good.jpi", len(entries))
	}
	if entries, _ := os.ReadDir(os.TempDir()); len(entries) != 0 {
		t.Errorf("wrote to the shared temp directory: %v", entries)
	}
}
