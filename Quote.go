package arn

import (
	"errors"

	"sort"

	"github.com/aerogo/nano"
	"github.com/fatih/color"
)

// Quote ...
type Quote struct {
	ID            string    `json:"id"`
	Text          QuoteText `json:"text" editable:"true"`
	CharacterID   string    `json:"characterId" editable:"true"`
	AnimeID       string    `json:"animeId" editable:"true"`
	EpisodeNumber int       `json:"episode" editable:"true"`
	Time          int       `json:"time" editable:"true"`
	IsDraft       bool      `json:"isDraft"`

	HasCreator
	HasEditor
	HasLikes
}

// IsMainQuote returns true if the quote is the main quote of the character.
func (quote *Quote) IsMainQuote() bool {
	return quote.CharacterID != "" && quote.Character().MainQuoteID == quote.ID
}

// Link returns a single quote.
func (quote *Quote) Link() string {
	return "/quote/" + quote.ID
}

// Publish checks the quote and publishes it when no errors were found.
func (quote *Quote) Publish() error {
	// No draft
	if !quote.IsDraft {
		return errors.New("Not a draft")
	}

	// No description
	if quote.Text.English == "" {
		return errors.New("A description is required")
	}

	// No character
	if quote.CharacterID == "" {
		return errors.New("A character is required")
	}

	// No anime
	if quote.AnimeID == "" {
		return errors.New("An anime is required")
	}

	// // No episode number
	// if quote.EpisodeNumber == -1 {
	// 	return errors.New("An episode number is required")
	// }

	// // No time
	// if quote.Time == -1 {
	// 	return errors.New("Time in minutes is required")
	// }

	// Invalid anime ID
	anime := quote.Anime()

	if anime == nil {
		return errors.New("Invalid anime ID")
	}

	// Invalid episode number
	maxEpisodes := anime.EpisodeCount

	if maxEpisodes != 0 && quote.EpisodeNumber > maxEpisodes {
		return errors.New("Invalid episode number")
	}

	// Draft index
	draftIndex, err := GetDraftIndex(quote.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.QuoteID == "" {
		return errors.New("Quote draft doesn't exist in the user draft index")
	}

	// Invalid character ID
	_, characterErr := GetCharacter(quote.CharacterID)

	if characterErr != nil {
		return errors.New("Character does not exist")
	}

	// Publish
	quote.IsDraft = false
	draftIndex.QuoteID = ""
	draftIndex.Save()
	return nil
}

// OnLike is called when the quote receives a like.
func (quote *Quote) OnLike(likedBy *User) {
	if likedBy.ID == quote.CreatedBy {
		return
	}

	if !quote.Creator().Settings().Notification.QuoteLikes {
		return
	}

	go func() {
		quote.Creator().SendNotification(&PushNotification{
			Title:   likedBy.Nick + " liked your " + quote.Character().Name.Canonical + " quote",
			Message: quote.Text.English,
			Icon:    "https:" + likedBy.AvatarLink("large"),
			Link:    "https://notify.moe" + likedBy.Link(),
			Type:    NotificationTypeLike,
		})
	}()
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

// String implements the default string serialization.
func (quote *Quote) String() string {
	return quote.Text.English
}

// GetQuote returns a single quote.
func GetQuote(id string) (*Quote, error) {
	obj, err := DB.Get("Quote", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Quote), nil
}

// StreamQuotes returns a stream of all quotes.
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
	character, _ := GetCharacter(quote.CharacterID)
	return character
}

// Anime fetches the anime where the quote is said.
func (quote *Quote) Anime() *Anime {
	anime, err := GetAnime(quote.AnimeID)

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
