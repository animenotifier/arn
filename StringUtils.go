package arn

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var stripTagsRegex = regexp.MustCompile(`<[^>]*>`)
var sourceRegex = regexp.MustCompile(`\(Source: (.*?)\)`)
var writtenByRegex = regexp.MustCompile(`\[Written by (.*?)\]`)

// GetGenreIDByName ...
func GetGenreIDByName(genre string) string {
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
