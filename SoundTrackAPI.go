package arn

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ Publishable       = (*SoundTrack)(nil)
	_ Likeable          = (*SoundTrack)(nil)
	_ LikeEventReceiver = (*SoundTrack)(nil)
	_ api.Newable       = (*SoundTrack)(nil)
	_ api.Editable      = (*SoundTrack)(nil)
	_ api.Deletable     = (*SoundTrack)(nil)
)

// Actions
func init() {
	API.RegisterActions("SoundTrack", []*api.Action{
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

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "create", "SoundTrack", soundtrack.ID, "", "", "")
	logEntry.Save()

	return soundtrack.Unpublish()
}

// Edit updates the external media object.
func (soundtrack *SoundTrack) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "SoundTrack", soundtrack.ID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	// Verify service name
	if strings.HasPrefix(key, "Media[") && strings.HasSuffix(key, ".Service") {
		newService := newValue.String()
		found := false

		for _, option := range DataLists["media-services"] {
			if option.Label == newService {
				found = true
				break
			}
		}

		if !found {
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
		draftIndex := soundtrack.Creator().DraftIndex()
		draftIndex.SoundTrackID = ""
		draftIndex.Save()
	}

	DB.Delete("SoundTrack", soundtrack.ID)
	return nil
}

// Authorize returns an error if the given API POST request is not authorized.
func (soundtrack *SoundTrack) Authorize(ctx *aero.Context, action string) error {
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

// Save saves the soundtrack object in the database.
func (soundtrack *SoundTrack) Save() {
	DB.Set("SoundTrack", soundtrack.ID, soundtrack)
}
