package arn

// ThreadPostable implements the Postable interface following Go naming convetions.
type ThreadPostable struct {
	thread *Thread
}

// ID returns the thread ID.
func (postable *ThreadPostable) ID() string {
	return postable.thread.ID
}

// Title returns the thread title.
func (postable *ThreadPostable) Title() string {
	return postable.thread.Title
}

// Text returns the Markdown text.
func (postable *ThreadPostable) Text() string {
	return postable.thread.Text
}

// HTML returns the Markdown text.
func (postable *ThreadPostable) HTML() string {
	return postable.thread.HTML()
}

// Creator returns the user object representing the thread's author.
func (postable *ThreadPostable) Creator() *User {
	return postable.thread.Creator()
}

// Likes returns an array of user IDs for the post.
func (postable *ThreadPostable) Likes() []string {
	return postable.thread.Likes
}

// LikedBy tells you whether the given user has liked the thread.
func (postable *ThreadPostable) LikedBy(userID string) bool {
	return postable.thread.LikedBy(userID)
}

// Link returns the relative URL of the thread.
func (postable *ThreadPostable) Link() string {
	return postable.thread.Link()
}

// Thread returns the internal thread object.
func (postable *ThreadPostable) Thread() *Thread {
	return postable.thread
}

// ThreadID returns the thread ID.
func (postable *ThreadPostable) ThreadID() string {
	return postable.thread.ID
}

// Created returns the date the thread has been created.
func (postable *ThreadPostable) Created() string {
	return postable.thread.Created
}

// Type returns the name of the underlying type.
func (postable *ThreadPostable) Type() string {
	return "Thread"
}
