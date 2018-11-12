package arn

import (
	"errors"
	"reflect"
	"time"

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
	user := GetUserFromContext(ctx)

	if user == nil {
		return true, errors.New("Not logged in")
	}

	switch key {
	case "Episodes":
		fromEpisode := item.Episodes
		toEpisode := int(newValue.Float())

		// Create or update activity if new episode amount is higher
		if toEpisode > fromEpisode {
			lastActivity := user.LastActivityConsumeAnime(item.AnimeID)

			if lastActivity == nil || time.Since(lastActivity.GetCreatedTime()) > 1*time.Hour {
				// If there is no last activity for the given anime,
				// or if the last activity happened more than an hour ago,
				// create a new activity.
				activity := NewActivityConsumeAnime(item.AnimeID, fromEpisode, toEpisode, user.ID)
				activity.Save()
			} else if toEpisode > lastActivity.ToEpisode {
				// Otherwise, update the last activity if episode count is higher.
				lastActivity.ToEpisode = toEpisode
				lastActivity.Save()
			}
		}

		item.Episodes = toEpisode

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
