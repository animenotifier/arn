package arn

import (
	"github.com/aerogo/nano"
)

// Character ...
type Character struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Image       string                `json:"image"`
	Description string                `json:"description"`
	Attributes  []*CharacterAttribute `json:"attributes"`
	// Name        *CharacterName        `json:"name"`
	// Mappings    []*Mapping            `json:"mappings"`
}

// Link ...
func (character *Character) Link() string {
	return "/character/" + character.ID
}

// Anime returns a list of all anime the character appears in.
func (character *Character) Anime() []*Anime {
	var results []*Anime

	for animeCharacters := range StreamAnimeCharacters() {
		if animeCharacters.Contains(character.ID) {
			anime, err := GetAnime(animeCharacters.AnimeID)

			if err != nil {
				continue
			}

			results = append(results, anime)
		}
	}

	return results
}

// GetCharacter ...
func GetCharacter(id string) (*Character, error) {
	obj, err := DB.Get("Character", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Character), nil
}

// Quotes returns the list of quotes for this character.
func (character *Character) Quotes() []*Quote {
	return FilterQuotes(func(quote *Quote) bool {
		return !quote.IsDraft && quote.CharacterID == character.ID
	})
}

// StreamCharacters returns a stream of all characters.
func StreamCharacters() chan *Character {
	channel := make(chan *Character, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Character") {
			channel <- obj.(*Character)
		}

		close(channel)
	}()

	return channel
}

// AllCharacters returns a slice of all characters.
func AllCharacters() []*Character {
	var all []*Character

	stream := StreamCharacters()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}
