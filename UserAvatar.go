package arn

import (
	"github.com/aerogo/http/client"
)

// UserAvatar ...
type UserAvatar struct {
	Extension string `json:"extension"`
	Source    string `json:"source"`
}

// RefreshAvatar ...
func (user *User) RefreshAvatar() (client.Response, error) {
	return client.Get("http://127.0.0.1:8001/" + user.ID).End()
}
