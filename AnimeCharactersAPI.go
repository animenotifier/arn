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
	_ fmt.Stringer           = (*AnimeCharacters)(nil)
	_ api.Editable           = (*AnimeCharacters)(nil)
	_ api.ArrayEventListener = (*AnimeCharacters)(nil)
)

// Authorize returns an error if the given API POST request is not authorized.
func (characters *AnimeCharacters) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil || (user.Role != "editor" && user.Role != "admin") {
		return errors.New("Not logged in or not authorized to edit")
	}

	return nil
}

// Edit saves a log entry for the edit.
func (characters *AnimeCharacters) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "AnimeCharacters", characters.AnimeID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	return false, nil
}

// OnAppend saves a log entry.
func (characters *AnimeCharacters) OnAppend(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayAppend", "AnimeCharacters", characters.AnimeID, fmt.Sprintf("%s[%d]", key, index), "", fmt.Sprint(obj))
	logEntry.Save()
}

// OnRemove saves a log entry.
func (characters *AnimeCharacters) OnRemove(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayRemove", "AnimeCharacters", characters.AnimeID, fmt.Sprintf("%s[%d]", key, index), fmt.Sprint(obj), "")
	logEntry.Save()
}

// Save saves the character in the database.
func (characters *AnimeCharacters) Save() {
	DB.Set("AnimeCharacters", characters.AnimeID, characters)
}

// Delete deletes the character list from the database.
func (characters *AnimeCharacters) Delete() error {
	DB.Delete("AnimeCharacters", characters.AnimeID)
	return nil
}
