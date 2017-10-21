package arn

import (
	"bytes"
	"strconv"
	"strings"
)

// AnimeEpisodes ...
type AnimeEpisodes struct {
	AnimeID string             `json:"animeId"`
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

// String returns a text representation of the anime episodes.
func (episodes *AnimeEpisodes) String() string {
	b := bytes.Buffer{}

	for _, episode := range episodes.Items {
		b.WriteString(strconv.Itoa(episode.Number))
		b.WriteString(" | ")
		b.WriteString(episode.Title.Japanese)
		b.WriteString(" | ")
		b.WriteString(episode.AiringDate.StartDateHuman())
		b.WriteByte('\n')
	}

	return strings.TrimRight(b.String(), "\n")
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
