// Package words writes counts as English: "1 plugin", "2 plugins".
//
// Reports used to say "1 plugin(s)" and "1 setting(s) come", which reads as a
// form left half filled in. Every count a report prints goes through here, so
// the grammar is decided once.
package words

import "fmt"

// Count renders n with noun, adding "s" for anything but one.
func Count(n int, noun string) string {
	return CountOf(n, noun, noun+"s")
}

// CountOf renders n with one or many, for nouns whose plural is not noun+"s".
func CountOf(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

// Verb picks the verb form that agrees with n: Verb(1, "is", "are") is "is".
func Verb(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}
