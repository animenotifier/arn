package arn

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
	"github.com/aerogo/markdown"
	"github.com/animenotifier/arn/autocorrect"
)

// Force interface implementations
var (
	_ Postable          = (*Post)(nil)
	_ Likeable          = (*Post)(nil)
	_ LikeEventReceiver = (*Post)(nil)
	_ PostParent        = (*Post)(nil)
	_ fmt.Stringer      = (*Post)(nil)
	_ api.Newable       = (*Post)(nil)
	_ api.Editable      = (*Post)(nil)
	_ api.Actionable    = (*Post)(nil)
	_ api.Deletable     = (*Post)(nil)
)

// Actions
func init() {
	API.RegisterActions("Post", []*api.Action{
		// Like post
		LikeAction(),

		// Unlike post
		UnlikeAction(),
	})
}

// Authorize returns an error if the given API POST request is not authorized.
func (post *Post) Authorize(ctx *aero.Context, action string) error {
	if !ctx.HasSession() {
		return errors.New("Neither logged in nor in session")
	}

	if action == "edit" {
		user := GetUserFromContext(ctx)

		if post.CreatedBy != user.ID && user.Role != "admin" {
			return errors.New("Can't edit the posts of other users")
		}
	}

	return nil
}

// Create sets the data for a new post with data we received from the API request.
func (post *Post) Create(ctx *aero.Context) error {
	data, err := ctx.Request().Body().JSONObject()

	if err != nil {
		return err
	}

	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	post.ID = GenerateID("Post")
	post.Text, _ = data["text"].(string)
	post.CreatedBy = user.ID
	post.ParentID, _ = data["parentId"].(string)
	post.ParentType, _ = data["parentType"].(string)
	post.Likes = []string{}
	post.Created = DateTimeUTC()
	post.Edited = ""

	// Check parent type
	if !DB.HasType(post.ParentType) {
		return errors.New("Invalid parent type: " + post.ParentType)
	}

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

	// Thread
	parent := post.Parent()

	if parent == nil {
		return errors.New(post.ParentType + " does not exist")
	}

	// Is the parent locked?
	if IsLocked(parent) {
		return errors.New(post.ParentType + " is locked")
	}

	// Bind to local variable for the upcoming goroutine.
	posts := parent.Posts()

	// Notifications
	go func() {
		// Build a list of users to notify
		notifyUsers := map[string]bool{}

		// Mark the creator of the parent
		notifyUsers[parent.CreatorID()] = true

		// Mark every user who participated in the discussion
		for _, post := range posts {
			notifyUsers[post.CreatedBy] = true
		}

		// Exclude author of the new post
		delete(notifyUsers, post.CreatedBy)

		// Notify
		for notifyUserID := range notifyUsers {
			notifyUser, err := GetUser(notifyUserID)

			if notifyUser == nil || err != nil {
				continue
			}

			notifyUser.SendNotification(&PushNotification{
				Title:   user.Nick + " replied",
				Message: fmt.Sprintf("%s replied in the %s \"%s\".", user.Nick, strings.ToLower(post.ParentType), parent.TitleByUser(notifyUser)),
				Icon:    "https:" + user.AvatarLink("large"),
				Link:    post.Link(),
				Type:    NotificationTypeForumReply,
			})
		}
	}()

	// Append to posts
	parent.AddPost(post.ID)

	// Save the parent thread
	parent.Save()

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "create", "Post", post.ID, "", "", "")
	logEntry.Save()

	return nil
}

// Edit saves a log entry for the edit.
func (post *Post) Edit(ctx *aero.Context, key string, value reflect.Value, newValue reflect.Value) (bool, error) {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "edit", "Post", post.ID, key, fmt.Sprint(value.Interface()), fmt.Sprint(newValue.Interface()))
	logEntry.Save()

	return false, nil
}

// OnAppend saves a log entry.
func (post *Post) OnAppend(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayAppend", "Post", post.ID, fmt.Sprintf("%s[%d]", key, index), "", fmt.Sprint(obj))
	logEntry.Save()
}

// OnRemove saves a log entry.
func (post *Post) OnRemove(ctx *aero.Context, key string, index int, obj interface{}) {
	user := GetUserFromContext(ctx)
	logEntry := NewEditLogEntry(user.ID, "arrayRemove", "Post", post.ID, fmt.Sprintf("%s[%d]", key, index), fmt.Sprint(obj), "")
	logEntry.Save()
}

// DeleteInContext deletes the post in the given context.
func (post *Post) DeleteInContext(ctx *aero.Context) error {
	user := GetUserFromContext(ctx)

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "delete", "Post", post.ID, "", fmt.Sprint(post), "")
	logEntry.Save()

	return post.Delete()
}

// Delete deletes the post from the database.
func (post *Post) Delete() error {
	parent := post.Parent()

	if parent == nil {
		return fmt.Errorf("Invalid %s parent ID: %s", post.ParentType, post.ParentID)
	}

	// Remove the reference of the post in the thread that contains it
	if !parent.RemovePost(post.ID) {
		return fmt.Errorf("This post does not exist in the %s", strings.ToLower(post.ParentType))
	}

	parent.Save()
	DB.Delete("Post", post.ID)
	return nil
}

// AfterEdit updates the date it has been edited.
func (post *Post) AfterEdit(ctx *aero.Context) error {
	post.Edited = DateTimeUTC()
	post.html = markdown.Render(post.Text)
	return nil
}

// Save saves the post object in the database.
func (post *Post) Save() {
	DB.Set("Post", post.ID, post)
}
