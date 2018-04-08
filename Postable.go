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
	Thread() *Thread
	ThreadID() string
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

// FilterPostablesWithUniqueThreads removes posts with the same thread until we have enough posts.
func FilterPostablesWithUniqueThreads(posts []Postable, limit int) []Postable {
	filtered := []Postable{}
	threadsProcessed := map[string]bool{}

	for _, post := range posts {
		if len(filtered) >= limit {
			return filtered
		}

		_, found := threadsProcessed[post.ThreadID()]

		if found {
			continue
		}

		threadsProcessed[post.ThreadID()] = true
		filtered = append(filtered, post)
	}

	return filtered
}
