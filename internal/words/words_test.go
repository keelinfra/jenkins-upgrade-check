package words

import "testing"

func TestCount(t *testing.T) {
	for _, c := range []struct {
		got, want string
	}{
		{Count(0, "plugin"), "0 plugins"},
		{Count(1, "plugin"), "1 plugin"},
		{Count(2, "plugin"), "2 plugins"},
		{CountOf(1, "dependency", "dependencies"), "1 dependency"},
		{CountOf(3, "dependency", "dependencies"), "3 dependencies"},
		{Verb(1, "is", "are"), "is"},
		{Verb(2, "is", "are"), "are"},
	} {
		if c.got != c.want {
			t.Errorf("got %q, want %q", c.got, c.want)
		}
	}
}
