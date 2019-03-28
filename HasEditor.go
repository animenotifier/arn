package arn

import (
	"github.com/aerogo/aero"
)

// HasEditor includes user ID and date for the last edit of this object.
type HasEditor struct {
	Edited   string `json:"edited"`
	EditedBy string `json:"editedBy"`
}

// Editor returns the user who last edited this object.
func (obj *HasEditor) Editor() *User {
	user, _ := GetUser(obj.EditedBy)
	return user
}

// AfterEdit updates the metadata.
func (obj *HasEditor) AfterEdit(ctx *aero.Context) error {
	obj.Edited = DateTimeUTC()
	obj.EditedBy = GetUserFromContext(ctx).ID
	return nil
}
