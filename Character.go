package arn

import (
	"fmt"

	"github.com/aerogo/nano"
)

// Character ...
type Character struct {
	ID          string                `json:"id"`
	Name        CharacterName         `json:"name" editable:"true"`
	Image       CharacterImage        `json:"image"`
	MainQuoteID string                `json:"mainQuoteId" editable:"true"`
	Description string                `json:"description" editable:"true" type:"textarea"`
	Spoilers    []Spoiler             `json:"spoilers" editable:"true"`
	Attributes  []*CharacterAttribute `json:"attributes" editable:"true"`
	HasMappings
}

// Link ...
func (character *Character) Link() string {
	return "/character/" + character.ID
}

// String returns the canonical name of the character.
func (character *Character) String() string {
	return character.Name.Canonical
}

// MainQuote ...
func (character *Character) MainQuote() *Quote {
	quote, _ := GetQuote(character.MainQuoteID)
	return quote
}

// AverageColor returns the average color of the image.
func (character *Character) AverageColor() string {
	color := character.Image.AverageColor

	if color.Hue == 0 && color.Saturation == 0 && color.Lightness == 0 {
		return ""
	}

	return color.String()
}

// ImageLink ...
func (character *Character) ImageLink(size string) string {
	extension := ".jpg"

	if size == "original" {
		extension = character.Image.Extension
	}

	return fmt.Sprintf("//%s/images/characters/%s/%s%s?%v", MediaHost, size, character.ID, extension, character.Image.LastModified)
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

// FilterCharacters filters all characters by a custom function.
func FilterCharacters(filter func(*Character) bool) []*Character {
	var filtered []*Character

	channel := DB.All("Character")

	for obj := range channel {
		realObject := obj.(*Character)

		if filter(realObject) {
			filtered = append(filtered, realObject)
		}
	}

	return filtered
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
