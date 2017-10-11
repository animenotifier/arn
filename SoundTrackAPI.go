package arn

import (
	"errors"
	"strings"

	"github.com/aerogo/aero"
	"github.com/animenotifier/arn/autocorrect"
)

// Authorize returns an error if the given API POST request is not authorized.
func (soundtrack *SoundTrack) Authorize(ctx *aero.Context, action string) error {
	if !ctx.HasSession() {
		return errors.New("Neither logged in nor in session")
	}

	return nil
}

// Create sets the data for a new soundtrack with data we received from the API request.
func (soundtrack *SoundTrack) Create(ctx *aero.Context) error {
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

	soundtrack.ID = GenerateID("SoundTrack")
	soundtrack.Likes = []string{}
	soundtrack.Created = DateTimeUTC()
	soundtrack.CreatedBy = user.ID
	soundtrack.Media = []*ExternalMedia{}

	// Soundcloud
	var soundcloud *ExternalMedia
	url, _ := data["soundcloud"].(string)

	if url != "" {
		soundcloud, err = GetSoundCloudMedia(url)

		if err != nil {
			return err
		}

		// Check that the track hasn't been posted yet
		_, err = DB.Get("SoundCloudToSoundTrack", soundcloud.ServiceID)

		if err == nil {
			return errors.New("This Soundcloud track has already been posted")
		}

		// Add to media
		soundtrack.Media = append(soundtrack.Media, soundcloud)
	}

	// Youtube
	var youtube *ExternalMedia
	url, _ = data["youtube"].(string)

	if url != "" {
		youtube, err = GetYoutubeMedia(url)

		if err != nil {
			return err
		}

		// Check that the video hasn't been posted yet
		_, err = DB.Get("YoutubeToSoundTrack", youtube.ServiceID)

		if err == nil {
			return errors.New("This Youtube video has already been posted")
		}

		// Add to media
		soundtrack.Media = append(soundtrack.Media, youtube)
	}

	// Tags
	tags, _ := data["tags"].([]interface{})
	soundtrack.Tags = make([]string, 0)

	animeFound := false
	for i := range tags {
		tag := tags[i].(string)
		tag = autocorrect.FixTag(tag)

		if strings.HasPrefix(tag, "anime:") {
			animeID := strings.TrimPrefix(tag, "anime:")
			_, err := GetAnime(animeID)

			if err != nil {
				return errors.New("Invalid anime ID")
			}

			animeFound = true
		}

		if tag != "" {
			soundtrack.Tags = append(soundtrack.Tags, tag)
		}
	}

	// No media added
	if len(soundtrack.Media) == 0 {
		return errors.New("No media specified (at least 1 media source is required)")
	}

	// No anime found
	if !animeFound {
		return errors.New("Need to specify at least one anime")
	}

	// No tags
	if len(tags) < 1 {
		return errors.New("Need to specify at least one tag")
	}

	// Save Soundcloud reference
	if soundcloud != nil {
		err = DB.Set("SoundCloudToSoundTrack", soundcloud.ServiceID, &SoundCloudToSoundTrack{
			ID:           soundcloud.ServiceID,
			SoundTrackID: soundtrack.ID,
		})

		if err != nil {
			return err
		}
	}

	// Save Youtube reference
	if youtube != nil {
		err = DB.Set("YoutubeToSoundTrack", youtube.ServiceID, &YoutubeToSoundTrack{
			ID:           youtube.ServiceID,
			SoundTrackID: soundtrack.ID,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// AfterEdit updates the metadata.
func (soundtrack *SoundTrack) AfterEdit(ctx *aero.Context) error {
	soundtrack.Edited = DateTimeUTC()
	soundtrack.EditedBy = GetUserFromContext(ctx).ID
	return nil
}

// Save saves the soundtrack object in the database.
func (soundtrack *SoundTrack) Save() error {
	return DB.Set("SoundTrack", soundtrack.ID, soundtrack)
}
