package arn

import (
	"errors"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ Likeable          = (*Activity)(nil)
	_ LikeEventReceiver = (*Activity)(nil)
	_ api.Deletable     = (*Activity)(nil)
)

// Actions
func init() {
	API.RegisterActions("Activity", []*api.Action{
		// Like
		LikeAction(),

		// Unlike
		UnlikeAction(),
	})
}

// Authorize returns an error if the given API request is not authorized.
func (activity *Activity) Authorize(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	return nil
}

// DeleteInContext deletes the activity in the given context.
func (activity *Activity) DeleteInContext(ctx *aero.Context) error {
	return activity.Delete()
}

// Delete deletes the object from the database.
func (activity *Activity) Delete() error {
	DB.Delete("Activity", activity.ID)
	return nil
}

// Save saves the activity object in the database.
func (activity *Activity) Save() {
	DB.Set("Activity", activity.ID, activity)
}
