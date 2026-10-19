package dryrun

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
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
