package arn

import (
	"sort"
	"strings"

	"github.com/aerogo/flow"
)

// MinimumStringSimilarity is the minimum JaroWinkler distance we accept for search results.
const MinimumStringSimilarity = 0.9

// SearchIndex ...
type SearchIndex struct {
	TextToID map[string]string `json:"textToId"`
}

// NewSearchIndex ...
func NewSearchIndex() *SearchIndex {
	return &SearchIndex{
		TextToID: make(map[string]string),
	}
}

// GetSearchIndex ...
func GetSearchIndex(id string) (*SearchIndex, error) {
	obj, err := DB.Get("SearchIndex", id)
	return obj.(*SearchIndex), err
}

// Search is a fuzzy search.
func Search(term string, maxUsers, maxAnime int) ([]*User, []*Anime) {
	term = strings.ToLower(term)

	if term == "" {
		return nil, nil
	}

	var userResults []*User
	var animeResults []*Anime

	type SearchItem struct {
		text       string
		similarity float64
	}

	// Search everything in parallel
	flow.Parallel(func() {
		// Search userResults
		var user *User

		userSearchIndex, err := GetSearchIndex("User")

		if err != nil {
			return
		}

		textToID := userSearchIndex.TextToID

		// Search items
		items := make([]*SearchItem, 0)

		for name := range textToID {
			s := StringSimilarity(term, name)

			if strings.Index(name, term) != -1 {
				s += 0.5
			}

			if s < MinimumStringSimilarity {
				continue
			}

			items = append(items, &SearchItem{
				text:       name,
				similarity: s,
			})
		}

		// Sort
		sort.Slice(items, func(i, j int) bool {
			return items[i].similarity > items[j].similarity
		})

		// Limit
		if len(items) >= maxUsers {
			items = items[:maxUsers]
		}

		// Fetch data
		for _, item := range items {
			user, err = GetUser(textToID[item.text])

			if err != nil {
				continue
			}

			userResults = append(userResults, user)
		}
	}, func() {
		// Search anime
		var anime *Anime

		animeSearchIndex, err := GetSearchIndex("Anime")

		if err != nil {
			return
		}

		textToID := animeSearchIndex.TextToID

		// Search items
		items := make([]*SearchItem, 0)
		animeIDAdded := map[string]*SearchItem{}

		for name, id := range textToID {
			s := StringSimilarity(term, name)

			if strings.Index(name, term) != -1 {
				s += 0.5
			}

			if s < MinimumStringSimilarity {
				continue
			}

			addedEntry, found := animeIDAdded[id]

			// Skip existing anime IDs
			if found {
				// But update existing entry with new similarity if it's higher
				if s > addedEntry.similarity {
					addedEntry.similarity = s
				}

				continue
			}

			item := &SearchItem{
				text:       name,
				similarity: s,
			}
			items = append(items, item)

			animeIDAdded[id] = item
		}

		// Sort
		sort.Slice(items, func(i, j int) bool {
			return items[i].similarity > items[j].similarity
		})

		// Limit
		if len(items) >= maxAnime {
			items = items[:maxAnime]
		}

		// Fetch data
		for _, item := range items {
			anime, err = GetAnime(textToID[item.text])

			if err != nil {
				continue
			}

			animeResults = append(animeResults, anime)
		}
	})

	return userResults, animeResults
}
