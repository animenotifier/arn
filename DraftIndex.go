package arn

// DraftIndex has references to unpublished drafts a user created.
type DraftIndex struct {
	UserID       string `json:"userId"`
	SoundTrackID string `json:"soundTrackId"`
}

// NewDraftIndex ...
func NewDraftIndex(userID string) *DraftIndex {
	return &DraftIndex{
		UserID: userID,
	}
}
