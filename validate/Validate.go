package validate

import (
	"strings"
	"time"

	"github.com/animenotifier/arn/autocorrect"
)

// Nick tests if the given nickname is valid.
func Nick(nick string) bool {
	if len(nick) < 2 {
		return false
	}

	return nick == autocorrect.UserNick(nick)
}

// Date tells you whether the date is valid.
func Date(date string) bool {
	if date == "" || strings.HasPrefix(date, "0001") {
		return false
	}

	_, err := time.Parse(time.RFC3339, date)
	return err == nil
}

// Email tests if the given email address is valid.
func Email(email string) bool {
	if email == "" {
		return false
	}

	// TODO: ...

	return true
}
