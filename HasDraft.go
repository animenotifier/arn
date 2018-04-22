package arn

// HasDraft includes a boolean indicating whether the object is a draft.
type HasDraft struct {
	IsDraft bool `json:"isDraft" editable:"true"`
}
