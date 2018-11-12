package arn

// ActivityConsumeAnime is a user activity that consumes anime.
type ActivityConsumeAnime struct {
	AnimeID     string `json:"animeId"`
	FromEpisode int    `json:"fromEpisode"`
	ToEpisode   int    `json:"toEpisode"`

	HasID
	HasCreator
}

// NewActivityConsumeAnime creates a new activity.
func NewActivityConsumeAnime(objectType string, objectID string, userID string) *ActivityConsumeAnime {
	return &ActivityConsumeAnime{
		HasID: HasID{
			ID: GenerateID("ActivityConsumeAnime"),
		},
		HasCreator: HasCreator{
			Created:   DateTimeUTC(),
			CreatedBy: userID,
		},
	}
}

// Type returns the type name.
func (activity *ActivityConsumeAnime) Type() string {
	return "ActivityConsumeAnime"
}

// // OnLike is called when the activity receives a like.
// func (activity *Activity) OnLike(likedBy *User) {
// 	if likedBy.ID == activity.CreatedBy {
// 		return
// 	}

// 	go func() {
// 		notifyUser := activity.Creator()

// 		notifyUser.SendNotification(&PushNotification{
// 			Title:   likedBy.Nick + " liked your activity",
// 			Message: activity.TextByUser(notifyUser),
// 			Icon:    "https:" + likedBy.AvatarLink("large"),
// 			Link:    "https://notify.moe" + activity.Link(),
// 			Type:    NotificationTypeLike,
// 		})
// 	}()
// }
