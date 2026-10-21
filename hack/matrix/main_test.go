package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLastPatches(t *testing.T) {
	meta := `<versions><version>2.58</version><version>2.60.1</version><version>2.60.3</version>
<version>2.60.2</version><version>2.541.1</version><version>2.541.3</version><version>2.580.1</version>
<version>2.46.3</version><version>2.585</version></versions>`
	got := LastPatches(meta, "2.580.1", 60)
	want := []string{"2.60.3", "2.541.3"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LastPatches = %v, want %v", got, want)
	}
}

func TestTable(t *testing.T) {
	m, err := Load("testdata/results")
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile("testdata/jenkins.csv")
	if err != nil {
		t.Fatal(err)
	}
	if err := m.AddShares(string(b)); err != nil {
		t.Fatal(err)
	}
	if m.To != "2.580.1" || len(m.Rows) != 2 || m.Rows[0].From != "2.346.3" {
		t.Fatalf("rows not loaded in version order: %+v", m.Rows)
	}
	if m.Rows[0].Share != 0.3 || m.Rows[1].Share != 0.6 {
		t.Errorf("shares = %v, %v; want 0.3 and 0.6 (weeklies are not part of an LTS line)", m.Rows[0].Share, m.Rows[1].Share)
	}
	md := m.Markdown()
	for _, want := range []string{
		"Jenkins does not start from 1 of them, and starts with plugins that fail to load from 1 more; 90.0% of all reporting controllers",
		"every stop loads cleanly from 2 of 2",
		"1 starting points already have plugins that do not load",
		"| [2.346.3](2.346.3-top300/summary.md) | 30.0% | ✅ all 349 load | ❌ Jenkins does not start: instance-identity | 2.541.3 (Java 17) → 2.580.1 (Java 21) | ✅ all 349 load |",
		"| ⚠️ 2 of 330 fail: gitlab-plugin |",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("table lacks %q:\n%s", want, md)
		}
	}
}
