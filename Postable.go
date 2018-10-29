package arn

import (
	"reflect"
	"sort"
)

// Postable is a generic interface for Threads, Posts and Messages.
type Postable interface {
	ID() string
	Title() string
	Text() string
	HTML() string
	Likes() []string
	LikedBy(userID string) bool
	Parent() PostParent
	ParentID() string
	Link() string
	Type() string
	Creator() *User
	Created() string
}

// CanBePostable is a type that defines the ToPostable() conversion.
type CanBePostable interface {
	ToPostable() Postable
}

// ToPostable converts a specific type to a generic postable.
func ToPostable(post CanBePostable) Postable {
	return post.ToPostable()
}

// ToPostables converts a slice of specific types to a slice of generic postables.
func ToPostables(sliceOfPosts interface{}) []Postable {
	var postables []Postable

	v := reflect.ValueOf(sliceOfPosts)

	for i := 0; i < v.Len(); i++ {
		canBePostable := v.Index(i).Interface().(CanBePostable)
		postables = append(postables, canBePostable.ToPostable())
	}

	return postables
}

// SortPostablesLatestFirst ...
func SortPostablesLatestFirst(posts []Postable) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Created() > posts[j].Created()
	})
}
