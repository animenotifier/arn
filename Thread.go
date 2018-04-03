package arn

import (
	"sort"

	"github.com/aerogo/markdown"
	"github.com/aerogo/nano"
)

// Thread represents a forum thread.
type Thread struct {
	ID       string   `json:"id"`
	Title    string   `json:"title" editable:"true"`
	Text     string   `json:"text" editable:"true"`
	AuthorID string   `json:"authorId"`
	Sticky   int      `json:"sticky"`
	Tags     []string `json:"tags"`
	Posts    []string `json:"posts"`
	Created  string   `json:"created"`
	Edited   string   `json:"edited"`
	HasLikes

	html string
}

// Author returns the thread author.
func (thread *Thread) Author() *User {
	author, _ := GetUser(thread.AuthorID)
	return author
}

// Link returns the relative URL of the thread.
func (thread *Thread) Link() string {
	return "/thread/" + thread.ID
}

// HTML returns the HTML representation of the thread.
func (thread *Thread) HTML() string {
	if thread.html != "" {
		return thread.html
	}

	thread.html = markdown.Render(thread.Text)
	return thread.html
}

// String implements the default string serialization.
func (thread *Thread) String() string {
	return thread.Title
}

// OnLike is called when the thread receives a like.
func (thread *Thread) OnLike(likedBy *User) {
	if !thread.Author().Settings().Notification.ForumLikes {
		return
	}

	go func() {
		thread.Author().SendNotification(&PushNotification{
			Title:   likedBy.Nick + " liked your thread",
			Message: likedBy.Nick + " liked your thread \"" + thread.Title + "\"",
			Icon:    "https:" + likedBy.AvatarLink("large"),
			Link:    "https://notify.moe" + likedBy.Link(),
			Type:    NotificationTypeLike,
		})
	}()
}

// Remove post from the post list.
func (thread *Thread) Remove(postID string) bool {
	for index, item := range thread.Posts {
		if item == postID {
			thread.Posts = append(thread.Posts[:index], thread.Posts[index+1:]...)
			return true
		}
	}

	return false
}

// ToPostable converts a thread into an object that implements the Postable interface.
func (thread *Thread) ToPostable() Postable {
	return &ThreadPostable{thread}
}

// GetThread ...
func GetThread(id string) (*Thread, error) {
	obj, err := DB.Get("Thread", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Thread), nil
}

// GetThreadsByTag ...
func GetThreadsByTag(tag string) []*Thread {
	var threads []*Thread
	allTags := (tag == "" || tag == "<nil>")

	for thread := range StreamThreads() {
		if (allTags && !Contains(thread.Tags, "update")) || Contains(thread.Tags, tag) {
			threads = append(threads, thread)
		}
	}

	return threads
}

// GetThreadsByUser ...
func GetThreadsByUser(user *User) []*Thread {
	var threads []*Thread

	for thread := range StreamThreads() {
		if thread.AuthorID == user.ID {
			threads = append(threads, thread)
		}
	}

	return threads
}

// StreamThreads ...
func StreamThreads() chan *Thread {
	channel := make(chan *Thread, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Thread") {
			channel <- obj.(*Thread)
		}

		close(channel)
	}()

	return channel
}

// AllThreads ...
func AllThreads() []*Thread {
	var all []*Thread

	for obj := range StreamThreads() {
		all = append(all, obj)
	}

	return all
}

// SortThreads sorts a slice of threads for the forum view (stickies first).
func SortThreads(threads []*Thread) {
	sort.Slice(threads, func(i, j int) bool {
		a := threads[i]
		b := threads[j]

		if a.Sticky != b.Sticky {
			return a.Sticky > b.Sticky
		}

		return a.Created > b.Created
	})
}

// SortThreadsLatestFirst sorts a slice of threads by creation date.
func SortThreadsLatestFirst(threads []*Thread) {
	sort.Slice(threads, func(i, j int) bool {
		return threads[i].Created > threads[j].Created
	})
}
