package arn

import "github.com/aerogo/nano"

// AnimeCharacters ...
type AnimeCharacters struct {
	AnimeID string            `json:"animeId" mainID:"true"`
	Items   []*AnimeCharacter `json:"items" editable:"true"`
}

// Contains tells you whether the given character ID exists.
func (characters *AnimeCharacters) Contains(characterID string) bool {
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
