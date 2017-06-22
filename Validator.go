package arn

import (
	"regexp"
	"strings"
)

var fixNickRegex = regexp.MustCompile(`[\W\s\d]`)

// FixUserNick automatically corrects a username.
func FixUserNick(nick string) string {
	nick = fixNickRegex.ReplaceAllString(nick, "")

	if nick == "" {
		return nick
	}

	for len(nick) > 1 && strings.HasPrefix(nick, "_") {
		nick = nick[1:]
	}

	if nick == "_" {
		return ""
	}

	return strings.ToUpper(string(nick[0])) + nick[1:]
}

// IsValidNick tests if the given nickname is valid.
func IsValidNick(nick string) bool {
	if len(nick) < 2 {
		return false
	}

	return nick == FixUserNick(nick)
}

// IsValidEmail tests if the given email address is valid.
func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}

	return true
}
