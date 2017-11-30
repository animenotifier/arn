package arn

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

var whitespace = rune(' ')

// RemoveSpecialCharacters ...
func RemoveSpecialCharacters(s string) string {
	return strings.Map(
		func(r rune) rune {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return whitespace
			}

			return r
		},
		s,
	)
}

// Capitalize returns the string with the first letter capitalized.
func Capitalize(s string) string {
	if s == "" {
		return ""
	}

	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// ToString converts anything into a string.
func ToString(v interface{}) string {
	return fmt.Sprint(v)
}

// Plural returns the number concatenated to the proper pluralization of the word.
func Plural(count int, singular string) string {
	if count == 1 || count == -1 {
		return ToString(count) + " " + singular
	}

	return ToString(count) + " " + singular + "s"
}

// ContainsUnicodeLetters tells you if unicode characters are inside the string.
func ContainsUnicodeLetters(s string) bool {
	return len(s) != len([]rune(s))
}
