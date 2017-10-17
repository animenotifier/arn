package arn

import (
	"sort"

	"github.com/aerogo/markdown"
)

// Post represents a forum post.
type Post struct {
	ID       PostID   `json:"id"`
	Text     string   `json:"text" editable:"true"`
	AuthorID UserID   `json:"authorId"`
	ThreadID ThreadID `json:"threadId"`
	Tags     []string `json:"tags"`
	Likes    []string `json:"likes"`
	Created  UTCDate  `json:"created"`
	Edited   UTCDate  `json:"edited"`

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
	return "/post/" + post.ID
}

// HTML returns the HTML representation of the post.
func (post *Post) HTML() string {
	if post.html != "" {
		return post.html
	}

	post.html = markdown.Render(post.Text)
	return post.html
}

// ToPostable converts a post into an object that implements the Postable interface.
func (post *Post) ToPostable() Postable {
	return &PostPostable{post}
}

// GetPost ...
func GetPost(id string) (*Post, error) {
	obj, err := DB.Get("Post", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Post), nil
}

// StreamPosts returns a stream of all posts.
func StreamPosts() (chan *Post, error) {
	objects, err := DB.All("Post")
	return objects.(chan *Post), err
}

// AllPosts returns a slice of all posts.
func AllPosts() ([]*Post, error) {
	var all []*Post

	stream, err := StreamPosts()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
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

// FilterPostsWithUniqueThreads removes posts with the same thread until we have enough posts.
func FilterPostsWithUniqueThreads(posts []*Post, limit int) []*Post {
	filtered := []*Post{}
	threadsProcessed := map[string]bool{}

	for _, post := range posts {
		if len(filtered) >= limit {
			return filtered
		}

		_, found := threadsProcessed[post.ThreadID]

		if found {
			continue
		}

		threadsProcessed[post.ThreadID] = true
		filtered = append(filtered, post)
	}

	return filtered
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

// Like ...
func (post *Post) Like(userID string) {
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
			Message: likedBy.Nick + " liked your post in the thread \"" + post.Thread().Title + "\"",
			Icon:    "https:" + likedBy.LargeAvatar(),
			Link:    "https://notify.moe" + likedBy.Link(),
		})
	}()
}

// Unlike ...
func (post *Post) Unlike(userID string) {
	for index, id := range post.Likes {
		if id == userID {
			post.Likes = append(post.Likes[:index], post.Likes[index+1:]...)
			return
		}
	}
}
