package arn

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aerogo/aero"
)

// AnimeList ...
type AnimeList struct {
	UserID string           `json:"userId"`
	Items  []*AnimeListItem `json:"items"`

	user *User
}

// Find returns the list item with the specified anime ID, if available.
func (list *AnimeList) Find(animeID string) *AnimeListItem {
	for _, item := range list.Items {
		if item.AnimeID == animeID {
			return item
		}
	}

	return nil
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

			list.Items[index] = item

			return nil
		}
	}

	return errors.New("Not found")
}

// Add adds an anime to the list if it hasn't been added yet.
func (list *AnimeList) Add(id interface{}) error {
	animeID := id.(string)

	if list.Contains(animeID) {
		return errors.New("Anime " + animeID + " has already been added")
	}

	creationDate := time.Now().UTC().Format(time.RFC3339)

	newItem := &AnimeListItem{
		AnimeID: animeID,
		Status:  AnimeListStatusPlanned,
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

// Authorize returns an error if the given API request is not authorized.
func (list *AnimeList) Authorize(ctx *aero.Context) error {
	if !ctx.HasSession() {
		return errors.New("Neither logged in nor in session")
	}

	userID, ok := ctx.Session().Get("userId").(string)

	if !ok || userID == "" {
		return errors.New("Not logged in")
	}

	if userID != ctx.Get("id") {
		return errors.New("Can not modify data from other users")
	}

	return nil
}

// User returns the user this anime list belongs to.
func (list *AnimeList) User() *User {
	if list.user == nil {
		list.user, _ = GetUser(list.UserID)
	}

	return list.user
}

// TransformBody returns an item that is passed to methods like Add, Remove, etc.
func (list *AnimeList) TransformBody(body []byte) interface{} {
	if len(body) > 0 && body[0] == '{' {
		item := &AnimeListItem{}
		err := json.Unmarshal(body, item)

		if err != nil {
			panic(err)
		}

		return item
	}

	return string(body)
}

// Save saves the anime list in the database.
func (list *AnimeList) Save() error {
	return DB.Set("AnimeList", list.UserID, list)
}

// GetAnimeList ...
func GetAnimeList(userID string) (*AnimeList, error) {
	obj, err := DB.Get("AnimeList", userID)
	return obj.(*AnimeList), err
}
