package arn

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	return strings.TrimSpace(description)
}

// FixGender ...
func FixGender(gender string) string {
	if gender != "male" && gender != "female" {
		return ""
	}

	return gender
}

// Capitalize returns the string with the first letter capitalized.
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// AnimeRatingStars displays the rating in Unicode stars.
func AnimeRatingStars(rating float64) string {
	stars := int(rating/20 + 0.5)
	return strings.Repeat("★", stars) + strings.Repeat("☆", 5-stars)
}

// EpisodesToString shows a question mark if the episode count is zero.
func EpisodesToString(episodes int) string {
	if episodes == 0 {
		return "?"
	}

	return ToString(episodes)
}

// EpisodeCountMax is used for the max value of number input on episodes.
func EpisodeCountMax(episodes int) string {
	if episodes == 0 {
		return ""
	}

	return strconv.Itoa(episodes)
}

// DateTimeUTC returns the current UTC time in RFC3339 format.
func DateTimeUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
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

// PanicOnError will panic if the error is not nil.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// PrettyPrint prints the object as indented JSON data on the console.
func PrettyPrint(obj interface{}) {
	pretty, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(pretty))
}
