package arn

import "github.com/parnurzeal/gorequest"

// UserAvatar ...
type UserAvatar struct {
	Extension string `json:"extension"`
	Source    string `json:"source"`
}

// RefreshAvatar ...
func (user *User) RefreshAvatar() (gorequest.Response, string, []error) {
	return gorequest.New().Get("http://media.notify.moe:8001/" + user.ID).End()
}
