package arn

import "sort"
import "github.com/aerogo/aero"

// Post represents a forum post.
type Post struct {
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	AuthorID string   `json:"authorId"`
	ThreadID string   `json:"threadId"`
	Tags     []string `json:"tags"`
	Likes    []string `json:"likes"`
	Created  string   `json:"created"`
	Edited   string   `json:"edited"`

	author *User
	thread *Thread
	html   string
}

// Author returns the post author.
func (post *Post) Author() *User {
	if post.author != nil {
		return post.author
	}

	post.author, _ = GetUser(post.AuthorID)
	return post.author
}

// Thread returns the thread this post was posted in.
func (post *Post) Thread() *Thread {
	if post.thread != nil {
		return post.thread
	}

	post.thread, _ = GetThread(post.ThreadID)
	return post.thread
}

// Link returns the relative URL of the post.
func (post *Post) Link() string {
	return "/posts/" + post.ID
}

// HTML returns the HTML representation of the post.
func (post *Post) HTML() string {
	if post.html != "" {
		return post.html
	}

	post.html = aero.Markdown(post.Text)
	return post.html
}

// ToPostable converts a post into an object that implements the Postable interface.
func (post *Post) ToPostable() Postable {
	return &PostPostable{post}
}

// GetPost ...
func GetPost(id string) (*Post, error) {
	obj, err := DB.Get("Post", id)
	return obj.(*Post), err
}

// GetPosts ...
func GetPosts() ([]*Post, error) {
	var posts []*Post

	channel := make(chan *Post)
	err := DB.Scan("Post", channel)

	for post := range channel {
		posts = append(posts, post)
	}

	return posts, err
}

// SortPostsLatestFirst sorts the slice of posts.
func SortPostsLatestFirst(posts []*Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Created > posts[j].Created
	})
}

// SortPostsLatestLast sorts the slice of posts.
func SortPostsLatestLast(posts []*Post) {
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Created < posts[j].Created
	})
}

// GetPostsByUser ...
func GetPostsByUser(user *User) ([]*Post, error) {
	var posts []*Post

	stream := make(chan *Post)
	err := DB.Scan("Post", stream)

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

// FilterPosts filters all forum posts by a custom function.
func FilterPosts(filter func(*Post) bool) ([]*Post, error) {
	var filtered []*Post

	channel := make(chan *Post)
	err := DB.Scan("Post", channel)

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
