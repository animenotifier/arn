package arn

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aerogo/markdown"
	"github.com/aerogo/nano"
)

// Post is a comment related to any parent type in the database.
type Post struct {
	Text       string   `json:"text" editable:"true" type:"textarea"`
	Tags       []string `json:"tags" editable:"true"`
	ThreadID   string   `json:"threadId"` // DEPRECATED
	ParentID   string   `json:"parentId"`
	ParentType string   `json:"parentType"`
	Edited     string   `json:"edited"`

	HasID
	HasCreator
	HasLikes

	html string
}

// Thread returns the thread this post was posted in.
// DEPRECATED
func (post *Post) Thread() *Thread {
	thread, _ := GetThread(post.ParentID)
	return thread
}

// Parent returns the object this post was posted in.
func (post *Post) Parent() PostParent {
	obj, _ := DB.Get(post.ParentType, post.ParentID)
	return obj.(PostParent)
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

// String implements the default string serialization.
func (post *Post) String() string {
	const maxLen = 170

	if len(post.Text) > maxLen {
		return post.Text[:maxLen-3] + "..."
	}

	return post.Text
}

// OnLike is called when the post receives a like.
func (post *Post) OnLike(likedBy *User) {
	if !post.Creator().Settings().Notification.ForumLikes {
		return
	}

	go func() {
		post.Creator().SendNotification(&PushNotification{
			Title:   likedBy.Nick + " liked your post",
			Message: fmt.Sprintf(`%s liked your post in the %s "%s"`, likedBy.Nick, strings.ToLower(post.ParentType), post.Parent().TitleByUser(post.Creator())),
			Icon:    "https:" + likedBy.AvatarLink("large"),
			Link:    "https://notify.moe" + likedBy.Link(),
			Type:    NotificationTypeLike,
		})
	}()
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
func StreamPosts() chan *Post {
	channel := make(chan *Post, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Post") {
			channel <- obj.(*Post)
		}

		close(channel)
	}()

	return channel
}

// AllPosts returns a slice of all posts.
func AllPosts() []*Post {
	var all []*Post

	for obj := range StreamPosts() {
		all = append(all, obj)
	}

	return all
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

		_, found := threadsProcessed[post.ParentID]

		if found {
			continue
		}

		threadsProcessed[post.ParentID] = true
		filtered = append(filtered, post)
	}

	return filtered
}

// GetPostsByUser ...
func GetPostsByUser(user *User) ([]*Post, error) {
	var posts []*Post

	for post := range StreamPosts() {
		if post.CreatedBy == user.ID {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

// FilterPosts filters all forum posts by a custom function.
func FilterPosts(filter func(*Post) bool) ([]*Post, error) {
	var filtered []*Post

	for post := range StreamPosts() {
		if filter(post) {
			filtered = append(filtered, post)
		}
	}

	return filtered, nil
}
