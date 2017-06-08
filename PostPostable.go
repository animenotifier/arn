package arn

// PostPostable implements the Postable interface following Go naming convetions.
type PostPostable struct {
	post *Post
}

// ID returns the post ID.
func (postable *PostPostable) ID() string {
	return postable.post.ID
}

// Title returns the title of thread this post belongs to.
func (postable *PostPostable) Title() string {
	return postable.post.Thread().Title
}

// Text returns the Markdown text.
func (postable *PostPostable) Text() string {
	return postable.post.Text
}

// Author returns the user object representing the post's author.
func (postable *PostPostable) Author() *User {
	return postable.post.Author()
}

// Likes returns an array of user IDs for the post.
func (postable *PostPostable) Likes() []string {
	return postable.post.Likes
}

// Link returns the relative URL of the post.
func (postable *PostPostable) Link() string {
	return postable.post.Link()
}

// Type returns the name of the underlying type.
func (postable *PostPostable) Type() string {
	return "Post"
}
