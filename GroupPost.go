package arn

import (
	"sort"

	"github.com/aerogo/markdown"
	"github.com/aerogo/nano"
)

// GroupPost represents a group post.
type GroupPost struct {
	ID       string   `json:"id"`
	Text     string   `json:"text" editable:"true"`
	AuthorID string   `json:"authorId"`
	GroupID  string   `json:"groupId"`
	Tags     []string `json:"tags"`
	Likes    []string `json:"likes"`
	IsDraft  bool     `json:"isDraft" editable:"true"`
	Created  string   `json:"created"`
	Edited   string   `json:"edited"`
}

// Author returns the group post's author.
func (post *GroupPost) Author() *User {
	author, _ := GetUser(post.AuthorID)
	return author
}

// Group returns the group post's group.
func (post *GroupPost) Group() *Group {
	group, _ := GetGroup(post.GroupID)
	return group
}

// Link returns the relative URL of the group post.
func (post *GroupPost) Link() string {
	return "/grouppost/" + post.ID
}

// HTML returns the HTML representation of the group post.
func (post *GroupPost) HTML() string {
	return markdown.Render(post.Text)
}

// GetGroupPost ...
func GetGroupPost(id string) (*GroupPost, error) {
	obj, err := DB.Get("GroupPost", id)

	if err != nil {
		return nil, err
	}

	return obj.(*GroupPost), nil
}

// StreamGroupPosts returns a stream of all posts.
func StreamGroupPosts() chan *GroupPost {
	channel := make(chan *GroupPost, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("GroupPost") {
			channel <- obj.(*GroupPost)
		}

		close(channel)
	}()

	return channel
}

// AllGroupPosts returns a slice of all posts.
func AllGroupPosts() ([]*GroupPost, error) {
	var all []*GroupPost

	stream := StreamGroupPosts()

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}

// SortGroupPostsLatestFirst sorts the slice of posts.
func SortGroupPostsLatestFirst(posts []*GroupPost) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Created > posts[j].Created
	})
}

// SortGroupPostsLatestLast sorts the slice of posts.
func SortGroupPostsLatestLast(posts []*GroupPost) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Created < posts[j].Created
	})
}

// GetGroupPostsByUser ...
func GetGroupPostsByUser(user *User) ([]*GroupPost, error) {
	var posts []*GroupPost

	for post := range StreamGroupPosts() {
		if post.AuthorID == user.ID {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

// FilterGroupPosts filters all group posts by a custom function.
func FilterGroupPosts(filter func(*GroupPost) bool) ([]*GroupPost, error) {
	var filtered []*GroupPost

	for post := range StreamGroupPosts() {
		if filter(post) {
			filtered = append(filtered, post)
		}
	}

	return filtered, nil
}

// Like ...
func (post *GroupPost) Like(userID string) {
	for _, id := range post.Likes {
		if id == userID {
			return
		}
	}

	post.Likes = append(post.Likes, userID)

	// Notify author of the post
	go func() {
		likedBy, err := GetUser(userID)

		if err != nil {
			return
		}

		post.Author().SendNotification(&PushNotification{
			Title:   likedBy.Nick + " liked your post",
			Message: likedBy.Nick + " liked your post in the group \"" + post.Group().Name + "\"",
			Icon:    "https:" + likedBy.LargeAvatar(),
			Link:    "https://notify.moe" + likedBy.Link(),
			Type:    NotificationTypeLike,
		})
	}()
}

// Unlike ...
func (post *GroupPost) Unlike(userID string) {
	for index, id := range post.Likes {
		if id == userID {
			post.Likes = append(post.Likes[:index], post.Likes[index+1:]...)
			return
		}
	}
}
