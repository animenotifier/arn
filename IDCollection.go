package arn

import (
	"errors"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// IDCollection ...
type IDCollection interface {
	Add(id string) error
	Remove(id string) bool
	Save() error
}

// AddAction returns an API action that adds a new item to the IDCollection.
func AddAction() *api.Action {
	return &api.Action{
		Route: "/add/:item-id",
		Run: func(obj interface{}, ctx *aero.Context) error {
			list := obj.(IDCollection)
			itemID := ctx.Get("item-id")
			err := list.Add(itemID)

			if err != nil {
				return err
			}

			return list.Save()
		},
	}
}

// RemoveAction returns an API action that removes an item from the IDCollection.
func RemoveAction() *api.Action {
	return &api.Action{
		Route: "/remove/:item-id",
		Run: func(obj interface{}, ctx *aero.Context) error {
			list := obj.(IDCollection)
			itemID := ctx.Get("item-id")

			if !list.Remove(itemID) {
				return errors.New("This item does not exist in the list")
			}

			return list.Save()
		},
	}
}
