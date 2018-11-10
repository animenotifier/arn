package arn

const (
	// ActivityTypeCreate when new objects are created.
	ActivityTypeCreate = "create"

	// ActivityTypeConsume when media is consumed (anime episodes watched).
	ActivityTypeConsume = "consume"

	// ActivityTypeComplete when media is completed.
	ActivityTypeComplete = "complete"

	// ActivityTypeDrop when media is dropped.
	ActivityTypeDrop = "drop"
)

// Activity is a user activity that appears in the follower's feeds.
type Activity struct {
	Type       string            `json:"type"`
	ObjectType string            `json:"objectType"`
	ObjectID   string            `json:"objectId"`
	Meta       map[string]string `json:"meta"`

	HasID
	HasCreator
	HasLikes
}

// NewActivity creates a new activity.
func NewActivity(typ string, objectType string, objectID string, meta map[string]string, userID string) *Activity {
	return &Activity{
		HasID: HasID{
			ID: GenerateID("Activity"),
		},
		ObjectType: objectType,
		ObjectID:   objectID,
		HasCreator: HasCreator{
			Created:   DateTimeUTC(),
			CreatedBy: userID,
		},
	}
}

// Link returns the permalink for the Activity.
func (activity *Activity) Link() string {
	return "/activity/" + activity.ID
}

// Text returns the textual representation of the activity.
func (activity *Activity) Text() string {
	return "Watched episode 123"
}

// OnLike is called when the activity receives a like.
func (activity *Activity) OnLike(likedBy *User) {
	if likedBy.ID == activity.CreatedBy {
		return
	}

	go func() {
		activity.Creator().SendNotification(&PushNotification{
			Title:   likedBy.Nick + " liked your activity",
			Message: activity.Text(),
			Icon:    "https:" + likedBy.AvatarLink("large"),
			Link:    "https://notify.moe" + likedBy.Link(),
			Type:    NotificationTypeLike,
		})
	}()
}
