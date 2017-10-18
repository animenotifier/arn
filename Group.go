package arn

import (
	"errors"

	"github.com/fatih/color"
)

// Group ...
type Group struct {
	ID          GroupID        `json:"id"`
	Name        string         `json:"name" editable:"true"`
	Tagline     string         `json:"tagline" editable:"true"`
	Image       string         `json:"image" editable:"true"`
	Description string         `json:"description" editable:"true" type:"textarea"`
	Rules       string         `json:"rules" editable:"true" type:"textarea"`
	Tags        []string       `json:"tags" editable:"true"`
	Members     []*GroupMember `json:"members"`
	Neighbors   []GroupID      `json:"neighbors"`
	IsDraft     bool           `json:"isDraft" editable:"true"`
	Created     UTCDate        `json:"created"`
	CreatedBy   UserID         `json:"createdBy"`
	Edited      UTCDate        `json:"edited"`
	EditedBy    UserID         `json:"editedBy"`

	creator *User
}

// Link ...
func (group *Group) Link() string {
	return "/group/" + group.ID
}

// ImageURL ...
func (group *Group) ImageURL() string {
	if group.Image != "" {
		return group.Image
	}

	return "https://media.kitsu.io/groups/avatars/2138/medium.png"
	// return "/images/brand/144.png"
}

// Creator ...
func (group *Group) Creator() *User {
	if group.creator != nil {
		return group.creator
	}

	user, err := GetUser(group.CreatedBy)

	if err != nil {
		color.Red("Error fetching user: %v", err)
		return nil
	}

	group.creator = user
	return group.creator
}

// Publish ...
func (group *Group) Publish() error {
	if !group.IsDraft {
		return errors.New("Not a draft")
	}

	group.IsDraft = false
	draftIndex, err := GetDraftIndex(group.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.GroupID == "" {
		return errors.New("Group draft doesn't exist in the user draft index")
	}

	draftIndex.GroupID = ""
	return draftIndex.Save()
}

// Unpublish ...
func (group *Group) Unpublish() error {
	group.IsDraft = true
	draftIndex, err := GetDraftIndex(group.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.GroupID != "" {
		return errors.New("You still have an unfinished draft")
	}

	draftIndex.GroupID = group.ID

	return draftIndex.Save()
}

// GetGroup ...
func GetGroup(id string) (*Group, error) {
	obj, err := DB.Get("Group", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Group), nil
}

// StreamGroups returns a stream of all groups.
func StreamGroups() (chan *Group, error) {
	objects, err := DB.All("Group")
	return objects.(chan *Group), err
}

// MustStreamGroups returns a stream of all groups.
func MustStreamGroups() chan *Group {
	stream, err := StreamGroups()
	PanicOnError(err)
	return stream
}

// AllGroups returns a slice of all groups.
func AllGroups() ([]*Group, error) {
	var all []*Group

	stream, err := StreamGroups()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}

// FilterGroups filters all groups by a custom function.
func FilterGroups(filter func(*Group) bool) ([]*Group, error) {
	var filtered []*Group

	channel := make(chan *Group)
	err := DB.Scan("Group", channel)

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
