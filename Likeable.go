package arn

import (
	"errors"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Likeable ...
type Likeable interface {
	Like(userID string)
	Unlike(userID string)
	LikedBy(userID string) bool
	Link() string
	Save()
}

// LikeAction ...
func LikeAction() *api.Action {
	return &api.Action{
		Route: "/like",
		Run: func(obj interface{}, ctx *aero.Context) error {
			likeable := obj.(Likeable)
			user := GetUserFromContext(ctx)

			if user == nil {
				return errors.New("Not logged in")
			}

			likeable.Like(user.ID)
			likeable.Save()
			return nil
		},
	}
}

// UnlikeAction ...
func UnlikeAction() *api.Action {
	return &api.Action{
		Route: "/unlike",
		Run: func(obj interface{}, ctx *aero.Context) error {
			likeable := obj.(Likeable)
			user := GetUserFromContext(ctx)

			if user == nil {
				return errors.New("Not logged in")
			}

			likeable.Unlike(user.ID)
			likeable.Save()
			return nil
		},
	}
}
