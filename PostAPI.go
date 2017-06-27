package arn

import (
	"errors"

	"github.com/aerogo/aero"
)

// Authorize returns an error if the given API POST request is not authorized.
func (post *Post) Authorize(ctx *aero.Context) error {
	if !ctx.HasSession() {
		return errors.New("Neither logged in nor in session")
	}

	return nil
}

// Create sets the data for a new post with data we received from the API request.
func (post *Post) Create(json interface{}, ctx *aero.Context) error {
	data, formatOK := json.(map[string]interface{})

	if !formatOK {
		return errors.New("Invalid format (expected JSON)")
	}

	userID, ok := ctx.Session().Get("userId").(string)

	if !ok || userID == "" {
		return errors.New("Not logged in")
	}

	user, err := GetUser(userID)

	if err != nil {
		return err
	}

	post.ID = GenerateID("Post")
	post.Text, _ = data["text"].(string)
	post.AuthorID = user.ID
	post.ThreadID, _ = data["threadId"].(string)
	post.Likes = []string{}
	post.Created = DateTimeUTC()
	post.Edited = ""

	// Post-process text
	post.Text = FixPostText(post.Text)

	// Tags
	tags, _ := data["tags"].([]interface{})
	post.Tags = make([]string, len(tags))

	for i := range post.Tags {
		post.Tags[i] = tags[i].(string)
	}

	if len(post.Text) < 5 {
		return errors.New("Text too short: Should be at least 5 characters")
	}

	thread, threadErr := GetThread(post.ThreadID)

	if threadErr != nil {
		return errors.New("Thread does not exist")
	}

	thread.Posts = append(thread.Posts, post.ID)

	// Save the parent thread
	return thread.Save()
}

// Save saves the post object in the database.
func (post *Post) Save() error {
	return DB.Set("Post", post.ID, post)
}
