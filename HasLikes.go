package arn

// HasLikes implements common like and unlike methods.
type HasLikes struct {
	Likes []string `json:"likes"`
}

// Like makes the given user ID like the object.
func (obj *HasLikes) Like(userID string) {
	for _, id := range obj.Likes {
		if id == userID {
			return
		}
	}

	obj.Likes = append(obj.Likes, userID)
}

// Unlike makes the given user ID unlike the object.
func (obj *HasLikes) Unlike(userID string) {
	for index, id := range obj.Likes {
		if id == userID {
			obj.Likes = append(obj.Likes[:index], obj.Likes[index+1:]...)
			return
		}
	}
}

// LikedBy checks to see if the user has liked the object.
func (obj *HasLikes) LikedBy(userID string) bool {
	for _, id := range obj.Likes {
		if id == userID {
			return true
		}
	}

	return false
}

// CountLikes returns the number of likes the object has received.
func (obj *HasLikes) CountLikes() int {
	return len(obj.Likes)
}
