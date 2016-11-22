package arn

// Postable is a generic interface for Threads, Posts and Messages.
type Postable interface {
	ID() string
	Text() string
	Likes() []string
	Author() *User
	Link() string
	Type() string
}
