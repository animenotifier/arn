package arn

// AnimeEpisodes ...
type AnimeEpisodes struct {
	AnimeID string          `json:"animeId"`
	Items   []*AnimeEpisode `json:"items"`
}

// Merge combines the data of both episode slices to one.
func (episodes *AnimeEpisodes) Merge(b []*AnimeEpisode) {
	if b == nil {
		return
	}

	for index, episode := range b {
		if index >= len(episodes.Items) {
			episodes.Items = append(episodes.Items, episode)
		} else {
			episodes.Items[index].Merge(episode)
		}
	}
}

// AvailableCount counts the number of available episodes.
func (episodes *AnimeEpisodes) AvailableCount() int {
	available := 0

	for _, episode := range episodes.Items {
		if len(episode.Links) > 0 {
			available++
		}
	}

	return available
}

// Save saves the episodes in the database.
func (episodes *AnimeEpisodes) Save() error {
	return DB.Set("AnimeEpisodes", episodes.AnimeID, episodes)
}

// GetAnimeEpisodes ...
func GetAnimeEpisodes(id string) (*AnimeEpisodes, error) {
	obj, err := DB.Get("AnimeEpisodes", id)

	if err != nil {
		return nil, err
	}

	return obj.(*AnimeEpisodes), nil
}
