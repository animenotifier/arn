package arn

import (
	"errors"
	"reflect"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ api.Creatable = (*AnimeListItem)(nil)
)

// Create is the constructor.
func (item *AnimeListItem) Create(ctx *aero.Context) error {
	data, err := ctx.RequestBodyJSONObject()

	if err != nil {
		return err
	}

	item.AnimeID = data["AnimeID"].(string)
	item.Status = AnimeListStatusPlanned
	item.Rating = &AnimeRating{}
	item.Created = DateTimeUTC()
	item.Edited = item.Created

	if item.Anime() == nil {
		return errors.New("Invalid anime ID")
	}

	if GetUserFromContext(ctx).AnimeList().Contains(item.AnimeID) {
		return errors.New("Anime " + item.AnimeID + " has already been added")
	}

	return nil
}

// Edit ...
func (item *AnimeListItem) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	switch key {
	case "Episodes":
		item.Episodes = int(newValue.Float())

		if item.Episodes < 0 {
			item.Episodes = 0
		}

		item.OnEpisodesChange()
		return true, nil

	case "Status":
		item.Status = newValue.String()
		item.OnStatusChange()
		return true, nil
	}

	return false, nil
}

// AfterEdit is called after the item is edited.
func (item *AnimeListItem) AfterEdit(ctx *aero.Context) error {
	item.Rating.Clamp()
	item.Edited = DateTimeUTC()
	return nil
}
