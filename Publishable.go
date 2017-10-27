package arn

import (
	"github.com/aerogo/aero"
	"github.com/aerogo/api"
)

// Publishable ...
type Publishable interface {
	Publish() error
	Unpublish() error
	Save()
}

// PublishAction returns an API action that publishes the object.
func PublishAction() *api.Action {
	return &api.Action{
		Route: "/publish",
		Run: func(obj interface{}, ctx *aero.Context) error {
			draft := obj.(Publishable)
			err := draft.Publish()

			if err != nil {
				return err
			}

			draft.Save()
			return nil
		},
	}
}

// UnpublishAction returns an API action that unpublishes the object.
func UnpublishAction() *api.Action {
	return &api.Action{
		Route: "/unpublish",
		Run: func(obj interface{}, ctx *aero.Context) error {
			draft := obj.(Publishable)
			err := draft.Unpublish()

			if err != nil {
				return err
			}

			draft.Save()
			return nil
		},
	}
}
