package arn

import (
	"errors"
	"os"
	"os/exec"
	"path"
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
	IsDraft   bool             `json:"isDraft" editable:"true"`
	File      string           `json:"file"`
	Created   string           `json:"created"`
	CreatedBy string           `json:"createdBy"`
	Edited    string           `json:"edited"`
	EditedBy  string           `json:"editedBy"`
	LikeableImplementation
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
	var beatmaps []string

	for _, tag := range track.Tags {
		if strings.HasPrefix(tag, "osu-beatmap:") {
			osuID := strings.TrimPrefix(tag, "osu-beatmap:")
			beatmaps = append(beatmaps, osuID)
		}
	}

	return beatmaps
}

// MainAnime ...
func (track *SoundTrack) MainAnime() *Anime {
	allAnime := track.Anime()

	if len(allAnime) == 0 {
		return nil
	}

	return allAnime[0]
}

// Creator returns the user who created this track.
func (track *SoundTrack) Creator() *User {
	user, _ := GetUser(track.CreatedBy)
	return user
}

// EditedByUser returns the user who edited this track last.
func (track *SoundTrack) EditedByUser() *User {
	user, _ := GetUser(track.EditedBy)
	return user
}

// OnLike is called when the soundtrack receives a like.
func (track *SoundTrack) OnLike(likedBy *User) {
	if likedBy.ID == track.CreatedBy {
		return
	}

	go func() {
		track.Creator().SendNotification(&PushNotification{
			Title:   likedBy.Nick + " liked your soundtrack " + track.Title,
			Message: likedBy.Nick + " liked your soundtrack " + track.Title + ".",
			Icon:    "https:" + likedBy.AvatarLink("large"),
			Link:    "https://notify.moe" + likedBy.Link(),
			Type:    NotificationTypeLike,
		})
	}()
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

	draftIndex, err := GetDraftIndex(track.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.SoundTrackID == "" {
		return errors.New("Soundtrack draft doesn't exist in the user draft index")
	}

	track.IsDraft = false
	draftIndex.SoundTrackID = ""
	draftIndex.Save()
	return nil
}

// Unpublish ...
func (track *SoundTrack) Unpublish() error {
	draftIndex, err := GetDraftIndex(track.CreatedBy)

	if err != nil {
		return err
	}

	if draftIndex.SoundTrackID != "" {
		return errors.New("You still have an unfinished draft")
	}

	track.IsDraft = true
	draftIndex.SoundTrackID = track.ID
	draftIndex.Save()
	return nil
}

// Download downloads the track.
func (track *SoundTrack) Download() error {
	youtubeVideos := track.MediaByService("Youtube")

	if len(youtubeVideos) == 0 {
		return errors.New("No Youtube ID")
	}

	youtubeID := youtubeVideos[0].ServiceID

	// Check for existing file
	if track.File != "" {
		stat, err := os.Stat(path.Join(Root, "audio", track.File))

		if err == nil && !stat.IsDir() && stat.Size() > 0 {
			return errors.New("Already downloaded")
		}
	}

	audioDirectory := path.Join(Root, "audio")
	baseName := track.ID + "|" + youtubeID
	filePath := path.Join(audioDirectory, baseName)

	cmd := exec.Command("youtube-dl", "--extract-audio", "--audio-quality", "0", "--output", filePath+".%(ext)s", youtubeID)
	err := cmd.Start()

	if err != nil {
		return err
	}

	err = cmd.Wait()

	if err != nil {
		return err
	}

	fullPath := FindFileWithExtension(baseName, audioDirectory, []string{
		".opus",
		".webm",
		".ogg",
		".m4a",
		".mp3",
		".flac",
		".wav",
	})

	extension := path.Ext(fullPath)
	track.File = baseName + extension

	return nil
}

// String implements the default string serialization.
func (track *SoundTrack) String() string {
	return track.Title
}

// SortSoundTracksLatestFirst ...
func SortSoundTracksLatestFirst(tracks []*SoundTrack) {
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Created > tracks[j].Created
	})
}

// SortSoundTracksPopularFirst ...
func SortSoundTracksPopularFirst(tracks []*SoundTrack) {
	sort.Slice(tracks, func(i, j int) bool {
		aLikes := len(tracks[i].Likes)
		bLikes := len(tracks[j].Likes)

		if aLikes == bLikes {
			return tracks[i].Created > tracks[j].Created
		}

		return aLikes > bLikes
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
func AllSoundTracks() []*SoundTrack {
	var all []*SoundTrack

	for obj := range StreamSoundTracks() {
		all = append(all, obj)
	}

	return all
}

// FilterSoundTracks filters all soundtracks by a custom function.
func FilterSoundTracks(filter func(*SoundTrack) bool) []*SoundTrack {
	var filtered []*SoundTrack

	for obj := range StreamSoundTracks() {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered
}
