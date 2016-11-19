package arn

// Thread represents a forum thread.
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

// ToPostable converts a thread into an object that implements the Postable interface.
// Threads, posts and messages can be converted to the generic Postable type.
func (thread *Thread) ToPostable() *ThreadPostable {
	return &ThreadPostable{thread}
}

// ThreadPostable implements the Postable interface following Go naming convetions.
type ThreadPostable struct {
	thread *Thread
}

// ID returns the thread ID.
func (postable *ThreadPostable) ID() string {
	return postable.thread.ID
}

// Text returns the Markdown text.
func (postable *ThreadPostable) Text() string {
	return postable.thread.Text
}

// Author returns the user object representing the thread's author.
func (postable *ThreadPostable) Author() *User {
	return postable.thread.Author
}

// Likes returns an array of user IDs for the post.
func (postable *ThreadPostable) Likes() []string {
	return postable.thread.Likes
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
