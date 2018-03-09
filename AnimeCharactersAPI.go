package arn

import (
	"errors"

	"github.com/aerogo/aero"
)

// Authorize returns an error if the given API POST request is not authorized.
func (chars *AnimeCharacters) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil || (user.Role != "editor" && user.Role != "admin") {
		return errors.New("Not logged in or not authorized to edit")
	}

	return nil
}

// Save saves the character in the database.
func (chars *AnimeCharacters) Save() {
	DB.Set("AnimeCharacters", chars.AnimeID, chars)
}
