package arn

import (
	"github.com/aerogo/nano"
)

// Group represents a group of users.
type Group struct {
	Name        string         `json:"name" editable:"true"`
	Tagline     string         `json:"tagline" editable:"true"`
	Image       string         `json:"image" editable:"true"`
	Description string         `json:"description" editable:"true" type:"textarea"`
	Rules       string         `json:"rules" editable:"true" type:"textarea"`
	Tags        []string       `json:"tags" editable:"true"`
	Members     []*GroupMember `json:"members"`
	Neighbors   []string       `json:"neighbors"`

	// Mixins
	HasID
	HasPosts
	HasCreator
	HasEditor
	HasDraft
}

// Link returns the URI to the group page.
func (group *Group) Link() string {
	return "/group/" + group.ID
}

// TitleByUser returns the preferred title for the given user.
func (group *Group) TitleByUser(user *User) string {
	if group.Name == "" {
		return "untitled"
	}

	return group.Name
}

// String is the default text representation of the group.
func (group *Group) String() string {
	return group.TitleByUser(nil)
}

// ImageURL ...
func (group *Group) ImageURL() string {
	if group.Image != "" {
		return group.Image
	}

	return "https://media.kitsu.io/groups/avatars/2138/medium.png"
	// return "/images/brand/144.png"
}

// FindMember returns the group member by user ID, if available.
func (group *Group) FindMember(userID string) *GroupMember {
	for _, member := range group.Members {
		if member.UserID == userID {
			return member
		}
	}

	return nil
}

// TypeName returns the type name.
func (group *Group) TypeName() string {
	return "Group"
}

// Publish ...
func (group *Group) Publish() error {
	return publish(group)
}

// Unpublish ...
func (group *Group) Unpublish() error {
	return unpublish(group)
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
func StreamGroups() chan *Group {
	channel := make(chan *Group, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Group") {
			channel <- obj.(*Group)
		}

		close(channel)
	}()

	return channel
}

// AllGroups returns a slice of all groups.
func AllGroups() ([]*Group, error) {
	var all []*Group

	stream := StreamGroups()

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}

// FilterGroups filters all groups by a custom function.
func FilterGroups(filter func(*Group) bool) ([]*Group, error) {
	var filtered []*Group

	for obj := range DB.All("Group") {
		realObject := obj.(*Group)

		if filter(realObject) {
			filtered = append(filtered, realObject)
		}
	}

	return filtered, nil
}
