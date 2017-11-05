package arn

import (
	"errors"
	"sort"

	"github.com/aerogo/nano"
)

// AnimeList ...
type AnimeList struct {
	UserID string           `json:"userId"`
	Items  []*AnimeListItem `json:"items"`
}

// Add adds an anime to the list if it hasn't been added yet.
func (list *AnimeList) Add(animeID string) error {
	if list.Contains(animeID) {
		return errors.New("Anime " + animeID + " has already been added")
	}

	creationDate := DateTimeUTC()

	item := &AnimeListItem{
		AnimeID: animeID,
		Status:  AnimeListStatusPlanned,
		Rating:  &AnimeRating{},
		Created: creationDate,
		Edited:  creationDate,
	}

	if item.Anime() == nil {
		return errors.New("Invalid anime ID")
	}

	list.Items = append(list.Items, item)
	return nil
}

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

	// Temporary save it before changing the status
	// because status changes can modify the episode count.
	// This will prevent loss of "episodes watched" data.
	existingEpisodes := existing.Episodes

	// Status
	existing.Status = item.Status
	existing.OnStatusChange()

	// Episodes
	if item.Episodes > existingEpisodes {
		existing.Episodes = item.Episodes
	} else {
		existing.Episodes = existingEpisodes
	}

	existing.OnEpisodesChange()

	// Rating
	if existing.Rating.Overall == 0 {
		existing.Rating.Overall = item.Rating.Overall
		existing.Rating.Clamp()
	}

	if existing.Notes == "" {
		existing.Notes = item.Notes
	}

	if item.RewatchCount > existing.RewatchCount {
		existing.RewatchCount = item.RewatchCount
	}

	// Edited
	existing.Edited = DateTimeUTC()
}

// User returns the user this anime list belongs to.
func (list *AnimeList) User() *User {
	user, _ := GetUser(list.UserID)
	return user
}

// Sort ...
func (list *AnimeList) Sort() {
	sort.Slice(list.Items, func(i, j int) bool {
		a := list.Items[i]
		b := list.Items[j]

		if (a.Status != AnimeListStatusWatching && a.Status != AnimeListStatusPlanned) && (b.Status != AnimeListStatusWatching && b.Status != AnimeListStatusPlanned) {
			if a.Rating.Overall == b.Rating.Overall {
				return a.Anime().Title.Canonical < b.Anime().Title.Canonical
			}

			return a.Rating.Overall > b.Rating.Overall
		}

		epsA := a.Anime().UpcomingEpisode()
		epsB := b.Anime().UpcomingEpisode()

		if epsA == nil && epsB == nil {
			if a.Rating.Overall == b.Rating.Overall {
				return a.Anime().Title.Canonical < b.Anime().Title.Canonical
			}

			return a.Rating.Overall > b.Rating.Overall
		}

		if epsA == nil {
			return false
		}

		if epsB == nil {
			return true
		}

		return epsA.Episode.AiringDate.Start < epsB.Episode.AiringDate.Start
	})
}

// Watching ...
func (list *AnimeList) Watching() *AnimeList {
	return list.FilterStatus(AnimeListStatusWatching)
}

// FilterStatus ...
func (list *AnimeList) FilterStatus(status string) *AnimeList {
	newList := &AnimeList{
		UserID: list.UserID,
		Items:  []*AnimeListItem{},
	}

	for _, item := range list.Items {
		if item.Status == status { // (item.Status == AnimeListStatusPlanned)
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

// NormalizeRatings normalizes all ratings so that they are perfectly stretched among the full scale.
func (list *AnimeList) NormalizeRatings() {
	mapped := map[float64]float64{}
	all := []float64{}

	for _, item := range list.Items {
		// Zero rating counts as not rated
		if item.Rating.Overall == 0 {
			continue
		}

		_, found := mapped[item.Rating.Overall]

		if !found {
			mapped[item.Rating.Overall] = item.Rating.Overall
			all = append(all, item.Rating.Overall)
		}
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i] < all[j]
	})

	count := len(all)

	// Prevent division by zero
	if count <= 1 {
		return
	}

	step := 9.9 / float64(count-1)
	currentRating := 0.1

	for _, rating := range all {
		mapped[rating] = currentRating
		currentRating += step
	}

	for _, item := range list.Items {
		item.Rating.Overall = mapped[item.Rating.Overall]
		item.Rating.Clamp()
	}
}

// StreamAnimeLists returns a stream of all anime.
func StreamAnimeLists() chan *AnimeList {
	channel := make(chan *AnimeList, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("AnimeList") {
			channel <- obj.(*AnimeList)
		}

		close(channel)
	}()

	return channel
}

// AllAnimeLists returns a slice of all anime.
func AllAnimeLists() ([]*AnimeList, error) {
	var all []*AnimeList

	stream := StreamAnimeLists()

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}

// GetAnimeList ...
func GetAnimeList(userID string) (*AnimeList, error) {
	animeList, err := DB.Get("AnimeList", userID)

	if err != nil {
		return nil, err
	}

	return animeList.(*AnimeList), nil
}
