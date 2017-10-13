package arn

import (
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
	"github.com/animenotifier/arn/autocorrect"
	"github.com/fatih/color"
)

// Force interface implementations
var (
	_ api.Editable = (*User)(nil)
)

// Authorize returns an error if the given API POST request is not authorized.
func (user *User) Authorize(ctx *aero.Context, action string) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// Edit updates the user object.
func (user *User) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	// Automatically correct account nicks
	if strings.HasPrefix(key, "Accounts.") && strings.HasSuffix(key, ".Nick") {
		newNick := newValue.String()
		newNick = autocorrect.FixAccountNick(newNick)
		value.SetString(newNick)

		// Refresh osu info if the name changed
		if key == "Accounts.Osu.Nick" {
			go func() {
				err := user.RefreshOsuInfo()

				if err != nil {
					color.Red("Error refreshing osu info of user '%s' with osu nick '%s': %v", user.Nick, newNick, err)
				} else {
					color.Green("Refreshed osu info of user '%s' with osu nick '%s': %v", user.Nick, newNick, user.Accounts.Osu.PP)
				}

				user.Save()
			}()
		}

		return true, nil
	}

	switch key {
	case "Nick":
		newNick := newValue.String()
		err := user.SetNick(newNick)
		return true, err
	}

	return false, nil
}

// Save saves the user object in the database.
func (user *User) Save() error {
	return DB.Set("User", user.ID, user)
}

// Filter removes privacy critical fields from the user object.
func (user *User) Filter() {
	user.Email = ""
	user.Gender = ""
	user.FirstName = ""
	user.LastName = ""
	user.IP = ""
	user.LastLogin = ""
	user.LastSeen = ""
	user.Accounts.Facebook.ID = ""
	user.Accounts.Google.ID = ""
	user.Accounts.Twitter.ID = ""
	user.AgeRange = UserAgeRange{}
	user.Location = UserLocation{}
}

// ShouldFilter tells whether data needs to be filtered in the given context.
func (user *User) ShouldFilter(ctx *aero.Context) bool {
	ctxUser := GetUserFromContext(ctx)

	if ctxUser != nil && ctxUser.Role == "admin" {
		return false
	}

	return true
}
