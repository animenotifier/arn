package arn

// Notification represents a user-associated notification.
type Notification struct {
	ID      string `json:"id"`
	UserID  string `json:"userId"`
	Created string `json:"created"`
	Seen    string `json:"seen"`
	PushNotification
}

// CreateNotification creates a new notification.
func CreateNotification(userID string, pushNotification *PushNotification) *Notification {
	return &Notification{
		ID:               GenerateID("Notification"),
		UserID:           userID,
		Created:          DateTimeUTC(),
		Seen:             "",
		PushNotification: *pushNotification,
	}
}
