package arn

import "github.com/aerogo/nano"

// Notification represents a user-associated notification.
type Notification struct {
	ID      string `json:"id"`
	UserID  string `json:"userId"`
	Created string `json:"created"`
	Seen    string `json:"seen"`
	PushNotification
}

// User retrieves the user the notification was sent to.
func (notification *Notification) User() *User {
	user, _ := GetUser(notification.UserID)
	return user
}

// NewNotification creates a new notification.
func NewNotification(userID string, pushNotification *PushNotification) *Notification {
	return &Notification{
		ID:               GenerateID("Notification"),
		UserID:           userID,
		Created:          DateTimeUTC(),
		Seen:             "",
		PushNotification: *pushNotification,
	}
}

// GetNotification ...
func GetNotification(id string) (*Notification, error) {
	obj, err := DB.Get("Notification", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Notification), nil
}

// StreamNotifications returns a stream of all notifications.
func StreamNotifications() chan *Notification {
	channel := make(chan *Notification, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Notification") {
			channel <- obj.(*Notification)
		}

		close(channel)
	}()

	return channel
}

// AllNotifications returns a slice of all notifications.
func AllNotifications() ([]*Notification, error) {
	var all []*Notification

	for obj := range StreamNotifications() {
		all = append(all, obj)
	}

	return all, nil
}
