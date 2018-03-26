package arn

import (
	"sort"
	"sync"

	"github.com/aerogo/nano"
)

// AnimeRelations ...
type AnimeRelations struct {
	AnimeID string           `json:"animeId" mainID:"true"`
	Items   []*AnimeRelation `json:"items" editable:"true"`

	sync.Mutex
}

// SortByStartDate ...
func (relations *AnimeRelations) SortByStartDate() {
	relations.Lock()
	defer relations.Unlock()

	sort.Slice(relations.Items, func(i, j int) bool {
		a := relations.Items[i].Anime()
		b := relations.Items[j].Anime()

		if a.StartDate == b.StartDate {
			return a.Title.Canonical < b.Title.Canonical
		}

		return a.StartDate < b.StartDate
	})
}

// Anime returns the anime the relations list refers to.
func (relations *AnimeRelations) Anime() *Anime {
	anime, _ := GetAnime(relations.AnimeID)
	return anime
}

// String implements the default string serialization.
func (relations *AnimeRelations) String() string {
	return relations.Anime().String()
}

// Remove removes the anime ID from the relations.
func (relations *AnimeRelations) Remove(animeID string) bool {
	relations.Lock()
	defer relations.Unlock()

	for index, item := range relations.Items {
		if item.AnimeID == animeID {
			relations.Items = append(relations.Items[:index], relations.Items[index+1:]...)
			return true
		}
	}

	return false
}

// GetAnimeRelations ...
func GetAnimeRelations(animeID string) (*AnimeRelations, error) {
	obj, err := DB.Get("AnimeRelations", animeID)

	if err != nil {
		return nil, err
	}

	return obj.(*AnimeRelations), nil
}

// StreamAnimeRelations returns a stream of all anime relations.
func StreamAnimeRelations() chan *AnimeRelations {
	channel := make(chan *AnimeRelations, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("AnimeRelations") {
			channel <- obj.(*AnimeRelations)
		}

		close(channel)
	}()

	return channel
}
