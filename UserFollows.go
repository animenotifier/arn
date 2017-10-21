package arn

import "errors"

// UserFollows ...
type UserFollows struct {
	UserID string       `json:"userId"`
	Items  []string `json:"items"`
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

// Users returns a slice of all the users you are following.
func (list *UserFollows) Users() []*User {
	objects, err := DB.GetMany("User", list.Items)

	if err != nil {
		return nil
	}

	return objects.([]*User)
}

// GetUserFollows ...
func GetUserFollows(id string) (*UserFollows, error) {
	obj, err := DB.Get("UserFollows", id)

	if err != nil {
		return nil, err
	}

	return obj.(*UserFollows), nil
}

// StreamUserFollows returns a stream of all user follows.
func StreamUserFollows() (chan *UserFollows, error) {
	channel := make(chan *UserFollows)
	err := DB.Scan("UserFollows", channel)
	return channel, err
}

// MustStreamUserFollows returns a stream of all user follows.
func MustStreamUserFollows() chan *UserFollows {
	stream, err := StreamUserFollows()
	PanicOnError(err)
	return stream
}

// AllUserFollows returns a slice of all user follows.
func AllUserFollows() ([]*UserFollows, error) {
	var all []*UserFollows

	stream, err := StreamUserFollows()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}
