package uc

import (
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
