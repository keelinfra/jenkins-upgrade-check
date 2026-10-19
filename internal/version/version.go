// Package version orders Jenkins core and plugin versions the way Jenkins
// does. Jenkins' hudson.util.VersionNumber is a copy of Maven's
// ComparableVersion, so "1.10" > "1.9", "2.0-beta-1" < "2.0", and an
// incremental such as "1336.vf33a_a_9863911" orders by its leading number.
package version

import (
	"math/big"
	"strings"
	"unicode"
)

// Version is a parsed version string. The zero value is the empty version.
type Version struct {
	raw   string
	items *item
}

// Parse never fails: anything Jenkins accepts as a version, and anything it
// does not, has an order.
func Parse(s string) Version {
	return Version{raw: s, items: parse(s)}
}

// String returns the text the version was parsed from.
func (v Version) String() string { return v.raw }

// Compare returns -1, 0 or 1.
func (v Version) Compare(o Version) int {
	a, b := v.items, o.items
	if a == nil {
		a = &item{kind: listKind}
	}
	if b == nil {
		b = &item{kind: listKind}
	}
	return compare(a, b)
}

// Compare parses both strings and compares them.
func Compare(a, b string) int { return Parse(a).Compare(Parse(b)) }

// Less reports whether a orders before b.
func Less(a, b string) bool { return Compare(a, b) < 0 }

// AtLeast reports whether v is the same as or newer than min.
func AtLeast(v, min string) bool { return Compare(v, min) >= 0 }

// IsPrerelease reports whether a version is one the update sites hold back
// from regular users: alpha, beta, release candidate, milestone or snapshot.
func IsPrerelease(s string) bool {
	l := strings.ToLower(s)
	for _, q := range []string{"alpha", "beta", "snapshot", "-rc", ".rc", "-m1", "-m2", "-m3", "milestone"} {
		if strings.Contains(l, q) {
			return true
		}
	}
	return false
}

type kind int

const (
	intKind kind = iota
	stringKind
	listKind
)

type item struct {
	kind kind
	num  *big.Int
	str  string
	list []*item
}

// Qualifier order, lowest first. The empty string is a release.
var qualifiers = []string{"alpha", "beta", "milestone", "rc", "snapshot", "", "sp"}

var aliases = map[string]string{"ga": "", "final": "", "release": "", "cr": "rc"}

const releaseIndex = "5"

func comparableQualifier(q string) string {
	for i, k := range qualifiers {
		if k == q {
			return string(rune('0' + i))
		}
	}
	return "7-" + q
}

func newString(s string, followedByDigit bool) *item {
	if followedByDigit && len(s) == 1 {
		switch s {
		case "a":
			s = "alpha"
		case "b":
			s = "beta"
		case "m":
			s = "milestone"
		}
	}
	if a, ok := aliases[s]; ok {
		s = a
	}
	return &item{kind: stringKind, str: s}
}

func newInt(s string) *item {
	n := new(big.Int)
	if _, ok := n.SetString(s, 10); !ok {
		n.SetInt64(0)
	}
	return &item{kind: intKind, num: n}
}

func parseItem(isDigit bool, s string) *item {
	if isDigit {
		return newInt(s)
	}
	return newString(s, false)
}

func (it *item) isNull() bool {
	switch it.kind {
	case intKind:
		return it.num.Sign() == 0
	case stringKind:
		return comparableQualifier(it.str) == releaseIndex
	default:
		return len(it.list) == 0
	}
}

func (it *item) normalize() {
	for i := len(it.list) - 1; i >= 0; i-- {
		last := it.list[i]
		if last.isNull() {
			it.list = append(it.list[:i], it.list[i+1:]...)
		} else if last.kind != listKind {
			break
		}
	}
}

func parse(version string) *item {
	version = strings.ToLower(strings.TrimSpace(version))
	root := &item{kind: listKind}
	list := root
	stack := []*item{list}
	push := func() {
		next := &item{kind: listKind}
		list.list = append(list.list, next)
		list = next
		stack = append(stack, list)
	}
	isDigit := false
	start := 0
	runes := []rune(version)
	for i, c := range runes {
		switch {
		case c == '.':
			if i == start {
				list.list = append(list.list, newInt("0"))
			} else {
				list.list = append(list.list, parseItem(isDigit, string(runes[start:i])))
			}
			start = i + 1
		case c == '-':
			if i == start {
				list.list = append(list.list, newInt("0"))
			} else {
				list.list = append(list.list, parseItem(isDigit, string(runes[start:i])))
			}
			start = i + 1
			push()
		case unicode.IsDigit(c):
			if !isDigit && i > start {
				list.list = append(list.list, newString(string(runes[start:i]), true))
				start = i
				push()
			}
			isDigit = true
		default:
			if isDigit && i > start {
				list.list = append(list.list, parseItem(true, string(runes[start:i])))
				start = i
				push()
			}
			isDigit = false
		}
	}
	if len(runes) > start {
		list.list = append(list.list, parseItem(isDigit, string(runes[start:])))
	}
	for i := len(stack) - 1; i >= 0; i-- {
		stack[i].normalize()
	}
	return root
}

// compare orders a against b; b may be nil, which stands for padding.
func compare(a, b *item) int {
	switch a.kind {
	case intKind:
		if b == nil {
			return a.num.Sign()
		}
		switch b.kind {
		case intKind:
			return a.num.Cmp(b.num)
		default:
			return 1
		}
	case stringKind:
		if b == nil {
			return strings.Compare(comparableQualifier(a.str), releaseIndex)
		}
		switch b.kind {
		case stringKind:
			return strings.Compare(comparableQualifier(a.str), comparableQualifier(b.str))
		default:
			return -1
		}
	default:
		if b == nil {
			if len(a.list) == 0 {
				return 0
			}
			return compare(a.list[0], nil)
		}
		switch b.kind {
		case intKind:
			return -1
		case stringKind:
			return 1
		}
		n := len(a.list)
		if len(b.list) > n {
			n = len(b.list)
		}
		for i := 0; i < n; i++ {
			var l, r *item
			if i < len(a.list) {
				l = a.list[i]
			}
			if i < len(b.list) {
				r = b.list[i]
			}
			var res int
			if l == nil {
				if r != nil {
					res = -compare(r, nil)
				}
			} else {
				res = compare(l, r)
			}
			if res != 0 {
				return res
			}
		}
		return 0
	}
}
