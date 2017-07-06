package arn

import (
	"sort"

	"github.com/aerogo/aero"
)

// Thread represents a forum thread.
type Thread struct {
	ID       string   `json:"id"`
	Title    string   `json:"title" editable:"true"`
	Text     string   `json:"text" editable:"true"`
	AuthorID string   `json:"authorId"`
	Sticky   int      `json:"sticky"`
	Tags     []string `json:"tags"`
	Likes    []string `json:"likes"`
	Posts    []string `json:"posts"`
	Created  string   `json:"created"`
	Edited   string   `json:"edited"`

	author *User
	html   string
}

// Author returns the thread author.
func (thread *Thread) Author() *User {
	if thread.author != nil {
		return thread.author
	}

	thread.author, _ = GetUser(thread.AuthorID)
	return thread.author
}

// Link returns the relative URL of the thread.
func (thread *Thread) Link() string {
	return "/threads/" + thread.ID
}

// HTML returns the HTML representation of the thread.
func (thread *Thread) HTML() string {
	if thread.html != "" {
		return thread.html
	}

	thread.html = aero.Markdown(thread.Text)
	return thread.html
}

// ToPostable converts a thread into an object that implements the Postable interface.
func (thread *Thread) ToPostable() Postable {
	return &ThreadPostable{thread}
}

// GetThread ...
func GetThread(id string) (*Thread, error) {
	obj, err := DB.Get("Thread", id)
	return obj.(*Thread), err
}

// GetThreadsByTag ...
func GetThreadsByTag(tag string) ([]*Thread, error) {
	var threads []*Thread

	scan := make(chan *Thread)
	err := DB.Scan("Thread", scan)
	allTags := (tag == "" || tag == "<nil>")

	for thread := range scan {
		if allTags || Contains(thread.Tags, tag) {
			threads = append(threads, thread)
		}
	}

	return threads, err
}

// GetThreadsByUser ...
func GetThreadsByUser(user *User) ([]*Thread, error) {
	var threads []*Thread

	scan := make(chan *Thread)
	err := DB.Scan("Thread", scan)

	for thread := range scan {
		if thread.AuthorID == user.ID {
			threads = append(threads, thread)
		}
	}

	return threads, err
}

// StreamThreads ...
func StreamThreads() (chan *Thread, error) {
	threads, err := DB.All("Thread")
	return threads.(chan *Thread), err
}

// AllThreads ...
func AllThreads() ([]*Thread, error) {
	var all []*Thread

	stream, err := StreamThreads()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
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
