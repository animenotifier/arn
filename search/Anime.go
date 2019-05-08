package search

import (
	"sort"
	"strings"

	"github.com/animenotifier/arn"
	"github.com/animenotifier/arn/stringutils"
)

// Anime searches all anime.
func Anime(originalTerm string, maxLength int) []*arn.Anime {
	term := strings.ToLower(stringutils.RemoveSpecialCharacters(originalTerm))
	results := make([]*Result, 0, maxLength)

	check := func(text string) float64 {
		if text == "" {
			return 0
		}

		return stringutils.AdvancedStringSimilarity(term, strings.ToLower(stringutils.RemoveSpecialCharacters(text)))
	}

	add := func(anime *arn.Anime, similarity float64) {
		similarity += float64(anime.Popularity.Total()) * popularityDamping

		results = append(results, &Result{
			obj:        anime,
			similarity: similarity,
		})
	}

	for anime := range arn.StreamAnime() {
		if anime.ID == originalTerm {
			return []*arn.Anime{anime}
		}

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

		// Romaji
		similarity = check(anime.Title.Romaji)

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
	final := make([]*arn.Anime, len(results))

	for i, result := range results {
		final[i] = result.obj.(*arn.Anime)
	}

	return final
}
