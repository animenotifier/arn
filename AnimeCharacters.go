package arn

import (
	"sync"

	"github.com/aerogo/nano"
)

// AnimeCharacters ...
type AnimeCharacters struct {
	AnimeID string            `json:"animeId" mainID:"true"`
	Items   []*AnimeCharacter `json:"items" editable:"true"`

	sync.Mutex
}

// Anime returns the anime the characters refer to.
func (characters *AnimeCharacters) Anime() *Anime {
	anime, _ := GetAnime(characters.AnimeID)
	return anime
}

// String implements the default string serialization.
func (characters *AnimeCharacters) String() string {
	return characters.Anime().String()
}

// Contains tells you whether the given character ID exists.
func (characters *AnimeCharacters) Contains(characterID string) bool {
	characters.Lock()
	defer characters.Unlock()

	for _, item := range characters.Items {
		if item.CharacterID == characterID {
			return true
		}
	}

	return false
}

// GetAnimeCharacters ...
func GetAnimeCharacters(animeID string) (*AnimeCharacters, error) {
	obj, err := DB.Get("AnimeCharacters", animeID)

	if err != nil {
		return nil, err
	}

	return obj.(*AnimeCharacters), nil
}

// StreamAnimeCharacters returns a stream of all anime characters.
func StreamAnimeCharacters() chan *AnimeCharacters {
	channel := make(chan *AnimeCharacters, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("AnimeCharacters") {
			channel <- obj.(*AnimeCharacters)
		}

		close(channel)
	}()

	return channel
}
