package arn

import (
	"errors"
	"fmt"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
	"github.com/aerogo/markdown"
	"github.com/animenotifier/arn/autocorrect"
)

// Force interface implementations
var (
	_ Likeable          = (*GroupPost)(nil)
	_ LikeEventReceiver = (*GroupPost)(nil)
	_ fmt.Stringer      = (*GroupPost)(nil)
	_ api.Newable       = (*GroupPost)(nil)
	_ api.Editable      = (*GroupPost)(nil)
	_ api.Actionable    = (*GroupPost)(nil)
	_ api.Deletable     = (*GroupPost)(nil)
)

// Actions
func init() {
	API.RegisterActions("GroupPost", []*api.Action{
		// Like post
		LikeAction(),

		// Unlike post
		UnlikeAction(),
	})
}

// Authorize returns an error if the given API POST request is not authorized.
func (post *GroupPost) Authorize(ctx *aero.Context, action string) error {
	if !ctx.HasSession() {
		return errors.New("Neither logged in nor in session")
	}

	if action == "edit" {
		user := GetUserFromContext(ctx)

		if post.CreatedBy != user.ID {
			return errors.New("Can't edit the posts of other users")
		}
	}

	return nil
}

// Create sets the data for a new post with data we received from the API request.
func (post *GroupPost) Create(ctx *aero.Context) error {
	data, err := ctx.Request().Body().JSONObject()

	if err != nil {
		return err
	}

	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	post.ID = GenerateID("GroupPost")
	post.Text, _ = data["text"].(string)
	post.CreatedBy = user.ID
	post.GroupID = data["groupId"].(string)
	post.Likes = []string{}
	post.ChildIDs = []string{}
	post.Created = DateTimeUTC()
	post.Edited = ""

	// Post-process text
	post.Text = autocorrect.PostText(post.Text)

	if len(post.Text) < 5 {
		return errors.New("Text too short: Should be at least 5 characters")
	}

	// Tags
	tags, _ := data["tags"].([]interface{})
	post.Tags = make([]string, len(tags))

	for i := range post.Tags {
		post.Tags[i] = tags[i].(string)
	}

	_, err = GetGroup(post.GroupID)

	if err != nil {
		return errors.New("Group does not exist")
	}

	return nil
}

// DeleteInContext deletes the group post in the given context.
func (post *GroupPost) DeleteInContext(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "delete", "GroupPost", post.ID, "", fmt.Sprint(post), "")
	logEntry.Save()

	return post.Delete()
}

// Delete deletes the post from the database.
func (post *GroupPost) Delete() error {
	DB.Delete("GroupPost", post.ID)
	return nil
}

// AfterEdit updates the date it has been edited.
func (post *GroupPost) AfterEdit(ctx *aero.Context) error {
	post.Edited = DateTimeUTC()
	post.html = markdown.Render(post.Text)
	return nil
}

// Save saves the post object in the database.
func (post *GroupPost) Save() {
	DB.Set("GroupPost", post.ID, post)
}
