package arn

import (
	"errors"

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
	API.RegisterActions("Group", []*api.Action{
		// Publish
		PublishAction(),

		// Unpublish
		UnpublishAction(),
	})
}

// Create ...
func (group *Group) Create(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	group.ID = GenerateID("Group")
	group.Tags = []string{}
	group.Neighbors = []string{}
	group.Created = DateTimeUTC()
	group.CreatedBy = user.ID
	group.Edited = group.Created
	group.EditedBy = group.CreatedBy

	group.Members = []*GroupMember{
		&GroupMember{
			UserID: user.ID,
			Role:   "founder",
			Joined: group.Created,
		},
	}

	return group.Unpublish()
}

// AfterEdit updates the metadata.
func (group *Group) AfterEdit(ctx *aero.Context) error {
	group.Edited = DateTimeUTC()
	group.EditedBy = GetUserFromContext(ctx).ID
	return nil
}

// Delete deletes the object from the database.
func (group *Group) Delete() error {
	if group.IsDraft {
		draftIndex := group.Creator().DraftIndex()
		draftIndex.GroupID = ""
		err := draftIndex.Save()

		if err != nil {
			return err
		}
	}

	_, err := DB.Delete("Group", group.ID)
	return err
}

// Authorize returns an error if the given API POST request is not authorized.
func (group *Group) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	if action == "edit" && group.CreatedBy != user.ID {
		return errors.New("Can't edit groups from other people")
	}

	return nil
}

// Save saves the group in the database.
func (group *Group) Save() error {
	return DB.Set("Group", group.ID, group)
}
