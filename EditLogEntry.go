package arn

import "github.com/aerogo/nano"

// EditLogEntry is an entry in the editor log.
type EditLogEntry struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	Action     string `json:"action"`
	ObjectType string `json:"objectType"` // The typename of what was edited
	ObjectID   string `json:"objectId"`   // The ID of what was edited
	Key        string `json:"key"`
	OldValue   string `json:"oldValue"`
	NewValue   string `json:"newValue"`
	Created    string `json:"created"`
}

// NewEditLogEntry ...
func NewEditLogEntry(userID, action, objectType, objectID, key, oldValue, newValue string) *EditLogEntry {
	return &EditLogEntry{
		ID:         GenerateID("EditLogEntry"),
		UserID:     userID,
		Action:     action,
		ObjectType: objectType,
		ObjectID:   objectID,
		Key:        key,
		OldValue:   oldValue,
		NewValue:   newValue,
		Created:    DateTimeUTC(),
	}
}

// User returns the user the log entry belongs to.
func (entry *EditLogEntry) User() *User {
	user, _ := GetUser(entry.UserID)
	return user
}

// StreamEditLogEntries returns a stream of all log entries.
func StreamEditLogEntries() chan *EditLogEntry {
	channel := make(chan *EditLogEntry, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("EditLogEntry") {
			channel <- obj.(*EditLogEntry)
		}

		close(channel)
	}()

	return channel
}

// AllEditLogEntries returns a slice of all log entries.
func AllEditLogEntries() []*EditLogEntry {
	var all []*EditLogEntry

	stream := StreamEditLogEntries()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}

// FilterEditLogEntries filters all log entries by a custom function.
func FilterEditLogEntries(filter func(*EditLogEntry) bool) []*EditLogEntry {
	var filtered []*EditLogEntry

	channel := DB.All("EditLogEntry")

	for obj := range channel {
		realObject := obj.(*EditLogEntry)

		if filter(realObject) {
			filtered = append(filtered, realObject)
		}
	}

	return filtered
}
