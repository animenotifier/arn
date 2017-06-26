package arn

// Postable is a generic interface for Threads, Posts and Messages.
type Postable interface {
	ID() string
	Title() string
	Text() string
	HTML() string
	Likes() []string
	Author() *User
	Link() string
	Type() string
}

// CanBePostable is a type that defines the ToPostable() conversion.
type CanBePostable interface {
	ToPostable() Postable
}

// ToPostable converts a specific type to a generic postable.
func ToPostable(post CanBePostable) Postable {
	return post.ToPostable()
}

// ThreadsToPostables converts a slice of specific types to a slice of generic postables.
func ThreadsToPostables(threads []*Thread) []Postable {
	var postables []Postable

	for _, post := range threads {
		postables = append(postables, post.ToPostable())
	}

	return postables
}
