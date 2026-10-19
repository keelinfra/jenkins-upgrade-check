package version

import "strings"

// javaSupport is the Jenkins Java support policy, by weekly line: an LTS
// release follows its weekly base (2.541.3 runs what 2.541 runs). Source:
// https://www.jenkins.io/doc/book/platform-information/support-policy-java/
// The policy applies to the controller and to every agent.
var javaSupport = []struct {
	Java        int
	FirstWeekly string
	LastWeekly  string // empty while still supported
}{
	{8, "2.54", "2.356"},
	{11, "2.164", "2.462"},
	{17, "2.340", "2.544"},
	{21, "2.463", ""},
	{25, "2.534", ""},
}

// Base returns the weekly line a core version belongs to: "2.541" for
// "2.541.3", and the version itself for a weekly.
func Base(core string) string {
	parts := strings.Split(core, ".")
	if len(parts) > 2 {
		return strings.Join(parts[:2], ".")
	}
	return core
}

// SupportedJava lists the Java versions a core runs on, oldest first.
func SupportedJava(core string) []int {
	base := Base(core)
	var out []int
	for _, j := range javaSupport {
		if !AtLeast(base, j.FirstWeekly) {
			continue
		}
		if j.LastWeekly != "" && Less(j.LastWeekly, base) {
			continue
		}
		out = append(out, j.Java)
	}
	return out
}

// SupportsJava reports whether core runs on the given Java feature release.
func SupportsJava(core string, java int) bool {
	for _, j := range SupportedJava(core) {
		if j == java {
			return true
		}
	}
	return false
}

// CommonJava lists the Java versions both cores run on, oldest first.
func CommonJava(a, b string) []int {
	var out []int
	for _, j := range SupportedJava(a) {
		if SupportsJava(b, j) {
			out = append(out, j)
		}
	}
	return out
}

// Bridges are the last LTS releases of each Java generation. Each one runs
// on the Java it is the last for and on the next, so a controller can change
// Java there and keep going.
var Bridges = []string{"2.346.3", "2.462.3", "2.541.3"}

// Hop is one step of an upgrade: switch the controller (and its agents) to
// Java, then upgrade the core to Core.
type Hop struct {
	Core string `json:"core"`
	Java int    `json:"java"`
}

// Path returns the hops from one core to another. A single hop means the
// target can be reached directly; more hops mean no Java version runs both
// ends, so the controller has to stop at a bridge release first.
func Path(from, to string) []Hop {
	var hops []Hop
	cur := from
	for guard := 0; guard < 8; guard++ {
		if common := CommonJava(cur, to); len(common) > 0 {
			return append(hops, Hop{Core: to, Java: common[0]})
		}
		next := ""
		for i := len(Bridges) - 1; i >= 0; i-- {
			b := Bridges[i]
			if Less(cur, b) && Less(b, to) && len(CommonJava(cur, b)) > 0 {
				next = b
				break
			}
		}
		if next == "" {
			return append(hops, Hop{Core: to})
		}
		hops = append(hops, Hop{Core: next, Java: CommonJava(cur, next)[0]})
		cur = next
	}
	return hops
}
