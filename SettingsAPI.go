package arn

import (
	"errors"
	"reflect"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ api.Editable = (*Settings)(nil)
	_ api.Filter   = (*Settings)(nil)
)

// Authorize returns an error if the given API POST request is not authorized.
func (settings *Settings) Authorize(ctx *aero.Context, action string) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// Edit updates the settings object.
func (settings *Settings) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	switch key {
	// case "Avatar.Source":
	// 	settings.Avatar.Source = newValue.String()
	// 	settings.Save() // Save needed here because RefreshAvatar fetches the settings on a DIFFERENT server
	// 	settings.User().RefreshAvatar()
	// 	return true, nil

	// case "Avatar.SourceURL":
	// 	settings.Avatar.SourceURL = newValue.String()
	// 	settings.Save() // Save needed here because RefreshAvatar fetches the settings on a DIFFERENT server
	// 	settings.User().RefreshAvatar()
	// 	return true, nil

	case "Theme":
		if settings.User().IsPro() {
			settings.Theme = newValue.String()
		} else {
			return true, errors.New("PRO accounts only")
		}

		return true, nil

	case "CalendarSettings.ShowUserList":
		settings.CalendarSettings.ShowUserList = newValue.Bool()
		return true, nil
	}

	return false, nil
}

// Filter removes privacy critical fields from the settings object.
func (settings *Settings) Filter() {
	settings.Notification.Email = ""
}

// ShouldFilter tells whether data needs to be filtered in the given context.
func (settings *Settings) ShouldFilter(ctx *aero.Context) bool {
	ctxUser := GetUserFromContext(ctx)

	if ctxUser != nil && ctxUser.Role == "admin" {
		return false
	}

	return true
}

// Save saves the settings in the database.
func (settings *Settings) Save() {
	DB.Set("Settings", settings.UserID, settings)
}
