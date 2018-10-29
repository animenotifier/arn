package arn

// HasCreator includes user ID and date for the creation of this object.
type HasCreator struct {
	Created   string `json:"created"`
	CreatedBy string `json:"createdBy"`
}

// Creator returns the user who created this object.
func (obj *HasCreator) Creator() *User {
	user, _ := GetUser(obj.CreatedBy)
	return user
}

// CreatorID returns the ID of the user who created this object.
func (obj *HasCreator) CreatorID() string {
	return obj.CreatedBy
}

// GetCreatedBy returns the ID of the user who created this object.
func (obj *HasCreator) GetCreatedBy() string {
	return obj.CreatedBy
}
