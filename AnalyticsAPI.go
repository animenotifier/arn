package arn

import (
	"encoding/json"

	"github.com/aerogo/aero"
)

// Authorize returns an error if the given API POST request is not authorized.
func (analytics *Analytics) Authorize(ctx *aero.Context) error {
	return AuthorizeIfLoggedIn(ctx)
}

// Create creates a new analytics object.
func (analytics *Analytics) Create(ctx *aero.Context) error {
	err := json.Unmarshal(ctx.RequestBody(), analytics)

	if err != nil {
		return err
	}

	analytics.UserID = GetUserFromContext(ctx).ID

	return nil
}

// Save saves the analytics in the database.
func (analytics *Analytics) Save() error {
	return DB.Set("Analytics", analytics.UserID, analytics)
}
