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

// Users returns a slice of all the users you are following.
func (list *UserFollows) Users() []*User {
	objects, err := DB.GetMany("User", list.Items)

	if err != nil {
		return nil
	}

	return objects.([]*User)
}

// GetUserFollows ...
func GetUserFollows(id string) (*UserFollows, error) {
	obj, err := DB.Get("UserFollows", id)

	if err != nil {
		return nil, err
	}

	return obj.(*UserFollows), nil
}

// StreamUserFollows returns a stream of all user follows.
func StreamUserFollows() (chan *UserFollows, error) {
	channel := make(chan *UserFollows)
	err := DB.Scan("UserFollows", channel)
	return channel, err
}

// MustStreamUserFollows returns a stream of all user follows.
func MustStreamUserFollows() chan *UserFollows {
	stream, err := StreamUserFollows()
	PanicOnError(err)
	return stream
}

// AllUserFollows returns a slice of all user follows.
func AllUserFollows() ([]*UserFollows, error) {
	var all []*UserFollows

	stream, err := StreamUserFollows()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}
