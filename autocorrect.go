package arn

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var stripTagsRegex = regexp.MustCompile(`<[^>]*>`)
var sourceRegex = regexp.MustCompile(`\(Source: (.*?)\)`)
var writtenByRegex = regexp.MustCompile(`\[Written by (.*?)\]`)

// FixGenre ...
func FixGenre(genre string) string {
	genre = strings.Replace(genre, "-", "", -1)
	genre = strings.Replace(genre, " ", "", -1)
	genre = strings.ToLower(genre)
	return genre
}

// FixAnimeDescription ...
func FixAnimeDescription(description string) string {
	description = stripTagsRegex.ReplaceAllString(description, "")
	description = sourceRegex.ReplaceAllString(description, "")
	description = writtenByRegex.ReplaceAllString(description, "")
	return description
}

// Capitalize returns the string with the first letter capitalized.
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}
