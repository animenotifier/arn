package arn

import (
	"encoding/json"
	"errors"

	"github.com/aerogo/aero"
)

// Add adds an anime to the list if it hasn't been added yet.
func (list *PushSubscriptions) Add(sub interface{}) error {
	subscription := sub.(*PushSubscription)

	if list.Contains(subscription) {
		return nil
		// return errors.New("PushSubscription " + subscription.ID() + " has already been added")
	}

	subscription.Created = DateTimeUTC()

	list.Items = append(list.Items, subscription)

	return nil
}

// Remove removes the anime ID from the list.
func (list *PushSubscriptions) Remove(id interface{}) bool {
	subscription := id.(*PushSubscription)

	for index, item := range list.Items {
		if item.ID() == subscription.ID() {
			list.Items = append(list.Items[:index], list.Items[index+1:]...)
			return true
		}
	}

	return false
}

// Contains checks if the list contains the anime ID already.
func (list *PushSubscriptions) Contains(id interface{}) bool {
	subscription := id.(*PushSubscription)

	for _, item := range list.Items {
		if item.ID() == subscription.ID() {
			return true
		}
	}

	return false
}

// Get ...
func (list *PushSubscriptions) Get(sub interface{}) (interface{}, error) {
	item := list.Find(sub.(*PushSubscription).ID())

	if item == nil {
		return nil, errors.New("Not found")
	}

	return item, nil
}

// Set ...
func (list *PushSubscriptions) Set(id interface{}, value interface{}) error {
	// subscription := id.(*PushSubscription)

	// for index, item := range list.Items {
	// 	if item.ID() == subscription.ID() {
	// 		item, ok := value.(*PushSubscription)

	// 		if !ok {
	// 			return errors.New("Missing push subscription properties")
	// 		}

	// 		if item.ID() != animeID {
	// 			return errors.New("Incorrect animeId property")
	// 		}

	// 		item.Edited = DateTimeUTC()
	// 		list.Items[index] = item

	// 		return nil
	// 	}
	// }

	return errors.New("Not implemented")
}

// Update ...
func (list *PushSubscriptions) Update(id interface{}, updatesObj interface{}) error {
	// updates := updatesObj.(map[string]interface{})
	// subscription := id.(*PushSubscription)

	// for _, item := range list.Items {
	// 	if item.PushSubscriptionID == animeID {
	// 		err := SetObjectProperties(item, updates, nil)
	// 		item.Edited = DateTimeUTC()

	// 		return err
	// 	}
	// }

	return errors.New("Not implemented")
}

// Authorize returns an error if the given API request is not authorized.
func (list *PushSubscriptions) Authorize(ctx *aero.Context) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}

// PostBody returns an item that is passed to methods like Add, Remove, etc.
func (list *PushSubscriptions) PostBody(body []byte) interface{} {
	var sub *PushSubscription
	PanicOnError(json.Unmarshal(body, &sub))
	return sub
}

// Save saves the push subscriptions in the database.
func (list *PushSubscriptions) Save() error {
	return DB.Set("PushSubscriptions", list.UserID, list)
}
