package arn

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/fatih/color"
)

// Authorize returns an error if the given API POST request is not authorized.
func (user *User) Authorize(ctx *aero.Context) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// PostBody reads the POST body and returns an object
// that is passed to methods like Update, Add, Remove, etc.
func (user *User) PostBody(body []byte) interface{} {
	if len(body) > 0 && body[0] == '{' {
		var updates interface{}
		PanicOnError(json.Unmarshal(body, &updates))
		return updates.(map[string]interface{})
	}

	return string(body)
}

// Update updates the user object with the data we received from the PostBody method.
func (user *User) Update(data interface{}) error {
	updates := data.(map[string]interface{})

	return SetObjectProperties(user, updates, func(fullKeyName string, field reflect.StructField, property reflect.Value, newValue reflect.Value) bool {
		// Automatically correct account nicks
		if strings.HasPrefix(fullKeyName, "Accounts.") && strings.HasSuffix(fullKeyName, ".Nick") {
			newNick := newValue.String()
			newNick = FixAccountNick(newNick)
			property.SetString(newNick)
			return true
		}

		switch fullKeyName {
		case "Nick":
			newNick := newValue.String()
			err := user.SetNick(newNick)

			if err != nil {
				color.Red(err.Error())
			}

			return true

		default:
			return false
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
