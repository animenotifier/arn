package arn

import (
	"errors"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Actions
func init() {
	API.RegisterActions([]*api.Action{
		// Add follow
		&api.Action{
			Table: "UserFollows",
			Route: "/add/:userId",
			Run: func(obj interface{}, ctx *aero.Context) error {
				userFollows := obj.(*UserFollows)
				userID := ctx.Get("userId")
				err := userFollows.Add(userID)

				if err != nil {
					return err
				}

				return userFollows.Save()
			},
		},

		// Remove follow
		&api.Action{
			Table: "UserFollows",
			Route: "/remove/:userId",
			Run: func(obj interface{}, ctx *aero.Context) error {
				userFollows := obj.(*UserFollows)
				userID := ctx.Get("userId")

				if !userFollows.Remove(userID) {
					return errors.New("You are not following this user")
				}

				return userFollows.Save()
			},
		},
	})
}

// Add adds an user to the list if it hasn't been added yet.
func (list *UserFollows) Add(userID string) error {
	if userID == list.UserID {
		return errors.New("You can't follow yourself")
	}

	if list.Contains(userID) {
		return errors.New("User " + userID + " has already been added")
	}

	list.Items = append(list.Items, userID)

	// Send notification
	user, err := GetUser(userID)

	if err == nil {
		follower, err := GetUser(list.UserID)

		if err == nil {
			user.SendNotification(&Notification{
				Title:   "You have a new follower!",
				Message: follower.Nick + " started following you.",
				Icon:    "https:" + follower.LargeAvatar(),
				Link:    "https://notify.moe" + follower.Link(),
			})
		}
	}

	return nil
}

// Remove removes the user ID from the list.
func (list *UserFollows) Remove(userID string) bool {
	for index, item := range list.Items {
		if item == userID {
			list.Items = append(list.Items[:index], list.Items[index+1:]...)
			return true
		}
	}

	return false
}

// Contains checks if the list contains the user ID already.
func (list *UserFollows) Contains(userID string) bool {
	for _, item := range list.Items {
		if item == userID {
			return true
		}
	}

	return false
}

// Authorize returns an error if the given API request is not authorized.
func (list *UserFollows) Authorize(ctx *aero.Context, action string) error {
	return AuthorizeIfLoggedInAndOwnData(ctx, "id")
}
