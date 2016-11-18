package arn

import "sort"

// Thread ...
type Thread struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Text     string   `json:"text"`
	Author   *User    `json:"-"`
	AuthorID string   `json:"authorId"`
	Tags     []string `json:"tags"`
	Likes    []string `json:"likes"`
	Sticky   bool     `json:"sticky"`
	Replies  int      `json:"replies"`
	Created  string   `json:"created"`
}

// ThreadList ...
type ThreadList []*Thread

func (c ThreadList) Len() int {
	return len(c)
}

func (c ThreadList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ThreadList) Less(i, j int) bool {
	a := c[i]
	b := c[j]

	if a.Sticky != b.Sticky {
		if a.Sticky {
			return true
		}

		if b.Sticky {
			return false
		}
	}

	return a.Created > b.Created
}

// GetThread ...
func GetThread(id string) (*Thread, error) {
	thread := new(Thread)
	err := GetObject("Threads", id, thread)
	return thread, err
}

// GetThreadsByTag ...
func GetThreadsByTag(tag string) ([]*Thread, error) {
	var threads ThreadList

	scan := make(chan *Thread)
	err := Scan("Threads", scan)

	for thread := range scan {
		if tag == "<nil>" || Contains(thread.Tags, tag) {
			threads = append(threads, thread)
		}
	}

	sort.Sort(threads)

	return threads, err
}
