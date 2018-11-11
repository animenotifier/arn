package arn

const (
	// ActivityTypeCreate when new objects are created.
	ActivityTypeCreate = "create"

	// type ActivityCreate struct {
	// 	ObjectID   string `json:"objectId"`
	// 	ObjectType string `json:"objectType"`
	// }

	// ActivityTypeConsume when media is consumed (anime episodes watched).
	ActivityTypeConsume = "consume"

	// type ActivityConsumeAnime struct {
	// 	AnimeID     string `json:"animeId"`
	// 	FromEpisode int    `json:"fromEpisode"`
	// 	ToEpisode   int    `json:"toEpisode"`
	// }

	// ActivityTypeComplete when media is completed.
	ActivityTypeComplete = "complete"

	// ActivityTypeDrop when media is dropped.
	ActivityTypeDrop = "drop"
)

// Activity is a user activity that appears in the follower's feeds.
type Activity struct {
	Type string                 `json:"type"`
	Meta map[string]interface{} `json:"meta"`

	HasID
	HasCreator
	HasLikes
}

// NewActivity creates a new activity.
func NewActivity(typ string, meta map[string]interface{}, userID string) *Activity {
	return &Activity{
		HasID: HasID{
			ID: GenerateID("Activity"),
		},
		HasCreator: HasCreator{
			Created:   DateTimeUTC(),
			CreatedBy: userID,
		},
		Type: typ,
		Meta: meta,
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
