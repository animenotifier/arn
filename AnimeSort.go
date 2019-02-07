package arn

import (
	"fmt"
	"sort"
	"time"
)

const (
	currentlyAiringBonus      = 5.0
	longSummaryBonus          = 0.1
	popularityThreshold       = 5
	popularityPenalty         = 8.0
	watchingPopularityWeight  = 0.07
	completedPopularityWeight = watchingPopularityWeight
	plannedPopularityWeight   = watchingPopularityWeight * (2.0 / 3.0)
	droppedPopularityWeight   = -plannedPopularityWeight
	visualsWeight             = 0.0075
	storyWeight               = 0.0075
	soundtrackWeight          = 0.0075
	movieBonus                = 0.28
	agePenalty                = 11.0
	ageThreshold              = 6 * 30 * 24 * time.Hour
	sortByTitle               = "Title"
	sortByStartDate           = "StartDate"
	sortByEpisodeCount        = "EpisodeCount"
	sortByEpisodeLength       = "EpisodeLength"
)

// SortAnimeByPopularity sorts the given slice of anime by popularity.
func SortAnimeByPopularity(animes []*Anime) {
	sort.Slice(animes, func(i, j int) bool {
		aPopularity := animes[i].Popularity.Total()
		bPopularity := animes[j].Popularity.Total()

		if aPopularity == bPopularity {
			return animes[i].Title.Canonical < animes[j].Title.Canonical
		}

		return aPopularity > bPopularity
	})
}

// SortAnimeByQuality sorts the given slice of anime by quality.
func SortAnimeByQuality(animes []*Anime) {
	SortAnimeByQualityDetailed(animes, "")
}

// SortAnimeByQualityDetailed sorts the given slice of anime by quality.
func SortAnimeByQualityDetailed(animes []*Anime, filterStatus string) {
	sort.Slice(animes, func(i, j int) bool {
		a := animes[i]
		b := animes[j]
		return CompareAnimeByQualityDetailed(a, b, filterStatus)
	})
}

// SortAnimeWithAlgo sorts the given slice of anime by the specified algo.
func SortAnimeWithAlgo(animes []*Anime, filterStatus string, algo string, user *User) {
	sort.Slice(animes, func(i, j int) bool {
		anime := animes[i]
		otherAnime := animes[j]

		switch algo {
		case sortByTitle:
			return CompareAnimeByTile(anime, otherAnime, user)
		case sortByStartDate:
			return CompareAnimeByStartDate(anime, otherAnime)
		case sortByEpisodeCount:
			return CompareAnimeByEpisodeCount(anime, otherAnime)
		case sortByEpisodeLength:
			return CompareAnimeByEpisodeLength(anime, otherAnime)
		default:
			return CompareAnimeByQualityDetailed(anime, otherAnime, filterStatus)
		}
	})
}

// CompareAnimeByTile compare 2 Animes by their title. Use the user preferred title if a user is passed; otherwise, it uses the canonical one.
func CompareAnimeByTile(anime *Anime, otherAnime *Anime, user *User) bool {
	if anime.Title.ByUser(user) == otherAnime.Title.ByUser(user) {
		return anime.Rating.Overall < otherAnime.Rating.Overall
	}

	return anime.Title.ByUser(user) < otherAnime.Title.ByUser(user)
}

// CompareAnimeByStartDate compare 2 Animes by their start date. return true if the first anime was aired strictly before the second.
func CompareAnimeByStartDate(anime *Anime, otherAnime *Anime) bool {
	if anime.StartDate == otherAnime.StartDate {
		return anime.Title.Canonical < otherAnime.Title.Canonical
	}

	return anime.StartDate > otherAnime.StartDate
}

// CompareAnimeByEpisodeCount compare 2 Animes by their episode counts. return true if the first anime has strictly more episodes the second.
func CompareAnimeByEpisodeCount(anime *Anime, otherAnime *Anime) bool {
	if anime.EpisodeCount == otherAnime.EpisodeCount {
		return anime.Title.Canonical < otherAnime.Title.Canonical
	}

	return anime.EpisodeCount > otherAnime.EpisodeCount
}

// CompareAnimeByEpisodeLength compare 2 Animes by their episode length. return true if the first anime episode are strictly longer the second.
func CompareAnimeByEpisodeLength(anime *Anime, otherAnime *Anime) bool {
	if anime.EpisodeLength == otherAnime.EpisodeLength {
		return anime.EpisodeLength < otherAnime.EpisodeLength
	}

	return anime.EpisodeLength > otherAnime.EpisodeLength
}

// CompareAnimeByQualityDetailed compare 2 Anime by their quality. return true if the first anime has a strictly better score than second.
func CompareAnimeByQualityDetailed(anime *Anime, otherAnime *Anime, filterStatus string) bool {
	scoreA := anime.Score()
	scoreB := otherAnime.Score()

	// If we show currently running shows, rank shows that started a long time ago a bit lower
	if filterStatus == "current" {
		if anime.StartDate != "" && time.Since(anime.StartDateTime()) > ageThreshold {
			scoreA -= agePenalty
		}

		if otherAnime.StartDate != "" && time.Since(otherAnime.StartDateTime()) > ageThreshold {
			scoreB -= agePenalty
		}
	}

	if scoreA == scoreB {
		return anime.Title.Canonical < otherAnime.Title.Canonical
	}

	return scoreA > scoreB
}

// Score returns the score used for the anime ranking.
func (anime *Anime) Score() float64 {
	score := anime.Rating.Overall
	score += anime.Rating.Story * storyWeight
	score += anime.Rating.Visuals * visualsWeight
	score += anime.Rating.Soundtrack * soundtrackWeight

	score += float64(anime.Popularity.Watching) * watchingPopularityWeight
	score += float64(anime.Popularity.Planned) * plannedPopularityWeight
	score += float64(anime.Popularity.Completed) * completedPopularityWeight
	score += float64(anime.Popularity.Dropped) * droppedPopularityWeight

	if anime.Status == "current" {
		score += currentlyAiringBonus
	}

	if anime.Type == "movie" {
		score += movieBonus
	}

	if anime.Popularity.Total() < popularityThreshold {
		score -= popularityPenalty
	}

	if len(anime.Summary) >= 140 {
		score += longSummaryBonus
	}

	return score
}

// ScoreHumanReadable returns the score used for the anime ranking in human readable format.
func (anime *Anime) ScoreHumanReadable() string {
	return fmt.Sprintf("%.1f", anime.Score())
}
