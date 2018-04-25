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

// HTML returns the Markdown text.
func (postable *PostPostable) HTML() string {
	return postable.post.HTML()
}

// Creator returns the user object representing the post's author.
func (postable *PostPostable) Creator() *User {
	return postable.post.Creator()
}

// Likes returns an array of user IDs for the post.
func (postable *PostPostable) Likes() []string {
	return postable.post.Likes
}

// LikedBy tells you whether the given user has liked the post.
func (postable *PostPostable) LikedBy(userID string) bool {
	return postable.post.LikedBy(userID)
}

// Link returns the relative URL of the post.
func (postable *PostPostable) Link() string {
	return postable.post.Link()
}

// Thread returns the thread the post belongs to.
func (postable *PostPostable) Thread() *Thread {
	return postable.post.Thread()
}

// ThreadID returns the thread the post belongs to.
func (postable *PostPostable) ThreadID() string {
	return postable.post.ThreadID
}

// Created returns the date the post has been created.
func (postable *PostPostable) Created() string {
	return postable.post.Created
}

// Type returns the name of the underlying type.
func (postable *PostPostable) Type() string {
	return "Post"
}
