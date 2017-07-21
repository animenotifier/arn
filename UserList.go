package arn

// UserFollows ...
type UserFollows struct {
	UserID string   `json:"userId"`
	Items  []string `json:"items"`
}

// Save saves the episodes in the database.
func (list *UserFollows) Save() error {
	return DB.Set("UserFollows", list.UserID, list)
}

// GetUserFollows ...
func GetUserFollows(id string) (*UserFollows, error) {
	obj, err := DB.Get("UserFollows", id)

	if err != nil {
		return nil, err
	}

	return obj.(*UserFollows), nil
}
