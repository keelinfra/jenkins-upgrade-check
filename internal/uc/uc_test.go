package uc

import (
	"errors"
	"strings"
	"testing"
)

func TestReadUpdateCenterJSONP(t *testing.T) {
	u, err := ReadUpdateCenter(strings.NewReader(`updateCenter.post(
{"plugins":{"git":{"name":"git","version":"5.2.1"}},"warnings":[],"deprecations":{}}
);`))
	if err != nil {
		t.Fatal(err)
	}
	if !u.Distributed("git") || u.Distributed("nope") {
		t.Errorf("got %+v", u.Plugins)
	}
}

func TestWarningAffects(t *testing.T) {
	u, err := ReadUpdateCenter(strings.NewReader(`{"warnings":[
		{"id":"A","name":"p","type":"plugin","versions":[{"lastVersion":"1.75","pattern":"1[.]([0-6][0-9]|7[0-5])(|[.-].*)"}]},
		{"id":"B","name":"p","type":"plugin","versions":[{"firstVersion":"2.0","lastVersion":"2.3"}]},
		{"id":"C","name":"p","type":"plugin","versions":[{"lastVersion":"9","pattern":"(?<!x)bad"}]}
	]}`))
	if err != nil {
		t.Fatal(err)
	}
	ids := func(v string) string {
		var s []string
		for _, w := range u.WarningsFor("plugin", "p", v) {
			s = append(s, w.ID)
		}
		return strings.Join(s, ",")
	}
	cases := map[string]string{
		"1.75":    "A,C",
		"1.76":    "C",
		"2.1":     "B,C",
		"2.4":     "C",
		"1.9":     "C",
		"10.0":    "",
		"1336.v1": "",
	}
	for v, want := range cases {
		if got := ids(v); got != want {
			t.Errorf("warnings for %s = %q, want %q", v, got, want)
		}
	}
}

func TestLatestFor(t *testing.T) {
	db := newReleases(map[string]map[string]*Release{
		"p": {
			"1.0":      {RequiredCore: "2.300"},
			"2.0":      {RequiredCore: "2.400"},
			"3.0-beta": {RequiredCore: "2.400"},
			"3.0":      {RequiredCore: "2.500"},
		},
	})
	for core, want := range map[string]string{"2.346.3": "1.0", "2.400": "2.0", "2.462.3": "2.0", "2.555.1": "3.0"} {
		if r := db.LatestFor("p", core); r == nil || r.Version != want {
			t.Errorf("LatestFor(p, %s) = %v, want %s", core, r, want)
		}
	}
	if db.LatestFor("p", "2.200") != nil {
		t.Error("no release runs on 2.200")
	}
	if db.Latest("p").Version != "3.0" {
		t.Error("Latest skips the beta")
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("disk on fire") }

func TestReadErrors(t *testing.T) {
	if _, err := ReadUpdateCenter(failingReader{}); err == nil || err.Error() != "disk on fire" {
		t.Errorf("ReadUpdateCenter: %v", err)
	}
	if _, err := ReadUpdateCenter(strings.NewReader("<html>")); err == nil || !strings.HasPrefix(err.Error(), "update center JSON: ") {
		t.Errorf("ReadUpdateCenter: %v", err)
	}
	if _, err := ReadReleases(strings.NewReader(`{"plugins":[]}`)); err == nil || !strings.HasPrefix(err.Error(), "plugin versions JSON: ") {
		t.Errorf("ReadReleases: %v", err)
	}
}

func TestDeprecatedAndWarningKinds(t *testing.T) {
	u, err := ReadUpdateCenter(strings.NewReader(`{
		"deprecations":{"old":{"url":"https://plugins.jenkins.io/old/"}},
		"warnings":[
			{"id":"CORE-1","name":"core","type":"core","versions":[{"lastVersion":"2.440"}]},
			{"id":"P-1","name":"p","type":"plugin","versions":[{"pattern":""},{"firstVersion":"1.0"}]},
			{"id":"P-2","name":"other","type":"plugin","versions":[{"lastVersion":"9"}]}
		]}`))
	if err != nil {
		t.Fatal(err)
	}
	if url, ok := u.Deprecated("old"); !ok || url != "https://plugins.jenkins.io/old/" {
		t.Errorf("Deprecated(old) = %q, %v", url, ok)
	}
	if _, ok := u.Deprecated("git"); ok {
		t.Error("git is not deprecated")
	}
	if ws := u.WarningsFor("core", "core", "2.426.3"); len(ws) != 1 || ws[0].ID != "CORE-1" {
		t.Errorf("core warnings = %+v", ws)
	}
	// A range with neither a pattern nor a last version affects nothing.
	if ws := u.WarningsFor("plugin", "p", "1.5"); len(ws) != 0 {
		t.Errorf("p warnings = %+v", ws)
	}
}

func TestReleases(t *testing.T) {
	rel, err := ReadReleases(strings.NewReader(`{"plugins":{
		"p":{
			"1.0":{"requiredCore":"2.300","url":"https://u/p-1.0.hpi"},
			"2.0":{"requiredCore":"2.400"},
			"2.1-beta":{"requiredCore":"2.400"},
			"3.0":{"requiredCore":"2.500"},
			"4.0":null
		},
		"beta-only":{"1.0-beta-1":{}}
	}}`))
	if err != nil {
		t.Fatal(err)
	}
	if !rel.Known("p") || rel.Known("q") {
		t.Error("Known")
	}
	if r := rel.Get("p", "1.0"); r == nil || r.Name != "p" || r.Version != "1.0" || r.URL != "https://u/p-1.0.hpi" {
		t.Errorf("Get(p, 1.0) = %+v", r)
	}
	if rel.Get("p", "4.0") != nil || rel.Get("p", "9.9") != nil {
		t.Error("a null entry or a version never published is not a release")
	}
	var got []string
	for _, r := range rel.LoadableFor("p", "2.462.3") {
		got = append(got, r.Version)
	}
	if strings.Join(got, ",") != "2.0,1.0" {
		t.Errorf("LoadableFor(p, 2.462.3) = %v, want newest first without the beta", got)
	}
	if rel.LoadableFor("p", "2.200") != nil {
		t.Error("nothing loads on 2.200")
	}
	if rel.Latest("beta-only") != nil || rel.Latest("q") != nil {
		t.Error("Latest skips pre-releases and unknown plugins")
	}
}

func TestArchiveURL(t *testing.T) {
	if got, want := ArchiveURL("momentjs", "1.1.1"), "https://updates.jenkins.io/download/plugins/momentjs/1.1.1/momentjs.hpi"; got != want {
		t.Errorf("ArchiveURL = %q, want %q", got, want)
	}
}
