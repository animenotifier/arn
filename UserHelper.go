package arn

import (
	"errors"
	"sort"
	"strconv"

	"github.com/aerogo/nano"
	"github.com/animenotifier/osu"
)

// GetUser ...
func GetUser(id string) (*User, error) {
	obj, err := DB.Get("User", id)

	if err != nil {
		return nil, err
	}

	return obj.(*User), nil
}

// GetUserByNick ...
func GetUserByNick(nick string) (*User, error) {
	obj, err := DB.Get("NickToUser", nick)

	if err != nil {
		return nil, err
	}

	userID := obj.(*NickToUser).UserID
	user, err := GetUser(userID)

	return user, err
}

// GetUserByEmail ...
func GetUserByEmail(email string) (*User, error) {
	obj, err := DB.Get("EmailToUser", email)

	if err != nil {
		return nil, err
	}

	userID := obj.(*EmailToUser).UserID
	user, err := GetUser(userID)

	return user, err
}

// GetUserByFacebookID ...
func GetUserByFacebookID(facebookID string) (*User, error) {
	obj, err := DB.Get("FacebookToUser", facebookID)

	if err != nil {
		return nil, err
	}

	userID := obj.(*FacebookToUser).UserID
	user, err := GetUser(userID)

	return user, err
}

// GetUserByGoogleID ...
func GetUserByGoogleID(googleID string) (*User, error) {
	obj, err := DB.Get("GoogleToUser", googleID)

	if err != nil {
		return nil, err
	}

	userID := obj.(*GoogleToUser).UserID
	user, err := GetUser(userID)

	return user, err
}

// StreamUsers returns a stream of all users.
func StreamUsers() chan *User {
	channel := make(chan *User, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("User") {
			channel <- obj.(*User)
		}

		close(channel)
	}()

	return channel
}

// AllUsers returns a slice of all users.
func AllUsers() ([]*User, error) {
	var all []*User

	for obj := range StreamUsers() {
		all = append(all, obj)
	}

	return all, nil
}

// FilterUsers filters all users by a custom function.
func FilterUsers(filter func(*User) bool) []*User {
	var filtered []*User

	for obj := range StreamUsers() {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered
}

// SortUsersLastSeen sorts a list of users by their last seen date.
func SortUsersLastSeen(users []*User) {
	sort.Slice(users, func(i, j int) bool {
		return users[i].LastSeen > users[j].LastSeen
	})
}

// RefreshOsuInfo refreshes a user's Osu information.
func (user *User) RefreshOsuInfo() error {
	if user.Accounts.Osu.Nick == "" {
		return errors.New("User doesn't have an osu username")
	}

	osu, err := osu.GetUser(user.Accounts.Osu.Nick)

	if err != nil {
		return err
	}

	user.Accounts.Osu.PP, _ = strconv.ParseFloat(osu.PPRaw, 64)
	user.Accounts.Osu.Level, _ = strconv.ParseFloat(osu.Level, 64)
	user.Accounts.Osu.Accuracy, _ = strconv.ParseFloat(osu.Accuracy, 64)

	return nil
}
