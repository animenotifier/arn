package arn

import (
	"errors"

	"github.com/aerogo/nano"
)

// Company ...
type Quote struct {
	ID          string   `json:"id"`
	Description string   `json:"description" editable:"true" type:"textarea"`
	Character   string   `json:"characterId" editable:"true"`
	Anime       string   `json:"animeId" editable:"true"`
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

// Publish ...
func (quote *Quote) Publish() error {
	// No draft
	if !quote.IsDraft {
		return errors.New("Not a draft")
	}

	// No title
	if quote.Description == "" {
		return errors.New("No quote")
	}
	draftIndex, err := GetDraftIndex(quote.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.QuoteID == "" {
		return errors.New("Quote draft doesn't exist in the user draft index")
	}

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

	quote.IsDraft = false
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
