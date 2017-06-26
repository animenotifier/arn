package arn

import (
	"sort"
	"strings"

	"github.com/aerogo/flow"
)

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

	// Search everything in parallel
	flow.Parallel(func() {
		// Search userResults
		var user *User

		userSearchIndex, err := GetSearchIndex("User")

		if err != nil {
			return
		}

		textToID := userSearchIndex.TextToID

		// Keys
		keys := make([]string, len(textToID))
		count := 0
		for name := range textToID {
			keys[count] = name
			count++
		}

		sort.Slice(keys, func(i, j int) bool {
			return StringSimilarity(term, keys[i]) > StringSimilarity(term, keys[j])
		})

		if len(keys) >= maxUsers {
			keys = keys[:maxUsers]
		}

		for _, key := range keys {
			user, err = GetUser(textToID[key])

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

		// Keys
		keys := make([]string, len(textToID))
		count := 0
		for name := range textToID {
			keys[count] = name
			count++
		}

		sort.Slice(keys, func(i, j int) bool {
			return StringSimilarity(term, keys[i]) > StringSimilarity(term, keys[j])
		})

		if len(keys) >= maxAnime {
			keys = keys[:maxAnime]
		}

		for _, key := range keys {
			anime, err = GetAnime(textToID[key])

			if err != nil {
				continue
			}

			animeResults = append(animeResults, anime)
		}
	})

	return userResults, animeResults
}
