package arn

import (
	"errors"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
	"reflect"
)

// Force interface implementations
var (
	_ Likeable      = (*Quote)(nil)
	_ Publishable   = (*Quote)(nil)
	_ api.Newable   = (*Quote)(nil)
	_ api.Editable  = (*Quote)(nil)
	_ api.Deletable = (*Quote)(nil)
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

	return quote.Unpublish()
}

// Edit remove the quote from it's previous linked character quote list if the new one is different
// and add it to the new one.
func (quote *Quote) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	if key == "CharacterId" {
		newCharacterId := newValue.String()
		previousCharacterId := newValue.String()
		if previousCharacterId != newCharacterId {
			newCharacter, err := GetCharacter(newCharacterId)
			previousCharacter, err := GetCharacter(previousCharacterId)

			if err != nil {
				return false, err
			}

			// Remove the reference of the quote in the previous character that contained it
			if !previousCharacter.Remove(quote.ID) {
				return false, errors.New("This quote does not exist")
			}

			previousCharacter.Save()
			value.SetString(newCharacterId)

			// Append to quotes Ids to the new character
			newCharacter.QuotesIds = append(newCharacter.QuotesIds, quote.ID)

			// Save the character
			newCharacter.Save()
			return true, nil
		}

	}
	return false, nil
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
	character, err := GetCharacter(quote.CharacterId)

	if err != nil {
		return err
	}

	// Remove the reference of the quote in the character that contains it
	if !character.Remove(quote.ID) {
		return errors.New("This quote does not exist")
	}

	character.Save()

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
