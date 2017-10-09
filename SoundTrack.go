package arn

import (
	"sort"
	"strings"

	"github.com/fatih/color"
)

// SoundTrack ...
type SoundTrack struct {
	ID        string           `json:"id"`
	Title     string           `json:"title" editable:"true"`
	Media     []*ExternalMedia `json:"media" editable:"true"`
	Tags      []string         `json:"tags" editable:"true"`
	Likes     []string         `json:"likes"`
	IsDraft   bool             `json:"isDraft"`
	Created   string           `json:"created"`
	CreatedBy string           `json:"createdBy"`
	Edited    string           `json:"edited"`
	EditedBy  string           `json:"editedBy"`

	mainAnime     *Anime
	createdByUser *User
}

// Link returns the permalink for the track.
func (track *SoundTrack) Link() string {
	return "/soundtrack/" + track.ID
}

// HasTag returns true if it contains the given tag.
func (track *SoundTrack) HasTag(search string) bool {
	for _, tag := range track.Tags {
		if tag == search {
			return true
		}
	}

	return false
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

// GetSoundTrack ...
func GetSoundTrack(id string) (*SoundTrack, error) {
	track, err := DB.Get("SoundTrack", id)

	if err != nil {
		return nil, err
	}

	return track.(*SoundTrack), nil
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

// GetSoundTracksByTag ...
func GetSoundTracksByTag(filterTag string) ([]*SoundTrack, error) {
	var filteredTracks []*SoundTrack
	tracks, err := StreamSoundTracks()

	if err != nil {
		return nil, err
	}

	for track := range tracks {
		for _, tag := range track.Tags {
			if tag == filterTag {
				filteredTracks = append(filteredTracks, track)
				break
			}
		}
	}

	return filteredTracks, nil
}

// StreamSoundTracks returns a stream of all soundtracks.
func StreamSoundTracks() (chan *SoundTrack, error) {
	tracks, err := DB.All("SoundTrack")
	return tracks.(chan *SoundTrack), err
}

// MustStreamSoundTracks returns a stream of all soundtracks.
func MustStreamSoundTracks() chan *SoundTrack {
	stream, err := StreamSoundTracks()
	PanicOnError(err)
	return stream
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

// FilterSoundTracks filters all soundtracks by a custom function.
func FilterSoundTracks(filter func(*SoundTrack) bool) ([]*SoundTrack, error) {
	var filtered []*SoundTrack

	channel, err := StreamSoundTracks()

	if err != nil {
		return filtered, err
	}

	for obj := range channel {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered, nil
}
