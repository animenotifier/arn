package arn

import (
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/animenotifier/arn/autocorrect"
	"github.com/fatih/color"
)

// Authorize returns an error if the given API POST request is not authorized.
func (user *User) Authorize(ctx *aero.Context) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// Update updates the user object with the data we received from the PostBody method.
func (user *User) Update(ctx *aero.Context, data interface{}) error {
	updates := data.(map[string]interface{})

	return SetObjectProperties(user, updates, func(fullKeyName string, field *reflect.StructField, property *reflect.Value, newValue reflect.Value) (bool, error) {
		// Automatically correct account nicks
		if strings.HasPrefix(fullKeyName, "Accounts.") && strings.HasSuffix(fullKeyName, ".Nick") {
			newNick := newValue.String()
			newNick = autocorrect.FixAccountNick(newNick)
			property.SetString(newNick)

			// Refresh osu info if the name changed
			if fullKeyName == "Accounts.Osu.Nick" {
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

		switch fullKeyName {
		case "Nick":
			newNick := newValue.String()
			err := user.SetNick(newNick)
			return true, err

		default:
			return false, nil
		}
	})
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
