package dryrun

import (
	"os"
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
