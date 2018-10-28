package arn

// PostParent is an interface that defines common functions for parent objects of posts.
type PostParent interface {
	Linkable
	TitleByUser(*User) string
}
