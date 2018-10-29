package arn

// HasPosts includes a list of Post IDs.
type HasPosts struct {
	PostIDs []string `json:"posts"`
}

// AddPost adds a post to the object.
func (obj *HasPosts) AddPost(postID string) {
	obj.PostIDs = append(obj.PostIDs, postID)
}

// RemovePost removes a post from the object.
func (obj *HasPosts) RemovePost(postID string) bool {
	for index, item := range obj.PostIDs {
		if item == postID {
			obj.PostIDs = append(obj.PostIDs[:index], obj.PostIDs[index+1:]...)
			return true
		}
	}

	return false
}

// Posts returns a slice of all posts.
func (obj *HasPosts) Posts() []*Post {
	objects := DB.GetMany("Post", obj.PostIDs)
	posts := []*Post{}

	for _, post := range objects {
		if post == nil {
			continue
		}

		posts = append(posts, post.(*Post))
	}

	return posts
}
