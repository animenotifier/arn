package arn

import (
	"errors"

	"github.com/aerogo/aero"
)

// Add adds an user to the list if it hasn't been added yet.
func (list *UserList) Add(id interface{}) error {
	userID := id.(string)

	if list.Contains(userID) {
		return errors.New("User " + userID + " has already been added")
	}

	list.Items = append(list.Items, userID)

	return nil
}

// Remove removes the user ID from the list.
func (list *UserList) Remove(id interface{}) bool {
	userID := id.(string)

	for index, item := range list.Items {
		if item == userID {
			list.Items = append(list.Items[:index], list.Items[index+1:]...)
			return true
		}
	}

	return false
}

// Contains checks if the list contains the user ID already.
func (list *UserList) Contains(id interface{}) bool {
	userID := id.(string)

	for _, item := range list.Items {
		if item == userID {
			return true
		}
	}

	return false
}

// Get ...
func (list *UserList) Get(id interface{}) (interface{}, error) {
	userID := id.(string)

	for _, item := range list.Items {
		if item == userID {
			return item, nil
		}
	}

	return nil, errors.New("Not found")
}

// Set ...
func (list *UserList) Set(id interface{}, value interface{}) error {
	return errors.New("Not applicable")
}

// Update ...
func (list *UserList) Update(id interface{}, updatesObj interface{}) error {
	return errors.New("Not applicable")
}

// Authorize returns an error if the given API request is not authorized.
func (list *UserList) Authorize(ctx *aero.Context) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// PostBody returns an item that is passed to methods like Add, Remove, etc.
func (list *UserList) PostBody(body []byte) interface{} {
	return string(body)
}
