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
	_ fmt.Stringer           = (*AnimeRelations)(nil)
	_ api.Editable           = (*AnimeRelations)(nil)
	_ api.ArrayEventListener = (*AnimeRelations)(nil)
)

// Authorize returns an error if the given API POST request is not authorized.
func (relations *AnimeRelations) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil || (user.Role != "editor" && user.Role != "admin") {
		return errors.New("Not logged in or not authorized to edit")
	}

	return nil
}

// Edit saves a log entry for the edit.
func (relations *AnimeRelations) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "AnimeRelations", relations.AnimeID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	return false, nil
}

// OnAppend saves a log entry.
func (relations *AnimeRelations) OnAppend(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayAppend", "AnimeRelations", relations.AnimeID, fmt.Sprintf("%s[%d]", key, index), "", fmt.Sprint(obj))
	logEntry.Save()
}

// OnRemove saves a log entry.
func (relations *AnimeRelations) OnRemove(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayRemove", "AnimeRelations", relations.AnimeID, fmt.Sprintf("%s[%d]", key, index), fmt.Sprint(obj), "")
	logEntry.Save()
}

// Save saves the anime relations object in the database.
func (relations *AnimeRelations) Save() {
	DB.Set("AnimeRelations", relations.AnimeID, relations)
}

// Delete deletes the relation list from the database.
func (relations *AnimeRelations) Delete() error {
	DB.Delete("AnimeRelations", relations.AnimeID)
	return nil
}
