package arn

import (
	"errors"

	"github.com/aerogo/aero"
)

// Authorize returns an error if the given API POST request is not authorized.
func (anime *Anime) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil || (user.Role != "editor" && user.Role != "admin") {
		return errors.New("Not logged in or not authorized to edit this anime")
	}

	return nil
}

// Save saves the anime in the database.
func (anime *Anime) Save() {
	DB.Set("Anime", anime.ID, anime)
}
