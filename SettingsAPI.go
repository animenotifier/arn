package arn

import (
	"github.com/aerogo/aero"
)

// Authorize returns an error if the given API POST request is not authorized.
func (settings *Settings) Authorize(ctx *aero.Context) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// Update updates the settings object.
func (settings *Settings) Update(ctx *aero.Context, data interface{}) error {
	updates := data.(map[string]interface{})
	return SetObjectProperties(settings, updates, nil)
}

// Save saves the settings in the database.
func (settings *Settings) Save() error {
	return DB.Set("Settings", settings.UserID, settings)
}
