package arn

import "sort"

// GetUser ...
func GetUser(id string) (*User, error) {
	obj, err := DB.Get("User", id)
	return obj.(*User), err
}

// GetUserByNick ...
func GetUserByNick(nick string) (*User, error) {
	return GetUserFromTable("NickToUser", nick)
}

// GetUserByEmail ...
func GetUserByEmail(email string) (*User, error) {
	return GetUserFromTable("EmailToUser", email)
}

// GetUserFromTable queries a table for the record with the given ID
// and returns the user that is referenced by record["userId"].
func GetUserFromTable(table string, id string) (*User, error) {
	rec, err := DB.GetMap(table, id)

	if err != nil {
		return nil, err
	}

	return GetUser(rec["userId"].(string))
}

// StreamUsers returns a stream of all users.
func StreamUsers() (chan *User, error) {
	channel := make(chan *User)
	err := DB.Scan("User", channel)
	return channel, err
}

// AllUsers returns a slice of all users.
func AllUsers() ([]*User, error) {
	var all []*User

	stream, err := StreamUsers()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}

// FilterUsers filters all users by a custom function.
func FilterUsers(filter func(*User) bool) ([]*User, error) {
	var filtered []*User

	channel, err := StreamUsers()

	if err != nil {
		return filtered, err
	}

	for obj := range channel {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered, nil
}

// SortUsersLastSeen sorts a list of users by their last seen date.
func SortUsersLastSeen(users []*User) []*User {
	sort.Slice(users, func(i, j int) bool {
		return users[i].LastSeen > users[j].LastSeen
	})

	return users
}
