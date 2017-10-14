package arn

import (
	"errors"
	"fmt"

	"github.com/aerogo/aero"
	"github.com/aerogo/api"
	"github.com/animenotifier/arn/autocorrect"
)

// Force interface implementations
var (
	_ api.Newable    = (*Post)(nil)
	_ api.Editable   = (*Post)(nil)
	_ api.Actionable = (*Post)(nil)
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

		if post.AuthorID != user.ID {
			return errors.New("Can't edit the posts of other users")
		}
	}

	return nil
}

// Create sets the data for a new post with data we received from the API request.
func (post *Post) Create(ctx *aero.Context) error {
	data, err := ctx.RequestBodyJSONObject()

	if err != nil {
		return err
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
	post.Text = autocorrect.FixPostText(post.Text)

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

	// Bind to local variable for the upcoming goroutine.
	oldPosts := thread.Posts

	// Notifications
	go func() {
		postsObj, err := DB.GetMany("Post", oldPosts)
		posts := postsObj.([]*Post)

		if err == nil {
			notifyUsers := map[string]bool{}
			notifyUsers[thread.AuthorID] = true

			for _, post := range posts {
				notifyUsers[post.AuthorID] = true
			}

			// Exclude author of the new post
			delete(notifyUsers, post.AuthorID)

			// Notify
			notification := &Notification{
				Title:   user.Nick + " replied",
				Message: fmt.Sprintf("%s replied in the thread \"%s\".", user.Nick, thread.Title),
				Icon:    "https://notify.moe/images/brand/300",
				Link:    post.Link(),
			}

			for notifyUserID := range notifyUsers {
				notifyUser, err := GetUser(notifyUserID)

				if notifyUser == nil || err != nil {
					continue
				}

				notifyUser.SendNotification(notification)
			}
		}
	}()

	// Append to posts
	thread.Posts = append(thread.Posts, post.ID)

	// Save the parent thread
	return thread.Save()
}

// AfterEdit updates the date it has been edited.
func (post *Post) AfterEdit(ctx *aero.Context) error {
	post.Edited = DateTimeUTC()
	return nil
}

// Save saves the post object in the database.
func (post *Post) Save() error {
	return DB.Set("Post", post.ID, post)
}
