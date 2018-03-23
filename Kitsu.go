package arn

// FindKitsuAnime tries to find a Kitsu anime in our anime database.
func FindKitsuAnime(searchID string) *Anime {
	anime, _ := GetAnime(searchID)
	return anime

	// For the future:

	// for _, anime := range allAnime {
	// 	if anime.GetMapping("kitsu/anime") == searchID {
	// 		return anime
	// 	}
	// }

	// return nil
}
