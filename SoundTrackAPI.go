package arn

import (
	"errors"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ api.Newable   = (*SoundTrack)(nil)
	_ api.Editable  = (*SoundTrack)(nil)
	_ api.Deletable = (*SoundTrack)(nil)
	_ Publishable   = (*SoundTrack)(nil)
)

// Actions
func init() {
	API.RegisterActions("SoundTrack", []*api.Action{
		// Publish
		PublishAction(),

		// Unpublish
		UnpublishAction(),
	})
}

// Create sets the data for a new soundtrack with data we received from the API request.
func (soundtrack *SoundTrack) Create(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	soundtrack.ID = GenerateID("SoundTrack")
	soundtrack.Likes = []string{}
	soundtrack.Created = DateTimeUTC()
	soundtrack.CreatedBy = user.ID
	soundtrack.Media = []*ExternalMedia{}
	soundtrack.Tags = []string{}

	return soundtrack.Unpublish()
}

// Edit updates the external media object.
func (soundtrack *SoundTrack) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	if strings.HasPrefix(key, "Media[") && strings.HasSuffix(key, ".Service") {
		newService := newValue.String()

		if !Contains(ExternalMediaServices, newService) {
			return true, errors.New("Invalid service name")
		}

		value.SetString(newService)
		return true, nil
	}

	return false, nil
}

// AfterEdit updates the metadata.
func (soundtrack *SoundTrack) AfterEdit(ctx *aero.Context) error {
	soundtrack.Edited = DateTimeUTC()
	soundtrack.EditedBy = GetUserFromContext(ctx).ID
	return nil
}

// Delete deletes the object from the database.
func (soundtrack *SoundTrack) Delete() error {
	if soundtrack.IsDraft {
		draftIndex := soundtrack.CreatedByUser().DraftIndex()
		draftIndex.SoundTrackID = ""
		err := draftIndex.Save()

		if err != nil {
			return err
		}
	}

	_, err := DB.Delete("SoundTrack", soundtrack.ID)
	return err
}

// Authorize returns an error if the given API POST request is not authorized.
func (soundtrack *SoundTrack) Authorize(ctx *aero.Context, action string) error {
	if !ctx.HasSession() {
		return errors.New("Neither logged in nor in session")
	}

	return nil
}

// Save saves the soundtrack object in the database.
func (soundtrack *SoundTrack) Save() error {
	return DB.Set("SoundTrack", soundtrack.ID, soundtrack)
}
