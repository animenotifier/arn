package arn

import (
	"errors"

	"github.com/aerogo/nano"
)

// Group ...
type Group struct {
	ID          string         `json:"id"`
	Name        string         `json:"name" editable:"true"`
	Tagline     string         `json:"tagline" editable:"true"`
	Image       string         `json:"image" editable:"true"`
	Description string         `json:"description" editable:"true" type:"textarea"`
	Rules       string         `json:"rules" editable:"true" type:"textarea"`
	Tags        []string       `json:"tags" editable:"true"`
	Members     []*GroupMember `json:"members"`
	Neighbors   []string       `json:"neighbors"`
	IsDraft     bool           `json:"isDraft" editable:"true"`
	Created     string         `json:"created"`
	CreatedBy   string         `json:"createdBy"`
	Edited      string         `json:"edited"`
	EditedBy    string         `json:"editedBy"`

	posts []*GroupPost
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

// Posts ...
func (group *Group) Posts() []*GroupPost {
	if group.posts == nil {
		group.posts, _ = FilterGroupPosts(func(post *GroupPost) bool {
			return post.GroupID == group.ID
		})
	}

	return group.posts
}

// Creator ...
func (group *Group) Creator() *User {
	creator, _ := GetUser(group.CreatedBy)
	return creator
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
	draftIndex.Save()

	return nil
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
	draftIndex.Save()

	return nil
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
