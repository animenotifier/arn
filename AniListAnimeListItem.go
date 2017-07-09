package arn

import "github.com/animenotifier/anilist"

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
