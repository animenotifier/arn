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
	_ fmt.Stringer           = (*AnimeEpisodes)(nil)
	_ api.Editable           = (*AnimeEpisodes)(nil)
	_ api.ArrayEventListener = (*AnimeEpisodes)(nil)
)

// Authorize returns an error if the given API POST request is not authorized.
func (episodes *AnimeEpisodes) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil || (user.Role != "editor" && user.Role != "admin") {
		return errors.New("Not logged in or not authorized to edit")
	}

	return nil
}

// Edit saves a log entry for the edit.
func (episodes *AnimeEpisodes) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "AnimeEpisodes", episodes.AnimeID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	return false, nil
}

// OnAppend saves a log entry.
func (episodes *AnimeEpisodes) OnAppend(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayAppend", "AnimeEpisodes", episodes.AnimeID, fmt.Sprintf("%s[%d]", key, index), "", fmt.Sprint(obj))
	logEntry.Save()
}

// OnRemove saves a log entry.
func (episodes *AnimeEpisodes) OnRemove(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayRemove", "AnimeEpisodes", episodes.AnimeID, fmt.Sprintf("%s[%d]", key, index), fmt.Sprint(obj), "")
	logEntry.Save()
}

// Save saves the episodes in the database.
func (episodes *AnimeEpisodes) Save() {
	DB.Set("AnimeEpisodes", episodes.AnimeID, episodes)
}

// Delete deletes the episode list from the database.
func (episodes *AnimeEpisodes) Delete() error {
	DB.Delete("AnimeEpisodes", episodes.AnimeID)
	return nil
}
