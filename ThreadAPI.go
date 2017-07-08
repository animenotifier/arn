package arn

import (
	"errors"

	"github.com/aerogo/aero"
	"github.com/animenotifier/arn/autocorrect"
)

// Authorize returns an error if the given API POST request is not authorized.
func (thread *Thread) Authorize(ctx *aero.Context) error {
	if !ctx.HasSession() {
		return errors.New("Neither logged in nor in session")
	}

	return nil
}

// Create sets the data for a new thread with data we received from the API request.
func (thread *Thread) Create(ctx *aero.Context) error {
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

	thread.ID = GenerateID("Thread")
	thread.Title, _ = data["title"].(string)
	thread.Text, _ = data["text"].(string)
	thread.AuthorID = user.ID
	thread.Sticky, _ = data["sticky"].(int)
	thread.Likes = []string{}
	thread.Posts = []string{}
	thread.Created = DateTimeUTC()
	thread.Edited = ""

	// Post-process text
	thread.Title = autocorrect.FixThreadTitle(thread.Title)
	thread.Text = autocorrect.FixPostText(thread.Text)

	// Tags
	tags, _ := data["tags"].([]interface{})
	thread.Tags = make([]string, len(tags))

	for i := range thread.Tags {
		thread.Tags[i] = tags[i].(string)
	}

	if len(tags) < 1 {
		return errors.New("Need to specify at least one tag")
	}

	if len(thread.Title) < 10 {
		return errors.New("Title too short: Should be at least 10 characters")
	}

	if len(thread.Text) < 10 {
		return errors.New("Text too short: Should be at least 10 characters")
	}

	return nil
}

// Update updates the thread object.
func (thread *Thread) Update(ctx *aero.Context, data interface{}) error {
	user := GetUserFromContext(ctx)

	if thread.AuthorID != user.ID {
		return errors.New("Can't edit the threads of other users")
	}

	thread.Edited = DateTimeUTC()

	updates := data.(map[string]interface{})
	return SetObjectProperties(thread, updates, nil)
}

// Action ...
func (thread *Thread) Action(ctx *aero.Context, action string) error {
	user := GetUserFromContext(ctx)

	if user == nil {
		return errors.New("Not logged in")
	}

	switch action {
	case "like":
		thread.Like(user.ID)
	case "unlike":
		thread.Unlike(user.ID)
	default:
		return errors.New("Unknown action: " + action)
	}
	return nil
}

// Save saves the thread object in the database.
func (thread *Thread) Save() error {
	return DB.Set("Thread", thread.ID, thread)
}
