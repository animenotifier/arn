package arn

import (
	"sort"
	"time"
)

const (
	currentlyAiringBonus      = 5.0
	longSummaryBonus          = 0.1
	popularityThreshold       = 5
	popularityPenalty         = 8.0
	watchingPopularityWeight  = 0.3
	completedPopularityWeight = 0.3
	plannedPopularityWeight   = 0.2
	droppedPopularityWeight   = -0.2
	agePenalty                = 11.0
	ageThreshold              = 6 * 30 * 24 * time.Hour
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

		scoreA := a.Rating.Overall
		scoreB := b.Rating.Overall

		if a.Status == "current" {
			scoreA += currentlyAiringBonus
		}

		if b.Status == "current" {
			scoreB += currentlyAiringBonus
		}

		if a.Popularity.Total() < popularityThreshold {
			scoreA -= popularityPenalty
		}

		if b.Popularity.Total() < popularityThreshold {
			scoreB -= popularityPenalty
		}

		if len(a.Summary) >= 140 {
			scoreA += longSummaryBonus
		}

		if len(b.Summary) >= 140 {
			scoreB += longSummaryBonus
		}

		// If we show currently running shows, rank shows that started a long time ago a bit lower
		if filterStatus == "current" {
			if a.StartDate != "" && time.Since(a.StartDateTime()) > ageThreshold {
				scoreA -= agePenalty
			}

			if b.StartDate != "" && time.Since(b.StartDateTime()) > ageThreshold {
				scoreB -= agePenalty
			}
		}

		scoreA += float64(a.Popularity.Watching) * watchingPopularityWeight
		scoreB += float64(b.Popularity.Watching) * watchingPopularityWeight

		scoreA += float64(a.Popularity.Planned) * plannedPopularityWeight
		scoreB += float64(b.Popularity.Planned) * plannedPopularityWeight

		scoreA += float64(a.Popularity.Completed) * completedPopularityWeight
		scoreB += float64(b.Popularity.Completed) * completedPopularityWeight

		scoreA += float64(a.Popularity.Dropped) * droppedPopularityWeight
		scoreB += float64(b.Popularity.Dropped) * droppedPopularityWeight

		return scoreA > scoreB
	})
}
