package validator

import (
	"strings"
	"time"

	"github.com/animenotifier/arn/autocorrect"
)

// IsValidNick tests if the given nickname is valid.
func IsValidNick(nick string) bool {
	if len(nick) < 2 {
		return false
	}

	return nick == autocorrect.FixUserNick(nick)
}

// IsValidDate tells you whether the date is valid.
func IsValidDate(date string) bool {
	if date == "" || strings.HasPrefix(date, "0001") {
		return false
	}

	_, err := time.Parse(time.RFC3339, date)
	return err == nil
}

// IsValidEmail tests if the given email address is valid.
func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}

	// TODO: ...

	return true
}
