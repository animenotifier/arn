package arn

import (
	"github.com/aerogo/aero"
)

// Authorize returns an error if the given API request is not authorized.
func (inventory *Inventory) Authorize(ctx *aero.Context) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// Save saves the push items in the database.
func (inventory *Inventory) Save() error {
	return DB.Set("Inventory", inventory.UserID, inventory)
}
