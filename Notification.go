package arn

// Notification ...
type Notification struct {
	ID      string `json:"id"`
	UserID  string `json:"userId"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Icon    string `json:"icon"`
	Image   string `json:"image"`
	Link    string `json:"link"`
	Created string `json:"created"`
	Seen    string `json:"seen"`
}

// CreateNotification creates a new notification.
func CreateNotification(userID string) *Notification {
	return &Notification{
		ID:      GenerateID("Notification"),
		UserID:  userID,
		Created: DateTimeUTC(),
	}
}
