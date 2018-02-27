package arn

import (
	"errors"

	"github.com/aerogo/nano"
)

// UserNotifications ...
type UserNotifications struct {
	UserID string   `json:"userId"`
	Items  []string `json:"items"`
}

// Add adds an user to the list if it hasn't been added yet.
func (list *UserNotifications) Add(notificationID string) error {
	if list.Contains(notificationID) {
		return errors.New("Notification " + notificationID + " has already been added")
	}

	list.Items = append(list.Items, notificationID)
	return nil
}

// Remove removes the notification ID from the list.
func (list *UserNotifications) Remove(notificationID string) bool {
	for index, item := range list.Items {
		if item == notificationID {
			list.Items = append(list.Items[:index], list.Items[index+1:]...)
			return true
		}
	}

	return false
}

// Contains checks if the list contains the notification ID already.
func (list *UserNotifications) Contains(notificationID string) bool {
	for _, item := range list.Items {
		if item == notificationID {
			return true
		}
	}

	return false
}

// Notifications returns a slice of all the notifications.
func (list *UserNotifications) Notifications() []*Notification {
	notificationsObj := DB.GetMany("Notification", list.Items)
	notifications := make([]*Notification, len(notificationsObj), len(notificationsObj))

	for i, obj := range notificationsObj {
		notifications[i] = obj.(*Notification)
	}

	return notifications
}

// GetUserNotifications ...
func GetUserNotifications(id string) (*UserNotifications, error) {
	obj, err := DB.Get("UserNotifications", id)

	if err != nil {
		return nil, err
	}

	return obj.(*UserNotifications), nil
}

// StreamUserNotifications returns a stream of all user notifications.
func StreamUserNotifications() chan *UserNotifications {
	channel := make(chan *UserNotifications, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("UserNotifications") {
			channel <- obj.(*UserNotifications)
		}

		close(channel)
	}()

	return channel
}

// AllUserNotifications returns a slice of all user notifications.
func AllUserNotifications() ([]*UserNotifications, error) {
	var all []*UserNotifications

	for obj := range StreamUserNotifications() {
		all = append(all, obj)
	}

	return all, nil
}
