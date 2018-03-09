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
	_ api.Editable       = (*Anime)(nil)
	_ api.CustomEditable = (*Anime)(nil)
)

// Edit creates an edit log entry.
func (anime *Anime) Edit(ctx *aero.Context, key string, oldValue reflect.Value, newValue reflect.Value) (consumed bool, err error) {
	fmt.Println(key, oldValue.String(), newValue.String())

	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "Anime", anime.ID, oldValue.String(), newValue.String())
	logEntry.Save()

	return false, nil
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
