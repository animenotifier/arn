package arn

import (
	"reflect"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ api.CustomEditable = (*AnimeListItem)(nil)
	_ api.AfterEditable  = (*AnimeListItem)(nil)
)

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
