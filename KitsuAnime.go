package arn

import (
	"github.com/aerogo/nano"
	"github.com/animenotifier/kitsu"
)

// StreamKitsuAnime returns a stream of all Kitsu anime.
func StreamKitsuAnime() chan *kitsu.Anime {
	channel := make(chan *kitsu.Anime, nano.ChannelBufferSize)

	go func() {
		for obj := range Kitsu.All("Anime") {
			channel <- obj.(*kitsu.Anime)
		}

		close(channel)
	}()

	return channel
}

// FilterKitsuAnime filters all Kitsu anime by a custom function.
func FilterKitsuAnime(filter func(*kitsu.Anime) bool) []*kitsu.Anime {
	var filtered []*kitsu.Anime

	channel := Kitsu.All("Anime")

	for obj := range channel {
		realObject := obj.(*kitsu.Anime)

		if filter(realObject) {
			filtered = append(filtered, realObject)
		}
	}

	return filtered
}

// AllKitsuAnime returns a slice of all Kitsu anime.
func AllKitsuAnime() []*kitsu.Anime {
	var all []*kitsu.Anime

	stream := StreamKitsuAnime()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}
