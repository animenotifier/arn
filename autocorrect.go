package arn

import "strings"

// FixGenre ...
func FixGenre(genre string) string {
	return strings.ToLower(genre)
}
