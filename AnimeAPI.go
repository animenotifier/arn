package arn

import (
	"errors"
	"reflect"

	"github.com/aerogo/aero"
)

// Authorize returns an error if the given API POST request is not authorized.
func (anime *Anime) Authorize(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	if user == nil || (user.Role != "editor" && user.Role != "admin") {
		return errors.New("Not logged in or not authorized to edit this anime")
	}

	return nil
}

// Update updates the anime object with the data we received from the PostBody method.
func (anime *Anime) Update(ctx *aero.Context, data interface{}) error {
	updates := data.(map[string]interface{})

	return SetObjectProperties(anime, updates, func(fullKeyName string, field *reflect.StructField, property *reflect.Value, newValue reflect.Value) bool {
		switch fullKeyName {
		case "Custom:ShoboiID":
			oldValue := anime.GetMapping("shoboi/anime")
			newValue := newValue.Interface().(string)

			anime.RemoveMapping("shoboi/anime", oldValue)

			if newValue != "" {
				user := GetUserFromContext(ctx)
				anime.AddMapping("shoboi/anime", newValue, user.ID)
			}

			return true

		default:
			return false
		}
	})
}

// Save saves the anime in the database.
func (anime *Anime) Save() error {
	return DB.Set("Anime", anime.ID, anime)
}
