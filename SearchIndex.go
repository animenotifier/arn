package arn

import (
	"strings"

	"github.com/aerogo/aero"
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
	aero.Parallel(func() {
		// Search userResults
		var user *User

		userSearchIndex, err := GetSearchIndex("User")

		if err != nil {
			return
		}

		for name, id := range userSearchIndex.TextToID {
			if strings.Index(name, term) != -1 {
				user, err = GetUser(id)

				if err != nil {
					continue
				}

				userResults = append(userResults, user)

				if len(userResults) >= maxUsers {
					break
				}
			}
		}
	}, func() {
		// Search anime
		var anime *Anime

		animeSearchIndex, err := GetSearchIndex("Anime")

		if err != nil {
			return
		}

		for title, id := range animeSearchIndex.TextToID {
			if strings.Index(title, term) != -1 {
				anime, err = GetAnime(id)

				if err != nil {
					continue
				}

				animeResults = append(animeResults, anime)

				if len(animeResults) >= maxAnime {
					break
				}
			}
		}
	})

	return userResults, animeResults
}
