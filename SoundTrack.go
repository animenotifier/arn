package arn

import (
	"errors"
	"sort"
	"strings"

	"github.com/aerogo/nano"
	"github.com/animenotifier/arn/autocorrect"
	"github.com/fatih/color"
)

// SoundTrack ...
type SoundTrack struct {
	ID        string           `json:"id"`
	Title     string           `json:"title" editable:"true"`
	Media     []*ExternalMedia `json:"media" editable:"true"`
	Tags      []string         `json:"tags" editable:"true" tooltip:"<ul><li><strong>anime:ID</strong> to connect it with anime</li><li><strong>opening</strong> for openings</li><li><strong>ending</strong> for endings</li><li><strong>cover</strong> for covers</li><li><strong>remix</strong> for remixes</li></ul>"`
	Likes     []string         `json:"likes"`
	IsDraft   bool             `json:"isDraft" editable:"true"`
	Created   string           `json:"created"`
	CreatedBy string           `json:"createdBy"`
	Edited    string           `json:"edited"`
	EditedBy  string           `json:"editedBy"`

	mainAnime    *Anime
	creator      *User
	editedByUser *User
	beatmaps     []string
}

// Link returns the permalink for the track.
func (track *SoundTrack) Link() string {
	return "/soundtrack/" + track.ID
}

// MediaByService ...
func (track *SoundTrack) MediaByService(service string) []*ExternalMedia {
	filtered := []*ExternalMedia{}

	for _, media := range track.Media {
		if media.Service == service {
			filtered = append(filtered, media)
		}
	}

	return filtered
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

// Beatmaps returns all osu beatmap IDs of the sound track.
func (track *SoundTrack) Beatmaps() []string {
	if track.beatmaps != nil {
		return track.beatmaps
	}

	for _, tag := range track.Tags {
		if strings.HasPrefix(tag, "osu-beatmap:") {
			osuID := strings.TrimPrefix(tag, "osu-beatmap:")
			track.beatmaps = append(track.beatmaps, osuID)
		}
	}

	return track.beatmaps
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

// Creator ...
func (track *SoundTrack) Creator() *User {
	if track.creator != nil {
		return track.creator
	}

	user, err := GetUser(track.CreatedBy)

	if err != nil {
		color.Red("Error fetching user: %v", err)
		return nil
	}

	track.creator = user
	return track.creator
}

// EditedByUser returns the user who edited this track last.
func (track *SoundTrack) EditedByUser() *User {
	if track.editedByUser != nil {
		return track.editedByUser
	}

	user, err := GetUser(track.EditedBy)

	if err != nil {
		color.Red("Error fetching user: %v", err)
		return nil
	}

	track.editedByUser = user
	return track.editedByUser
}

// Publish ...
func (track *SoundTrack) Publish() error {
	// No draft
	if !track.IsDraft {
		return errors.New("Not a draft")
	}

	// No media added
	if len(track.Media) == 0 {
		return errors.New("No media specified (at least 1 media source is required)")
	}

	animeFound := false

	for _, tag := range track.Tags {
		tag = autocorrect.FixTag(tag)

		if strings.HasPrefix(tag, "anime:") {
			animeID := strings.TrimPrefix(tag, "anime:")
			_, err := GetAnime(animeID)

			if err != nil {
				return errors.New("Invalid anime ID")
			}

			animeFound = true
		}
	}

	// No anime found
	if !animeFound {
		return errors.New("Need to specify at least one anime")
	}

	// No tags
	if len(track.Tags) < 1 {
		return errors.New("Need to specify at least one tag")
	}

	track.IsDraft = false
	draftIndex, err := GetDraftIndex(track.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.SoundTrackID == "" {
		return errors.New("Soundtrack draft doesn't exist in the user draft index")
	}

	draftIndex.SoundTrackID = ""
	draftIndex.Save()
	return nil
}

// Unpublish ...
func (track *SoundTrack) Unpublish() error {
	track.IsDraft = true
	draftIndex, err := GetDraftIndex(track.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.SoundTrackID != "" {
		return errors.New("You still have an unfinished draft")
	}

	draftIndex.SoundTrackID = track.ID
	draftIndex.Save()
	return nil
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

// StreamSoundTracks returns a stream of all soundtracks.
func StreamSoundTracks() chan *SoundTrack {
	channel := make(chan *SoundTrack, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("SoundTrack") {
			channel <- obj.(*SoundTrack)
		}

		close(channel)
	}()

	return channel
}

// AllSoundTracks ...
func AllSoundTracks() ([]*SoundTrack, error) {
	var all []*SoundTrack

	for obj := range StreamSoundTracks() {
		all = append(all, obj)
	}

	return all, nil
}

// FilterSoundTracks filters all soundtracks by a custom function.
func FilterSoundTracks(filter func(*SoundTrack) bool) ([]*SoundTrack, error) {
	var filtered []*SoundTrack

	for obj := range StreamSoundTracks() {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered, nil
}

// Like adds an user to the track's Like array if they aren't already in it.
func (track *SoundTrack) Like(userID string) {
	for _, id := range track.Likes {
		if id == userID {
			return
		}
	}

	track.Likes = append(track.Likes, userID)
}

// Unlike removes the user from the track's Likes array if they are in it.
func (track *SoundTrack) Unlike(userID string) {
	for index, id := range track.Likes {
		if id == userID {
			track.Likes = append(track.Likes[:index], track.Likes[index+1:]...)
			return
		}
	}
}
