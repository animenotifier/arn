package arn

import (
	"errors"
	"fmt"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Actions
func init() {
	API.RegisterActions([]*api.Action{
		// Use item
		&api.Action{
			Table: "Inventory",
			Route: "/use/:slot",
			Run: func(obj interface{}, ctx *aero.Context) error {
				inventory := obj.(*Inventory)
				slotIndex, err := ctx.GetInt("slot")

				if err != nil {
					return err
				}

				slot := inventory.Slots[slotIndex]

				if slot.IsEmpty() {
					return errors.New("No item in this slot")
				}

				user := GetUserFromContext(ctx)

				fmt.Println("Use item", slot.Item().Name, "on", user.Nick)

				return inventory.Save()
			},
		},
	})
}

// Authorize returns an error if the given API request is not authorized.
func (inventory *Inventory) Authorize(ctx *aero.Context) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// Save saves the push items in the database.
func (inventory *Inventory) Save() error {
	return DB.Set("Inventory", inventory.UserID, inventory)
}
