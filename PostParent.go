package arn

// PostParent is an interface that defines common functions for parent objects of posts.
type PostParent interface {
	Linkable
	Lockable
	GetID() string
	TitleByUser(*User) string
	Posts() []*Post
	Creator() *User
	CreatorID() string
	AddPost(string)
	RemovePost(string) bool
}
