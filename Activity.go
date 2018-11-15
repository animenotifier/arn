package arn

import (
	"sort"
	"sync"

	"github.com/aerogo/nano"
)

// Activity is a user activity that appears in the follower's feeds.
type Activity interface {
	Creator() *User
	TypeName() string
	GetID() string
	GetCreated() string
	GetCreatedBy() string
}

// SortActivitiesLatestFirst puts the latest entries on top.
func SortActivitiesLatestFirst(entries []Activity) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].GetCreated() > entries[j].GetCreated()
	})
}

// StreamActivities returns a stream of all activities.
func StreamActivities() chan Activity {
	channel := make(chan Activity, nano.ChannelBufferSize)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for obj := range DB.All("ActivityCreate") {
			channel <- obj.(Activity)
		}

		wg.Done()
	}()

	go func() {
		for obj := range DB.All("ActivityConsumeAnime") {
			channel <- obj.(Activity)
		}

		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(channel)
	}()

	return channel
}

// AllActivities returns a slice of all activities.
func AllActivities() []Activity {
	var all []Activity

	stream := StreamActivities()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}

// FilterActivities filters all Activities by a custom function.
func FilterActivities(filter func(Activity) bool) []Activity {
	var filtered []Activity

	for obj := range StreamActivities() {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered
}
