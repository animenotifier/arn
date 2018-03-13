package arn

import (
	"sort"
)

// SortCompaniesPopularFirst ...
func SortCompaniesPopularFirst(companies []*Company) {
	// Generate company ID to popularity map
	popularity := map[string]int{}

	for anime := range StreamAnime() {
		for _, studio := range anime.Studios() {
			popularity[studio.ID] += anime.Popularity.Watching + anime.Popularity.Completed
		}
	}

	// Sort by using the popularity map
	sort.Slice(companies, func(i, j int) bool {
		a := companies[i]
		b := companies[j]

		aPopularity := popularity[a.ID]
		bPopularity := popularity[b.ID]

		if aPopularity == bPopularity {
			return a.Name.English < b.Name.English
		}

		return aPopularity > bPopularity
	})
}
