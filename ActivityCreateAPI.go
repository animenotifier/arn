package arn

import (
	"github.com/aerogo/api"
)

// Force interface implementations
var (
	_ Activity    = (*ActivityCreate)(nil)
	_ api.Savable = (*ActivityCreate)(nil)
)

// Save saves the activity object in the database.
func (activity *ActivityCreate) Save() {
	DB.Set("ActivityCreate", activity.ID, activity)
}
