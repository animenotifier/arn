package arn

import (
	"sort"
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

// Import adds an anime to the list if it hasn't been added yet
// and if it did exist it will update episode, rating and notes.
func (list *AnimeList) Import(item *AnimeListItem) {
	existing := list.Find(item.AnimeID)

	// If it doesn't exist yet: Simply add it.
	if existing == nil {
		list.Items = append(list.Items, item)
		return
	}

	// If it exists: Copy the attributes to the existing item.
	existing.Status = item.Status
	existing.OnStatusChange()

	if item.Episodes > existing.Episodes {
		existing.Episodes = item.Episodes
		existing.OnEpisodesChange()
	}

	if existing.Rating.Overall == 0 {
		existing.Rating.Overall = item.Rating.Overall
	}

	if existing.Notes == "" {
		existing.Notes = item.Notes
	}

	if item.RewatchCount > existing.RewatchCount {
		existing.RewatchCount = item.RewatchCount
	}

	existing.Edited = DateTimeUTC()
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
			return list.Items[i].Rating.Overall > list.Items[j].Rating.Overall
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

// WatchingAndPlanned ...
func (list *AnimeList) WatchingAndPlanned() *AnimeList {
	newList := &AnimeList{
		UserID: list.UserID,
		Items:  []*AnimeListItem{},
	}

	for _, item := range list.Items {
		if item.Status == AnimeListStatusWatching || item.Status == AnimeListStatusPlanned {
			newList.Items = append(newList.Items, item)
		}
	}

	return newList
}

// SplitByStatus splits the anime list into multiple ones by status.
func (list *AnimeList) SplitByStatus() map[string]*AnimeList {
	statusToList := map[string]*AnimeList{}

	statusToList[AnimeListStatusWatching] = &AnimeList{
		UserID: list.UserID,
		Items:  []*AnimeListItem{},
	}

	statusToList[AnimeListStatusCompleted] = &AnimeList{
		UserID: list.UserID,
		Items:  []*AnimeListItem{},
	}

	statusToList[AnimeListStatusPlanned] = &AnimeList{
		UserID: list.UserID,
		Items:  []*AnimeListItem{},
	}

	statusToList[AnimeListStatusHold] = &AnimeList{
		UserID: list.UserID,
		Items:  []*AnimeListItem{},
	}

	statusToList[AnimeListStatusDropped] = &AnimeList{
		UserID: list.UserID,
		Items:  []*AnimeListItem{},
	}

	for _, item := range list.Items {
		statusList := statusToList[item.Status]
		statusList.Items = append(statusList.Items, item)
	}

	return statusToList
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
