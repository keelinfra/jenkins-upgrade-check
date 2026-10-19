package version

import (
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	ordered := [][2]string{
		{"1.9", "1.10"},
		{"2.0-beta-1", "2.0"},
		{"2.0-alpha-1", "2.0-beta-1"},
		{"2.0-rc-1", "2.0"},
		{"2.0", "2.0.1"},
		{"2.0", "2.0-sp1"},
		{"2.462.3", "2.475"},
		{"2.475", "2.479.1"},
		{"2.541.3", "2.555.1"},
		{"1.75", "1336.vf33a_a_9863911"},
		{"1335.vf07d9ce377a_e", "1336.vf33a_a_9863911"},
		{"2.13.2-21.v8e6b_3f1a_2d0b", "2.13.2-22.v1"},
		{"2.15.2", "2.15.2.1"},
		{"4.0.0-123.vabc", "4.1.0-1.vabc"},
	}
	for _, p := range ordered {
		if Compare(p[0], p[1]) >= 0 {
			t.Errorf("%s should order before %s", p[0], p[1])
		}
		if Compare(p[1], p[0]) <= 0 {
			t.Errorf("%s should order after %s", p[1], p[0])
		}
	}
	equal := [][2]string{
		{"1.0", "1"},
		{"1.0.0", "1"},
		{"1-final", "1"},
		{"1.0-ga", "1.0"},
		{"2.0-cr1", "2.0-rc1"},
	}
	for _, p := range equal {
		if c := Compare(p[0], p[1]); c != 0 {
			t.Errorf("%s and %s should be equal, got %d", p[0], p[1], c)
		}
	}
}

func TestIsPrerelease(t *testing.T) {
	for _, v := range []string{"2.0-beta-1", "1.0-alpha", "3.0-rc1", "1.0-SNAPSHOT"} {
		if !IsPrerelease(v) {
			t.Errorf("%s is a pre-release", v)
		}
	}
	for _, v := range []string{"1336.vf33a_a_9863911", "2.13.2-21.v8e6b_3f1a_2d0b", "5.2.1"} {
		if IsPrerelease(v) {
			t.Errorf("%s is a release", v)
		}
	}
}

func TestSupportedJava(t *testing.T) {
	cases := map[string][]int{
		"2.346.3": {8, 11, 17},
		"2.361.4": {11, 17},
		"2.462.3": {11, 17},
		"2.479.1": {17, 21},
		"2.541.3": {17, 21, 25},
		"2.544":   {17, 21, 25},
		"2.545":   {21, 25},
		"2.555.1": {21, 25},
		"2.222.4": {8, 11},
	}
	for core, want := range cases {
		if got := SupportedJava(core); !reflect.DeepEqual(got, want) {
			t.Errorf("SupportedJava(%s) = %v, want %v", core, got, want)
		}
	}
}

func TestPath(t *testing.T) {
	cases := []struct {
		from, to string
		want     []Hop
	}{
		{"2.541.3", "2.555.1", []Hop{{"2.555.1", 21}}},
		{"2.492.3", "2.555.1", []Hop{{"2.555.1", 21}}},
		{"2.462.3", "2.555.1", []Hop{{"2.541.3", 17}, {"2.555.1", 21}}},
		{"2.346.3", "2.555.1", []Hop{{"2.541.3", 17}, {"2.555.1", 21}}},
		{"2.222.4", "2.555.1", []Hop{{"2.462.3", 11}, {"2.541.3", 17}, {"2.555.1", 21}}},
		{"2.426.3", "2.479.3", []Hop{{"2.479.3", 17}}},
	}
	for _, c := range cases {
		if got := Path(c.from, c.to); !reflect.DeepEqual(got, c.want) {
			t.Errorf("Path(%s, %s) = %v, want %v", c.from, c.to, got, c.want)
		}
	}
}
