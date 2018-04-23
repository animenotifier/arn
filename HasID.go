package arn

// HasID includes an object ID.
type HasID struct {
	ID string `json:"id"`
}

// GetID returns the ID.
func (obj *HasID) GetID() string {
	return obj.ID
}
