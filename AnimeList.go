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
	if list.user == nil {
		list.user, _ = GetUser(list.UserID)
	}

	return list.user
}

// Sort ...
func (list *AnimeList) Sort() {
	sort.Slice(list.Items, func(i, j int) bool {
		a := list.Items[i]
		b := list.Items[j]

		if a.Status != AnimeListStatusWatching && b.Status != AnimeListStatusWatching {
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

// PrefetchAnime loads all the anime objects from the list into memory.
func (list *AnimeList) PrefetchAnime() {
	animeIDList := make([]string, len(list.Items), len(list.Items))

	for i, item := range list.Items {
		animeIDList[i] = item.AnimeID
	}

	// Prefetch anime objects
	animeObjects, _ := DB.GetMany("Anime", animeIDList)
	prefetchedAnime := animeObjects.([]*Anime)

	for i, anime := range prefetchedAnime {
		list.Items[i].anime = anime
	}
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
func GetAnimeList(userID string) (*AnimeList, error) {
	animeList := &AnimeList{
		UserID: userID,
		Items:  []*AnimeListItem{},
	}

	m, err := DB.GetMap("AnimeList", userID)

	if err != nil {
		return nil, err
	}

	itemList := m["items"].([]interface{})

	for _, itemMap := range itemList {
		item := itemMap.(map[interface{}]interface{})
		ratingMap := item["rating"].(map[interface{}]interface{})
		newItem := &AnimeListItem{
			AnimeID:      item["animeId"].(string),
			Status:       item["status"].(string),
			Episodes:     item["episodes"].(int),
			Notes:        item["notes"].(string),
			RewatchCount: item["rewatchCount"].(int),
			Private:      item["private"].(int) != 0,
			Edited:       item["edited"].(string),
			Created:      item["created"].(string),
			Rating: &AnimeRating{
				Overall:    ratingMap["overall"].(float64),
				Story:      ratingMap["story"].(float64),
				Visuals:    ratingMap["visuals"].(float64),
				Soundtrack: ratingMap["soundtrack"].(float64),
			},
		}

		animeList.Items = append(animeList.Items, newItem)
	}

	return animeList, nil
}
