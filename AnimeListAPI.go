package arn

import (
	"encoding/json"
	"errors"

	"github.com/aerogo/aero"
)

// Add adds an anime to the list if it hasn't been added yet.
func (list *AnimeList) Add(id interface{}) error {
	animeID := id.(string)

	if list.Contains(animeID) {
		return errors.New("Anime " + animeID + " has already been added")
	}

	creationDate := DateTimeUTC()

	newItem := &AnimeListItem{
		AnimeID: animeID,
		Status:  AnimeListStatusPlanned,
		Rating:  &AnimeRating{},
		Created: creationDate,
		Edited:  creationDate,
	}

	list.Items = append(list.Items, newItem)

	return nil
}

// Remove removes the anime ID from the list.
func (list *AnimeList) Remove(id interface{}) bool {
	animeID := id.(string)

	for index, item := range list.Items {
		if item.AnimeID == animeID {
			list.Items = append(list.Items[:index], list.Items[index+1:]...)
			return true
		}
	}

	return false
}

// Contains checks if the list contains the anime ID already.
func (list *AnimeList) Contains(id interface{}) bool {
	animeID := id.(string)

	for _, item := range list.Items {
		if item.AnimeID == animeID {
			return true
		}
	}

	return false
}

// Get ...
func (list *AnimeList) Get(id interface{}) (interface{}, error) {
	item := list.Find(id.(string))

	if item == nil {
		return nil, errors.New("Not found")
	}

	return item, nil
}

// Set ...
func (list *AnimeList) Set(id interface{}, value interface{}) error {
	animeID := id.(string)

	for index, item := range list.Items {
		if item.AnimeID == animeID {
			item, ok := value.(*AnimeListItem)

			if !ok {
				return errors.New("Missing anime list item properties")
			}

			if item.AnimeID != animeID {
				return errors.New("Incorrect animeId property")
			}

			item.Edited = DateTimeUTC()
			list.Items[index] = item

			return nil
		}
	}

	return errors.New("Not found")
}

// Update ...
func (list *AnimeList) Update(id interface{}, updatesObj interface{}) error {
	updates := updatesObj.(map[string]interface{})
	animeID := id.(string)

	for _, item := range list.Items {
		if item.AnimeID == animeID {
			err := SetObjectProperties(item, updates, nil)
			item.Edited = DateTimeUTC()

			item.Rating.Clamp()

			for key := range updates {
				switch key {
				case "Episodes":
					item.OnEpisodesChange()

				case "Status":
					item.OnStatusChange()
				}
			}

			return err
		}
	}

	return errors.New("Not found")
}

// Authorize returns an error if the given API request is not authorized.
func (list *AnimeList) Authorize(ctx *aero.Context) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// PostBody returns an item that is passed to methods like Add, Remove, etc.
func (list *AnimeList) PostBody(body []byte) interface{} {
	if len(body) > 0 && body[0] == '{' {
		var updates interface{}
		PanicOnError(json.Unmarshal(body, &updates))
		return updates.(map[string]interface{})
	}

	return string(body)
}

// Save saves the anime list in the database.
func (list *AnimeList) Save() error {
	return DB.Set("AnimeList", list.UserID, list)
}
