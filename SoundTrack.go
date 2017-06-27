package arn

import (
	"sort"
	"strings"

	"github.com/fatih/color"
)

// SoundTrack ...
type SoundTrack struct {
	ID        string           `json:"id"`
	Media     []*ExternalMedia `json:"media"`
	Tags      []string         `json:"tags"`
	Likes     []string         `json:"likes"`
	Created   string           `json:"created"`
	CreatedBy string           `json:"createdBy"`
	Edited    string           `json:"edited"`

	mainAnime     *Anime
	createdByUser *User
}

// ExternalMedia ...
type ExternalMedia struct {
	Service   string `json:"service"`
	ServiceID string `json:"serviceId"`
}

// Anime fetches all tagged anime of the sound track.
func (track *SoundTrack) Anime() []*Anime {
	var animeList []*Anime

	for _, tag := range track.Tags {
		if strings.HasPrefix(tag, "anime:") {
			animeID := strings.TrimPrefix(tag, "anime:")
			anime, err := GetAnime(animeID)

			if err != nil {
				color.Red("Error fetching anime: %v", err)
				continue
			}

			animeList = append(animeList, anime)
		}
	}

	return animeList
}

// MainAnime ...
func (track *SoundTrack) MainAnime() *Anime {
	if track.mainAnime != nil {
		return track.mainAnime
	}

	allAnime := track.Anime()

	if len(allAnime) == 0 {
		return nil
	}

	track.mainAnime = allAnime[0]
	return track.mainAnime
}

// CreatedByUser ...
func (track *SoundTrack) CreatedByUser() *User {
	if track.createdByUser != nil {
		return track.createdByUser
	}

	user, err := GetUser(track.CreatedBy)

	if err != nil {
		color.Red("Error fetching user: %v", err)
		return nil
	}

	track.createdByUser = user
	return track.createdByUser
}

// SortSoundTracksLatestFirst ...
func SortSoundTracksLatestFirst(tracks []*SoundTrack) {
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Created > tracks[j].Created
	})
}

// GetSoundTracksByUser ...
func GetSoundTracksByUser(user *User) ([]*SoundTrack, error) {
	var userTracks []*SoundTrack
	tracks, err := StreamSoundTracks()

	if err != nil {
		return nil, err
	}

	for track := range tracks {
		if track.CreatedBy == user.ID {
			userTracks = append(userTracks, track)
		}
	}

	return userTracks, nil
}

// StreamSoundTracks ...
func StreamSoundTracks() (chan *SoundTrack, error) {
	tracks, err := DB.All("SoundTrack")
	return tracks.(chan *SoundTrack), err
}

// AllSoundTracks ...
func AllSoundTracks() ([]*SoundTrack, error) {
	var all []*SoundTrack

	stream, err := StreamSoundTracks()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}
