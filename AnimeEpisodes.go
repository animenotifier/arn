package arn

import (
	"bytes"
	"strconv"
	"strings"
	"sync"

	"github.com/aerogo/nano"
)

// AnimeEpisodes is a list of episodes for an anime.
type AnimeEpisodes struct {
	AnimeID string          `json:"animeId" mainID:"true"`
	Items   []*AnimeEpisode `json:"items" editable:"true"`

	sync.Mutex
}

// Link returns the link for that object.
func (episodes *AnimeEpisodes) Link() string {
	return "/anime/" + episodes.AnimeID + "/episodes"
}

// Find finds the given episode number.
func (episodes *AnimeEpisodes) Find(episodeNumber int) (*AnimeEpisode, int) {
	episodes.Lock()
	defer episodes.Unlock()

	for index, episode := range episodes.Items {
		if episode.Number == episodeNumber {
			return episode, index
		}
	}

	return nil, -1
}

// Merge combines the data of both episode slices to one.
func (episodes *AnimeEpisodes) Merge(b []*AnimeEpisode) {
	if b == nil {
		return
	}

	episodes.Lock()
	defer episodes.Unlock()

	for index, episode := range b {
		if index >= len(episodes.Items) {
			episodes.Items = append(episodes.Items, episode)
		} else {
			episodes.Items[index].Merge(episode)
		}
	}
}

// LastReversed returns the last n items in reversed order.
func (episodes *AnimeEpisodes) LastReversed(count int) []*AnimeEpisode {
	episodes.Lock()
	defer episodes.Unlock()

	items := episodes.Items[len(episodes.Items)-count:]

	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}

	return items
}

// AvailableCount counts the number of available episodes.
func (episodes *AnimeEpisodes) AvailableCount() int {
	episodes.Lock()
	defer episodes.Unlock()

	available := 0

	for _, episode := range episodes.Items {
		if len(episode.Links) > 0 {
			available++
		}
	}

	return available
}

// Anime returns the anime the episodes refer to.
func (episodes *AnimeEpisodes) Anime() *Anime {
	anime, _ := GetAnime(episodes.AnimeID)
	return anime
}

// String implements the default string serialization.
func (episodes *AnimeEpisodes) String() string {
	return episodes.Anime().String()
}

// ListString returns a text representation of the anime episodes.
func (episodes *AnimeEpisodes) ListString() string {
	episodes.Lock()
	defer episodes.Unlock()

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

// StreamAnimeEpisodes returns a stream of all anime episodes.
func StreamAnimeEpisodes() chan *AnimeEpisodes {
	channel := make(chan *AnimeEpisodes, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("AnimeEpisodes") {
			channel <- obj.(*AnimeEpisodes)
		}

		close(channel)
	}()

	return channel
}

// GetAnimeEpisodes ...
func GetAnimeEpisodes(id string) (*AnimeEpisodes, error) {
	obj, err := DB.Get("AnimeEpisodes", id)

	if err != nil {
		return nil, err
	}

	return obj.(*AnimeEpisodes), nil
}
