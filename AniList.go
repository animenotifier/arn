package arn

import (
	"strconv"
	"strings"

	"github.com/animenotifier/anilist"
)

// FindAniListAnime tries to find an AniListAnime in our Anime database.
func FindAniListAnime(search *anilist.Anime, allAnime []*Anime) *Anime {
	var mostSimilar *Anime
	var similarity float64

	searchID := strconv.Itoa(search.ID)

	for _, anime := range allAnime {
		if anime.GetMapping("anilist/anime") == searchID {
			return anime
		}

		japanese := strings.Replace(anime.Title.Japanese, "2ndシーズン", "2", 1)
		romaji := strings.Replace(anime.Title.Romaji, " 2nd Season", " 2", 1)

		search.TitleJapanese = strings.TrimSpace(strings.Replace(search.TitleJapanese, "2ndシーズン", "2", 1))
		search.TitleRomaji = strings.TrimSpace(strings.Replace(search.TitleRomaji, " 2nd Season", " 2", 1))

		titleSimilarity := 0.0

		if anime.Title.Romaji != "" {
			titleSimilarity = StringSimilarity(anime.Title.Romaji, search.TitleRomaji)
		} else {
			titleSimilarity = StringSimilarity(anime.Title.Canonical, search.TitleRomaji)
		}

		if strings.ToLower(japanese) == strings.ToLower(search.TitleJapanese) {
			titleSimilarity += 1.0
		}

		if strings.ToLower(romaji) == strings.ToLower(search.TitleRomaji) {
			titleSimilarity += 1.0
		}

		if strings.ToLower(anime.Title.English) == strings.ToLower(search.TitleEnglish) {
			titleSimilarity += 1.0
		}

		if titleSimilarity > similarity {
			mostSimilar = anime
			similarity = titleSimilarity
		}
	}

	if mostSimilar.EpisodeCount != search.TotalEpisodes {
		similarity -= 0.02
	}

	if similarity >= 0.92 && mostSimilar.GetMapping("anilist/anime") == "" {
		// fmt.Printf("MATCH:    %s => %s (%.2f)\n", search.TitleRomaji, mostSimilar.Title.Romaji, similarity)
		mostSimilar.SetMapping("anilist/anime", strconv.Itoa(search.ID))
		mostSimilar.Save()
		return mostSimilar
	}

	// color.Red("MISMATCH: %s => %s (%.2f)", search.TitleRomaji, mostSimilar.Title.Romaji, similarity)
	return nil
}

// AniListAnimeListStatus returns the ARN version of the anime status.
func AniListAnimeListStatus(item *anilist.AnimeListItem) string {
	switch item.ListStatus {
	case "watching":
		return AnimeListStatusWatching
	case "completed":
		return AnimeListStatusCompleted
	case "plan to watch":
		return AnimeListStatusPlanned
	case "on-hold":
		return AnimeListStatusHold
	case "dropped":
		return AnimeListStatusDropped
	default:
		return AnimeListStatusPlanned
	}
}
