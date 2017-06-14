package arn

import (
	"errors"
)

// AnimeListStatus values for anime list items
const (
	AnimeListStatusWatching  = "watching"
	AnimeListStatusCompleted = "completed"
	AnimeListStatusPlanned   = "planned"
	AnimeListStatusDropped   = "dropped"
	AnimeListStatusHold      = "hold"
)

// AnimeList ...
type AnimeList struct {
	UserID string           `json:"userId"`
	Items  []*AnimeListItem `json:"items"`
}

// AnimeListItem ...
type AnimeListItem struct {
	AnimeID      string      `json:"animeId"`
	Status       string      `json:"status"`
	Episode      int         `json:"episode"`
	Rating       AnimeRating `json:"rating"`
	Notes        string      `json:"notes"`
	RewatchCount int         `json:"rewatchCount"`
	Private      bool        `json:"private"`

	anime *Anime
}

// Add adds an anime to the list if it hasn't been added yet.
func (list *AnimeList) Add(id interface{}) error {
	animeID := id.(string)

	if list.Contains(animeID) {
		return errors.New("Anime " + animeID + " has already been added")
	}

	newItem := &AnimeListItem{
		AnimeID: animeID,
		Status:  AnimeListStatusPlanned,
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

// TransformBody returns an item that is passed to methods like Add, Remove, etc.
func (list *AnimeList) TransformBody(body []byte) interface{} {
	return string(body)
}

// Anime fetches the associated anime data.
func (item *AnimeListItem) Anime() *Anime {
	if item.anime == nil {
		item.anime, _ = GetAnime(item.AnimeID)
	}

	return item.anime
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
