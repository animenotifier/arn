package arn

import (
	"sort"
	"strings"

	"github.com/aerogo/flow"
)

// SearchResult ...
type SearchResult struct {
	obj        interface{}
	similarity float64
}

// Search is a fuzzy search.
func Search(term string, maxUsers, maxAnime, maxPosts, maxThreads int) ([]*User, []*Anime, []*Post, []*Thread) {
	term = strings.ToLower(term)

	if term == "" {
		return nil, nil, nil, nil
	}

	var userResults []*User
	var animeResults []*Anime
	var postResults []*Post
	var threadResults []*Thread

	flow.Parallel(func() {
		userResults = SearchUsers(term, maxUsers)
	}, func() {
		animeResults = SearchAnime(term, maxAnime)
	})

	return userResults, animeResults, postResults, threadResults
}

// SearchUsers searches all users.
func SearchUsers(term string, maxLength int) []*User {
	var results []*SearchResult

	for user := range StreamUsers() {
		text := strings.ToLower(user.Nick)

		// Similarity check
		similarity := AdvancedStringSimilarity(term, text)

		if similarity < MinimumStringSimilarity {
			continue
		}

		results = append(results, &SearchResult{
			obj:        user,
			similarity: similarity,
		})
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].similarity > results[j].similarity
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	// Final list
	final := make([]*User, len(results), len(results))

	for i, result := range results {
		final[i] = result.obj.(*User)
	}

	return final
}

// SearchAnime searches all anime.
func SearchAnime(term string, maxLength int) []*Anime {
	var results []*SearchResult

	check := func(text string) float64 {
		return AdvancedStringSimilarity(term, RemoveSpecialCharacters(strings.ToLower(text)))
	}

	add := func(anime *Anime, similarity float64) {
		similarity += float64(anime.Popularity.Total()) / 50.0

		results = append(results, &SearchResult{
			obj:        anime,
			similarity: similarity,
		})
	}

	for anime := range StreamAnime() {
		// Canonical title
		similarity := check(anime.Title.Canonical)

		if similarity >= MinimumStringSimilarity {
			add(anime, similarity)
			continue
		}

		// Synonyms
		for _, synonym := range anime.Title.Synonyms {
			similarity := check(synonym)

			if similarity >= MinimumStringSimilarity {
				add(anime, similarity)
				continue
			}
		}

		// Japanese
		similarity = check(anime.Title.Japanese)

		if similarity >= MinimumStringSimilarity {
			add(anime, similarity)
			continue
		}
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].similarity > results[j].similarity
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	// Final list
	final := make([]*Anime, len(results), len(results))

	for i, result := range results {
		final[i] = result.obj.(*Anime)
	}

	return final
}
