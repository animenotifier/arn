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
	_ Likeable          = (*Post)(nil)
	_ LikeEventReceiver = (*Post)(nil)
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

		if post.AuthorID != user.ID {
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
	post.AuthorID = user.ID
	post.ThreadID, _ = data["threadId"].(string)
	post.Likes = []string{}
	post.Created = DateTimeUTC()
	post.Edited = ""

	// Post-process text
	post.Text = autocorrect.FixPostText(post.Text)

	if len(post.Text) < 5 {
		return errors.New("Text too short: Should be at least 5 characters")
	}

	// Tags
	tags, _ := data["tags"].([]interface{})
	post.Tags = make([]string, len(tags))

	for i := range post.Tags {
		post.Tags[i] = tags[i].(string)
	}

	thread, threadErr := GetThread(post.ThreadID)

	if threadErr != nil {
		return errors.New("Thread does not exist")
	}

	mentionedUsers := map[string]*User{}
	// Look for mentionedNicknames
	for _, match := range mentionNickRegex.FindAllStringSubmatch(post.Text, -1) {
		mentionedNickname := match[1]
		mentionedUser, err := GetUserByNick(mentionedNickname[1:])

		// Ignore the mention if the user is not found
		if err == nil {
			// Check if we've already passed the current mentionned user.
			if _, ok := mentionedUsers[mentionedUser.ID]; !ok {
				replacement := "${1}<@" + mentionedUser.ID + ">${2}"
				post.Text = TransformIDToMention(mentionedNickname, post.Text, replacement)
				mentionedUsers[mentionedUser.ID] = mentionedUser
			}
		}
	}

	// Bind to local variable for the upcoming goroutine.
	oldPosts := thread.Posts

	// Notifications
	go func() {
		postsObj := DB.GetMany("Post", oldPosts)
		posts := make([]*Post, len(postsObj), len(postsObj))

		for i, obj := range postsObj {
			posts[i] = obj.(*Post)
		}

		if err == nil {
			notifyUsers := map[string]bool{}
			notifyUsers[thread.AuthorID] = true

			for _, post := range posts {
				notifyUsers[post.AuthorID] = true
			}

			// Exclude author of the new post
			delete(notifyUsers, post.AuthorID)

			// Notify
			notification := &PushNotification{
				Title:   user.Nick + " replied",
				Message: fmt.Sprintf("%s replied in the thread \"%s\".", user.Nick, thread.Title),
				Icon:    "https:" + user.AvatarLink("large"),
				Link:    post.Link(),
				Type:    NotificationTypeForumReply,
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
	thread.Save()

	// Write log entry
	logEntry := NewEditLogEntry(user.ID, "create", "Post", post.ID, "", "", "")
	logEntry.Save()

	return nil
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
	thread, err := GetThread(post.ThreadID)

	if err != nil {
		return err
	}

	// Remove the reference of the post in the thread that contains it
	if !thread.Remove(post.ID) {
		return errors.New("This post does not exist in the thread")
	}

	thread.Save()
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
