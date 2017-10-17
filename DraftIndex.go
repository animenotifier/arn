package arn

// DraftIndex has references to unpublished drafts a user created.
type DraftIndex struct {
	UserID       UserID       `json:"userId"`
	SoundTrackID SoundTrackID `json:"soundTrackId"`
}

// NewDraftIndex ...
func NewDraftIndex(userID string) *DraftIndex {
	return &DraftIndex{
		UserID: userID,
	}
}

// GetDraftIndex ...
func GetDraftIndex(id string) (*DraftIndex, error) {
	obj, err := DB.Get("DraftIndex", id)

	if err != nil {
		return nil, err
	}

	return obj.(*DraftIndex), nil
}
