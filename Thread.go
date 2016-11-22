package arn

// Thread represents a forum thread.
type Thread struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Text     string   `json:"text"`
	AuthorID string   `json:"authorId"`
	Tags     []string `json:"tags"`
	Likes    []string `json:"likes"`
	Sticky   bool     `json:"sticky"`
	Replies  int      `json:"replies"`
	Created  string   `json:"created"`
	Edited   string   `json:"edited"`

	author *User
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

// ToPostable converts a thread into an object that implements the Postable interface.
func (thread *Thread) ToPostable() *ThreadPostable {
	return &ThreadPostable{thread}
}

// GetThread ...
func GetThread(id string) (*Thread, error) {
	thread := new(Thread)
	err := GetObject("Threads", id, thread)
	return thread, err
}

// GetThreadsByTag ...
func GetThreadsByTag(tag string) (ThreadList, error) {
	var threads ThreadList

	scan := make(chan *Thread)
	err := Scan("Threads", scan)

	for thread := range scan {
		if tag == "<nil>" || Contains(thread.Tags, tag) {
			threads = append(threads, thread)
		}
	}

	return threads, err
}
