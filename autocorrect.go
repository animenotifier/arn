package arn

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// FixGenre ...
func FixGenre(genre string) string {
	genre = strings.Replace(genre, "-", "", -1)
	genre = strings.Replace(genre, " ", "", -1)
	genre = strings.ToLower(genre)
	return genre
}

// Capitalize returns the string with the first letter capitalized.
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}
