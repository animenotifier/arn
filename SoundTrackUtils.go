package arn

import (
	"errors"
	"regexp"
)

var youtubeIDRegex = regexp.MustCompile(`youtu(?:.*\/v\/|.*v=|\.be\/)([A-Za-z0-9_-]{11})`)

// GetYoutubeMedia returns an ExternalMedia object for the given Youtube link.
func GetYoutubeMedia(url string) (*ExternalMedia, error) {
	matches := youtubeIDRegex.FindStringSubmatch(url)

	if len(matches) < 2 {
		return nil, errors.New("Invalid Youtube URL")
	}

	videoID := matches[1]

	media := &ExternalMedia{
		Service:   "Youtube",
		ServiceID: videoID,
	}

	return media, nil
}

// // GetSoundCloudMedia returns an ExternalMedia object for the given Soundcloud link.
// func GetSoundCloudMedia(url string) (*ExternalMedia, error) {
// 	var err error
// 	_, body, errs := gorequest.New().Get("https://api.soundcloud.com/resolve.json?url=" + url + "&client_id=" + APIKeys.SoundCloud.ID).EndBytes()

// 	if len(errs) > 0 {
// 		return nil, errs[0]
// 	}

// 	var soundcloud SoundCloudTrack
// 	err = jsoniter.Unmarshal(body, &soundcloud)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if soundcloud.ID == 0 {
// 		return nil, errors.New("Invalid Soundcloud response as the ID is not valid")
// 	}

// 	soundCloudID := strconv.Itoa(soundcloud.ID)

// 	return &ExternalMedia{
// 		Service:   "SoundCloud",
// 		ServiceID: soundCloudID,
// 	}, nil
// }
