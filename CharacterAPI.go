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
	_ fmt.Stringer = (*Character)(nil)
	_ api.Editable = (*Character)(nil)
)

// Actions
func init() {
	API.RegisterActions("Character", []*api.Action{
		// Like character
		LikeAction(),

		// Unlike character
		UnlikeAction(),
	})
}

// Authorize returns an error if the given API request is not authorized.
func (character *Character) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	// Allow custom actions (like, unlike) for normal users
	if action == "action" {
		return nil
	}

	if user.Role != "editor" && user.Role != "admin" {
		return errors.New("Insufficient permissions")
	}

	return nil
}

// Edit creates an edit log entry.
func (character *Character) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (consumed bool, err error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "Character", character.ID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	return false, nil
}

// Save saves the character in the database.
func (character *Character) Save() {
	DB.Set("Character", character.ID, character)
}
