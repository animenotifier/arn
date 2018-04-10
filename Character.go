package arn

import (
	"fmt"

	"github.com/aerogo/nano"
)

// Character ...
type Character struct {
	ID          string                `json:"id"`
	Name        CharacterName         `json:"name"`
	Image       CharacterImage        `json:"image"`
	Description string                `json:"description"`
	Attributes  []*CharacterAttribute `json:"attributes"`
	MainQuoteID string                `json:"mainQuoteId"`
	HasMappings
}

// Link ...
func (character *Character) Link() string {
	return "/character/" + character.ID
}

// MainQuote ...
func (character *Character) MainQuote() *Quote {
	quote, _ := GetQuote(character.MainQuoteID)
	return quote
}

// ImageLink ...
func (character *Character) ImageLink(size string) string {
	return fmt.Sprintf("//media.kitsu.io/characters/images/%s/original%s?%d", character.ID, character.Image.Extension, character.Image.LastModified)
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
