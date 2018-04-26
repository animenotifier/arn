package stringutils

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/xrash/smetrics"
)

var whitespace = rune(' ')

// RemoveSpecialCharacters ...
func RemoveSpecialCharacters(s string) string {
	return strings.Map(
		func(r rune) rune {
			if r == '-' {
				return -1
			}

			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return whitespace
			}

			return r
		},
		s,
	)
}

// AdvancedStringSimilarity is like StringSimilarity but boosts the value if a appears directly in b.
func AdvancedStringSimilarity(a string, b string) float64 {
	if a == b {
		return 1000000
	}

	s := StringSimilarity(a, b)

	if strings.Contains(b, a) {
		s += 0.5

		if strings.HasPrefix(b, a) {
			s += 5.0
		}
	}

	return s
}

// StringSimilarity returns 1.0 if the strings are equal and goes closer to 0 when they are different.
func StringSimilarity(a string, b string) float64 {
	return smetrics.JaroWinkler(a, b, 0.7, 4)
}

// Capitalize returns the string with the first letter capitalized.
func Capitalize(s string) string {
	if s == "" {
		return ""
	}

	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

// Plural returns the number concatenated to the proper pluralization of the word.
func Plural(count int, singular string) string {
	if count == 1 || count == -1 {
		return fmt.Sprint(count) + " " + singular
	}

	return fmt.Sprint(count) + " " + singular + "s"
}

// ContainsUnicodeLetters tells you if unicode characters are inside the string.
func ContainsUnicodeLetters(s string) bool {
	return len(s) != len([]rune(s))
}

// PrettyPrint prints the object as indented JSON data on the console.
func PrettyPrint(obj interface{}) {
	pretty, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(pretty))
}
