package arn

// FindMyAnimeListAnime tries to find a MyAnimeList anime in our anime database.
func FindMyAnimeListAnime(searchID string, allAnime []*Anime) *Anime {
	for _, anime := range allAnime {
		if anime.GetMapping("myanimelist/anime") == searchID {
			return anime
		}
	}

	return nil
}
