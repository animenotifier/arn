package arn

import (
	"errors"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ Likeable          = (*Quote)(nil)
	_ LikeEventReceiver = (*Quote)(nil)
	_ Publishable       = (*Quote)(nil)
	_ api.Newable       = (*Quote)(nil)
	_ api.Editable      = (*Quote)(nil)
	_ api.Deletable     = (*Quote)(nil)
)

// Actions
func init() {
	API.RegisterActions("Quote", []*api.Action{
		// Publish
		PublishAction(),

		// Unpublish
		UnpublishAction(),

		// Like
		LikeAction(),

		// Unlike
		UnlikeAction(),
	})
}

// Create sets the data for a new quote with data we received from the API request.
func (quote *Quote) Create(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	quote.ID = GenerateID("Quote")
	quote.Created = DateTimeUTC()
	quote.CreatedBy = user.ID
	quote.Likes = []string{}
	quote.EpisodeNumber = -1
	quote.Time = -1

	return quote.Unpublish()
}

// AfterEdit updates the metadata.
func (quote *Quote) AfterEdit(ctx *aero.Context) error {
	quote.Edited = DateTimeUTC()
	quote.EditedBy = GetUserFromContext(ctx).ID
	return nil
}

// Save saves the quote in the database.
func (quote *Quote) Save() {
	DB.Set("Quote", quote.ID, quote)
}

// Delete deletes the object from the database.
func (quote *Quote) Delete() error {
	if quote.IsDraft {
		draftIndex := quote.Creator().DraftIndex()
		draftIndex.QuoteID = ""
		draftIndex.Save()
	}

	DB.Delete("Quote", quote.ID)
	return nil
}

// Authorize returns an error if the given API request is not authorized.
func (quote *Quote) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	if action == "delete" {
		if user.Role != "editor" && user.Role != "admin" {
			return errors.New("Insufficient permissions")
		}
	}

	return nil
}
