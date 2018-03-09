package arn

import (
	"sort"
	"sync"
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

// GetAnimeRelations ...
func GetAnimeRelations(animeID string) (*AnimeRelations, error) {
	obj, err := DB.Get("AnimeRelations", animeID)

	if err != nil {
		return nil, err
	}

	return obj.(*AnimeRelations), nil
}
