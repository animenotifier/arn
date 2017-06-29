package arn

import "sort"

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

// User returns the user this anime list belongs to.
func (list *AnimeList) User() *User {
	if list.user == nil {
		list.user, _ = GetUser(list.UserID)
	}

	return list.user
}

// Sort ...
func (list *AnimeList) Sort() {
	sort.Slice(list.Items, func(i, j int) bool {
		a := list.Items[i].Anime().UpcomingEpisode()
		b := list.Items[j].Anime().UpcomingEpisode()

		if a == nil && b == nil {
			return list.Items[i].FinalRating() > list.Items[j].FinalRating()
		}

		if a == nil {
			return false
		}

		if b == nil {
			return true
		}

		return a.Episode.AiringDate.Start < b.Episode.AiringDate.Start
	})
}

// StreamAnimeLists returns a stream of all anime.
func StreamAnimeLists() (chan *AnimeList, error) {
	objects, err := DB.All("AnimeList")
	return objects.(chan *AnimeList), err
}

// AllAnimeLists returns a slice of all anime.
func AllAnimeLists() ([]*AnimeList, error) {
	var all []*AnimeList

	stream, err := StreamAnimeLists()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}

// GetAnimeList ...
func GetAnimeList(user *User) (*AnimeList, error) {
	obj, err := DB.Get("AnimeList", user.ID)
	return obj.(*AnimeList), err
}
