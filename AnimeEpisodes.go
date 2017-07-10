package arn

// AnimeEpisodes ...
type AnimeEpisodes struct {
	AnimeID string          `json:"animeId"`
	Items   []*AnimeEpisode `json:"items"`
}

// Save saves the episodes in the database.
func (episodes *AnimeEpisodes) Save() error {
	return DB.Set("AnimeEpisodes", episodes.AnimeID, episodes)
}
