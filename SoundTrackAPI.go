package arn

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/aerogo/aero"
	"github.com/parnurzeal/gorequest"
)

// SoundCloudToSoundTrack ...
type SoundCloudToSoundTrack struct {
	ID           string `json:"id"`
	SoundTrackID string `json:"soundTrackId"`
}

// Authorize returns an error if the given API POST request is not authorized.
func (soundtrack *SoundTrack) Authorize(ctx *aero.Context) error {
	if !ctx.HasSession() {
		return errors.New("Neither logged in nor in session")
	}

	return nil
}

// Create sets the data for a new soundtrack with data we received from the API request.
func (soundtrack *SoundTrack) Create(postBody interface{}, ctx *aero.Context) error {
	data, formatOK := postBody.(map[string]interface{})

	if !formatOK {
		return errors.New("Invalid format (expected JSON)")
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

	// Soundcloud ID
	url, _ := data["soundcloud"].(string)

	if url == "" {
		return errors.New("Need to specify a soundcloud link")
	}

	_, body, errs := gorequest.New().Get("https://api.soundcloud.com/resolve.json?url=" + url + "&client_id=" + APIKeys.SoundCloud.ID).EndBytes()

	if len(errs) > 0 {
		return errs[0]
	}

	var resp SoundCloudResolveResponse
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return err
	}

	soundCloudID := strconv.Itoa(resp.ID)

	soundtrack.Media = []*ExternalMedia{
		&ExternalMedia{
			Service:   "SoundCloud",
			ServiceID: soundCloudID,
		},
	}

	// Check that the track hasn't been posted yet
	_, err = DB.Get("SoundCloudToSoundTrack", soundCloudID)

	if err == nil {
		return errors.New("This track has already been posted")
	}

	// Tags
	tags, _ := data["tags"].([]interface{})
	soundtrack.Tags = make([]string, 0)

	animeFound := false
	for i := range tags {
		tag := tags[i].(string)
		tag = FixTag(tag)

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

	if !animeFound {
		return errors.New("Need to specify at least one anime")
	}

	if len(tags) < 1 {
		return errors.New("Need to specify at least one tag")
	}

	// Save reference
	return DB.Set("SoundCloudToSoundTrack", soundCloudID, &SoundCloudToSoundTrack{
		ID:           soundCloudID,
		SoundTrackID: soundtrack.ID,
	})
}

// Save saves the soundtrack object in the database.
func (soundtrack *SoundTrack) Save() error {
	return DB.Set("SoundTrack", soundtrack.ID, soundtrack)
}
