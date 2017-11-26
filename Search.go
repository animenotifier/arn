package arn

import (
	"sort"
	"strings"

	"github.com/aerogo/flow"
)

// MinimumStringSimilarity is the minimum JaroWinkler distance we accept for search results.
const MinimumStringSimilarity = 0.89

// popularityDamping reduces the factor of popularity in search results.
const popularityDamping = 0.048

// SearchResult ...
type SearchResult struct {
	obj        interface{}
	similarity float64
}

// Search is a fuzzy search.
func Search(term string, maxUsers, maxAnime, maxPosts, maxThreads, maxTracks int, maxCharacters int) ([]*User, []*Anime, []*Post, []*Thread, []*SoundTrack, []*Character) {
	term = strings.ToLower(term)

	if term == "" {
		return nil, nil, nil, nil, nil, nil
	}

	var userResults []*User
	var animeResults []*Anime
	var postResults []*Post
	var threadResults []*Thread
	var trackResults []*SoundTrack
	var characterResults []*Character

	flow.Parallel(func() {
		userResults = SearchUsers(term, maxUsers)
	}, func() {
		animeResults = SearchAnime(term, maxAnime)
	}, func() {
		postResults = SearchPosts(term, maxPosts)
	}, func() {
		threadResults = SearchThreads(term, maxThreads)
	}, func() {
		trackResults = SearchSoundTracks(term, maxTracks)
	}, func() {
		characterResults = SearchCharacters(term, maxCharacters)
	})

	return userResults, animeResults, postResults, threadResults, trackResults, characterResults
}

// SearchCharacters searches all characters.
func SearchCharacters(term string, maxLength int) []*Character {
	var results []*Character

	for character := range StreamCharacters() {
		text := strings.ToLower(character.Name)

		if !strings.Contains(text, term) {
			continue
		}

		results = append(results, character)
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name > results[j].Name
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	return results
}

// SearchSoundTracks searches all soundtracks.
func SearchSoundTracks(term string, maxLength int) []*SoundTrack {
	var results []*SearchResult

	for track := range StreamSoundTracks() {
		text := strings.ToLower(track.Title)

		// Similarity check
		similarity := AdvancedStringSimilarity(term, text)

		if similarity >= MinimumStringSimilarity {
			results = append(results, &SearchResult{
				obj:        track,
				similarity: similarity,
			})
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
	final := make([]*SoundTrack, len(results), len(results))

	for i, result := range results {
		final[i] = result.obj.(*SoundTrack)
	}

	return final
}

// SearchPosts searches all posts.
func SearchPosts(term string, maxLength int) []*Post {
	var results []*Post

	for post := range StreamPosts() {
		text := strings.ToLower(post.Text)

		if !strings.Contains(text, term) {
			continue
		}

		results = append(results, post)
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].Created > results[j].Created
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	return results
}

// SearchThreads searches all threads.
func SearchThreads(term string, maxLength int) []*Thread {
	var results []*Thread

	for thread := range StreamThreads() {
		text := strings.ToLower(thread.Text)

		if strings.Contains(text, term) {
			results = append(results, thread)
			continue
		}

		text = strings.ToLower(thread.Title)

		if strings.Contains(text, term) {
			results = append(results, thread)
			continue
		}
	}

	// Sort
	sort.Slice(results, func(i, j int) bool {
		return results[i].Created > results[j].Created
	})

	// Limit
	if len(results) >= maxLength {
		results = results[:maxLength]
	}

	return results
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
		similarity += float64(anime.Popularity.Total()) * popularityDamping

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

		// English
		similarity = check(anime.Title.English)

		if similarity >= MinimumStringSimilarity {
			add(anime, similarity)
			continue
		}

		// Synonyms
		for _, synonym := range anime.Title.Synonyms {
			similarity := check(synonym)

			if similarity >= MinimumStringSimilarity {
				add(anime, similarity)
				goto nextAnime
			}
		}

		// Japanese
		similarity = check(anime.Title.Japanese)

		if similarity >= MinimumStringSimilarity {
			add(anime, similarity)
			continue
		}

	nextAnime:
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
