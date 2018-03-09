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
	_ fmt.Stringer           = (*Anime)(nil)
	_ api.Editable           = (*Anime)(nil)
	_ api.CustomEditable     = (*Anime)(nil)
	_ api.ArrayEventListener = (*Anime)(nil)
)

// Edit creates an edit log entry.
func (anime *Anime) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (consumed bool, err error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "Anime", anime.ID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	return false, nil
}

// OnAppend saves a log entry.
func (anime *Anime) OnAppend(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayAppend", "Anime", anime.ID, fmt.Sprintf("%s[%d]", key, index), "", fmt.Sprint(obj))
	logEntry.Save()
}

// OnRemove saves a log entry.
func (anime *Anime) OnRemove(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayRemove", "Anime", anime.ID, fmt.Sprintf("%s[%d]", key, index), fmt.Sprint(obj), "")
	logEntry.Save()
}

// Authorize returns an error if the given API POST request is not authorized.
func (anime *Anime) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil || (user.Role != "editor" && user.Role != "admin") {
		return errors.New("Not logged in or not authorized to edit this anime")
	}

	return nil
}

// AfterEdit updates the metadata.
func (anime *Anime) AfterEdit(ctx *aero.Context) error {
	anime.Edited = DateTimeUTC()
	anime.EditedBy = GetUserFromContext(ctx).ID
	return nil
}

// Save saves the anime in the database.
func (anime *Anime) Save() {
	DB.Set("Anime", anime.ID, anime)
}
