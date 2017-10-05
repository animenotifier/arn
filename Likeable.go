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
	Save() error
}

// LikeAction ...
func LikeAction(table string) *api.Action {
	return &api.Action{
		Table: table,
		Route: "/like",
		Run: func(obj interface{}, ctx *aero.Context) error {
			likeable := obj.(Likeable)
			user := GetUserFromContext(ctx)

			if user == nil {
				return errors.New("Not logged in")
			}

			likeable.Like(user.ID)
			return likeable.Save()
		},
	}
}

// UnlikeAction ...
func UnlikeAction(table string) *api.Action {
	return &api.Action{
		Table: table,
		Route: "/unlike",
		Run: func(obj interface{}, ctx *aero.Context) error {
			likeable := obj.(Likeable)
			user := GetUserFromContext(ctx)

			if user == nil {
				return errors.New("Not logged in")
			}

			likeable.Unlike(user.ID)
			return likeable.Save()
		},
	}
}
