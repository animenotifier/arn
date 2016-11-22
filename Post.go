package arn

// Post represents a forum post.
type Post struct {
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	AuthorID string   `json:"authorId"`
	ThreadID string   `json:"threadId"`
	Likes    []string `json:"likes"`
	Created  string   `json:"created"`
	Edited   string   `json:"edited"`

	author *User
	thread *Thread
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

// ToPostable converts a post into an object that implements the Postable interface.
func (post *Post) ToPostable() *PostPostable {
	return &PostPostable{post}
}

// GetPost ...
func GetPost(id string) (*Post, error) {
	post := new(Post)
	err := GetObject("Posts", id, post)
	return post, err
}

// GetPosts ...
func GetPosts() (PostList, error) {
	var posts PostList

	scan := make(chan *Post)
	err := Scan("Posts", scan)

	for post := range scan {
		posts = append(posts, post)
	}

	return posts, err
}

// FilterPosts filters all forum posts by a custom function.
func FilterPosts(filter func(*Post) bool) (PostList, error) {
	var filtered PostList

	channel := make(chan *Post)
	err := Scan("Posts", channel)

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
