package arn

import (
	"strings"
	"unicode"
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
