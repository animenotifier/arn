package arn

import (
	"errors"
	"reflect"

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

// VirtualEdit updates virtual properties.
func (anime *Anime) VirtualEdit(ctx *aero.Context, key string, newValue reflect.Value) (bool, error) {
	switch key {
	case "Virtual:ShoboiID":
		oldValue := anime.GetMapping("shoboi/anime")
		newValue := newValue.Interface().(string)

		anime.RemoveMapping("shoboi/anime", oldValue)

		if newValue != "" {
			user := GetUserFromContext(ctx)
			anime.AddMapping("shoboi/anime", newValue, user.ID)
		}

		return true, nil

	case "Virtual:AniListID":
		oldValue := anime.GetMapping("anilist/anime")
		newValue := newValue.Interface().(string)

		anime.RemoveMapping("anilist/anime", oldValue)

		if newValue != "" {
			user := GetUserFromContext(ctx)
			anime.AddMapping("anilist/anime", newValue, user.ID)
		}

		return true, nil
	}

	return false, nil
}

// Save saves the anime in the database.
func (anime *Anime) Save() error {
	return DB.Set("Anime", anime.ID, anime)
}
