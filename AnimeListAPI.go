package arn

import (
	"github.com/aerogo/aero"
)

// // Add adds an anime to the list if it hasn't been added yet.
// func (list *AnimeList) Add(animeID string) error {
// 	animeID := id.(string)

// 	if list.Contains(animeID) {
// 		return errors.New("Anime " + animeID + " has already been added")
// 	}

// 	creationDate := DateTimeUTC()

// 	newItem := &AnimeListItem{
// 		AnimeID: animeID,
// 		Status:  AnimeListStatusPlanned,
// 		Rating:  &AnimeRating{},
// 		Created: creationDate,
// 		Edited:  creationDate,
// 	}

// 	list.Items = append(list.Items, newItem)

// 	return nil
// }

// Remove removes the anime ID from the list.
func (list *AnimeList) Remove(animeID string) bool {
	for index, item := range list.Items {
		if item.AnimeID == animeID {
			list.Items = append(list.Items[:index], list.Items[index+1:]...)
			return true
		}
	}

	return false
}

// Contains checks if the list contains the anime ID already.
func (list *AnimeList) Contains(animeID string) bool {
	for _, item := range list.Items {
		if item.AnimeID == animeID {
			return true
		}
	}

	return false
}

// Authorize returns an error if the given API request is not authorized.
func (list *AnimeList) Authorize(ctx *aero.Context, action string) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// Save saves the anime list in the database.
func (list *AnimeList) Save() error {
	return DB.Set("AnimeList", list.UserID, list)
}
