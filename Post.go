package arn

// Post represents a forum post.
type Post struct {
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	Author   *User    `json:"-"`
	AuthorID string   `json:"authorId"`
	ThreadID string   `json:"threadId"`
	Likes    []string `json:"likes"`
	Created  string   `json:"created"`
	Edited   string   `json:"edited"`
}

// Init fetches additional post data like the author and thread objects.
func (post *Post) Init() {
	post.Author, _ = GetUser(post.AuthorID)
}

// ToPostable converts a post into an object that implements the Postable interface.
func (post *Post) ToPostable() *PostPostable {
	return &PostPostable{post}
}

// PostPostable implements the Postable interface following Go naming convetions.
type PostPostable struct {
	post *Post
}

// ID returns the post ID.
func (postable *PostPostable) ID() string {
	return postable.post.ID
}

// Text returns the Markdown text.
func (postable *PostPostable) Text() string {
	return postable.post.Text
}

// Author returns the user object representing the post's author.
func (postable *PostPostable) Author() *User {
	return postable.post.Author
}

// Likes returns an array of user IDs for the post.
func (postable *PostPostable) Likes() []string {
	return postable.post.Likes
}

// PostList ...
type PostList []*Post

func (c PostList) Len() int {
	return len(c)
}

func (c PostList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c PostList) Less(i, j int) bool {
	return c[i].Created < c[j].Created
}

// FilterPosts filters all forum posts by a custom function.
func FilterPosts(filter func(*Post) bool) (PostList, error) {
	var filtered PostList

	channel := make(chan *Post)
	_, err := client.ScanAllObjects(scanPolicy, channel, "arn", "Posts")

	if err != nil {
		return filtered, err
	}

	for post := range channel {
		if filter(post) {
			post.Init()
			filtered = append(filtered, post)
		}
	}

	return filtered, nil
}
