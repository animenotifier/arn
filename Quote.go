package arn

import (
	"errors"

	"github.com/aerogo/nano"
	"github.com/fatih/color"
	"sort"
)

// Quote ...
type Quote struct {
	ID          string   `json:"id"`
	Description string   `json:"description" editable:"true" type:"textarea"`
	CharacterId string   `json:"characterId" editable:"true"`
	AnimeId     string   `json:"animeId" editable:"true"`
	Episode     int      `json:"episode" editable:"true"`
	Time        int      `json:"time" editable:"true"`
	Likes       []string `json:"likes"`
	IsDraft     bool     `json:"isDraft"`
	Created     string   `json:"created"`
	CreatedBy   string   `json:"createdBy"`
	Edited      string   `json:"edited"`
	EditedBy    string   `json:"editedBy"`
}

// Link returns a single quote.
func (quote *Quote) Link() string {
	return "/quote/" + quote.ID
}

// Creator returns the user who created this quote.
func (quote *Quote) Creator() *User {
	user, _ := GetUser(quote.CreatedBy)
	return user
}

// EditedByUser returns the user who edited this quote.
func (quote *Quote) EditedByUser() *User {
	user, _ := GetUser(quote.EditedBy)
	return user
}

// Publish ...
func (quote *Quote) Publish() error {
	// No draft
	if !quote.IsDraft {
		return errors.New("Not a draft")
	}

	// No description
	if quote.Description == "" {
		return errors.New("A description is required")
	}

	// No character
	if quote.CharacterId == "" {
		return errors.New("A character is required")
	}

	// Episode is provided without the anime
	if quote.AnimeId == "" && quote.Episode != 0 {
		return errors.New("An anime is required before adding an episode")
	}

	// The time is provided without the Episode
	if quote.Episode == 0 && quote.Time != 0 {
		return errors.New("An episode is required before adding a time")
	}

	draftIndex, err := GetDraftIndex(quote.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.QuoteID == "" {
		return errors.New("Quote draft doesn't exist in the user draft index")
	}

	character, characterErr := GetCharacter(quote.CharacterId)

	if characterErr != nil {
		return errors.New("Character does not exist")
	}

	// Append to quotes Ids
	character.QuotesIds = append(character.QuotesIds, quote.ID)

	// Save the character
	character.Save()

	quote.IsDraft = false
	draftIndex.QuoteID = ""
	draftIndex.Save()
	return nil
}

// Unpublish ...
func (quote *Quote) Unpublish() error {
	draftIndex, err := GetDraftIndex(quote.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.QuoteID != "" {
		return errors.New("You still have an unfinished draft")
	}

	quote.IsDraft = true
	draftIndex.QuoteID = quote.ID
	draftIndex.Save()
	return nil
}

// Like adds an user to the quote's Likes array if they aren't already in it.
func (quote *Quote) Like(userID string) {
	for _, id := range quote.Likes {
		if id == userID {
			return
		}
	}

	quote.Likes = append(quote.Likes, userID)
}

// Unlike removes the user from the quote's Likes array if they are in it.
func (quote *Quote) Unlike(userID string) {
	for index, id := range quote.Likes {
		if id == userID {
			quote.Likes = append(quote.Likes[:index], quote.Likes[index+1:]...)
			return
		}
	}
}

// LikedBy checks to see if the user has liked the quote.
func (quote *Quote) LikedBy(userID string) bool {
	for _, id := range quote.Likes {
		if id == userID {
			return true
		}
	}

	return false
}

// GetQuote returns a single quote.
func GetQuote(id string) (*Quote, error) {
	obj, err := DB.Get("Quote", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Quote), nil
}

// StreamCompanies returns a stream of all quotes.
func StreamQuotes() chan *Quote {
	channel := make(chan *Quote, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Quote") {
			channel <- obj.(*Quote)
		}

		close(channel)
	}()

	return channel
}

// AllQuotes returns a slice of all quotes.
func AllQuotes() []*Quote {
	var all []*Quote

	stream := StreamQuotes()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}

// Character returns the character cited in the quote
func (quote *Quote) Character() *Character {
	character, _ := GetCharacter(quote.CharacterId)
	return character
}

// Anime fetches the anime where the quote is said.
func (quote *Quote) Anime() *Anime {
	var anime *Anime
	anime, err := GetAnime(quote.AnimeId)

	if err != nil {
		color.Red("Error fetching anime: %v", err)
	}

	return anime
}

// SortQuotesLatestFirst ...
func SortQuotesLatestFirst(quotes []*Quote) {
	sort.Slice(quotes, func(i, j int) bool {
		return quotes[i].Created > quotes[j].Created
	})
}

// SortQuotesPopularFirst ...
func SortQuotesPopularFirst(quotes []*Quote) {
	sort.Slice(quotes, func(i, j int) bool {
		aLikes := len(quotes[i].Likes)
		bLikes := len(quotes[j].Likes)

		if aLikes == bLikes {
			return quotes[i].Created > quotes[j].Created
		}

		return aLikes > bLikes
	})
}

// FilterQuotes filters all quotes by a custom function.
func FilterQuotes(filter func(*Quote) bool) []*Quote {
	var filtered []*Quote

	for obj := range StreamQuotes() {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered
}
