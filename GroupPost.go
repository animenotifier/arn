package arn

import (
	"sort"

	"github.com/aerogo/markdown"
)

// GroupPost represents a group post.
type GroupPost struct {
	ID       string   `json:"id"`
	Text     string   `json:"text" editable:"true"`
	AuthorID UserID   `json:"authorId"`
	GroupID  GroupID  `json:"groupId"`
	Tags     []string `json:"tags"`
	Likes    []string `json:"likes"`
	IsDraft  bool     `json:"isDraft" editable:"true"`
	Created  UTCDate  `json:"created"`
	Edited   UTCDate  `json:"edited"`

	author *User
	group  *Group
	html   string
}

// Author returns the group post's author.
func (post *GroupPost) Author() *User {
	if post.author != nil {
		return post.author
	}

	post.author, _ = GetUser(post.AuthorID)
	return post.author
}

// Group returns the group post's group.
func (post *GroupPost) Group() *Group {
	if post.group != nil {
		return post.group
	}

	post.group, _ = GetGroup(post.GroupID)
	return post.group
}

// Link returns the relative URL of the group post.
func (post *GroupPost) Link() string {
	return "/grouppost/" + post.ID
}

// HTML returns the HTML representation of the group post.
func (post *GroupPost) HTML() string {
	if post.html != "" {
		return post.html
	}

	post.html = markdown.Render(post.Text)
	return post.html
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
func StreamGroupPosts() (chan *GroupPost, error) {
	objects, err := DB.All("GroupPost")
	return objects.(chan *GroupPost), err
}

// AllGroupPosts returns a slice of all posts.
func AllGroupPosts() ([]*GroupPost, error) {
	var all []*GroupPost

	stream, err := StreamGroupPosts()

	if err != nil {
		return nil, err
	}

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

	stream := make(chan *GroupPost)
	err := DB.Scan("GroupPost", stream)

	if err != nil {
		return nil, err
	}

	for post := range stream {
		if post.AuthorID == user.ID {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

// FilterGroupPosts filters all group posts by a custom function.
func FilterGroupPosts(filter func(*GroupPost) bool) ([]*GroupPost, error) {
	var filtered []*GroupPost

	channel := make(chan *GroupPost)
	err := DB.Scan("GroupPost", channel)

	if err != nil {
		return filtered, err
	}

	for post := range channel {
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

		post.Author().SendNotification(&Notification{
			Title:   likedBy.Nick + " liked your post",
			Message: likedBy.Nick + " liked your post in the group \"" + post.Group().Name + "\"",
			Icon:    "https:" + likedBy.LargeAvatar(),
			Link:    "https://notify.moe" + likedBy.Link(),
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
