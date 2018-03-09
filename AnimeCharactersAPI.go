package arn

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ fmt.Stringer = (*AnimeCharacters)(nil)
	_ api.Editable = (*AnimeCharacters)(nil)
)

// Authorize returns an error if the given API POST request is not authorized.
func (chars *AnimeCharacters) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil || (user.Role != "editor" && user.Role != "admin") {
		return errors.New("Not logged in or not authorized to edit")
	}

	return nil
}

// Edit saves a log entry for the edit.
func (chars *AnimeCharacters) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "AnimeCharacters", chars.AnimeID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	return false, nil
}

// Save saves the character in the database.
func (chars *AnimeCharacters) Save() {
	DB.Set("AnimeCharacters", chars.AnimeID, chars)
}
