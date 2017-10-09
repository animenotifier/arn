package arn

// ExternalMedia ...
type ExternalMedia struct {
	Service   string `json:"service"`
	ServiceID string `json:"serviceId"`
}

// EmbedLink ...
func (media *ExternalMedia) EmbedLink() string {
	switch media.Service {
	case "SoundCloud":
		return "https://w.soundcloud.com/player/?url=https://api.soundcloud.com/tracks/" + media.ServiceID + "?auto_play=false&hide_related=true&show_comments=false&show_user=false&show_reposts=false&visual=true"
	case "Youtube":
		return "https://www.youtube.com/embed/" + media.ServiceID
	default:
		return ""
	}
}

// // RefreshMetaData ...
// func (media *ExternalMedia) RefreshMetaData() {
// 	switch media.Service {
// 	case "SoundCloud":
// 		_, body, errs := gorequest.New().Get("https://api.soundcloud.com/tracks/" + media.ServiceID + ".json?client_id=" + APIKeys.SoundCloud.ID).EndBytes()

// 		if len(errs) > 0 {
// 			color.Red(errs[0].Error())
// 			return
// 		}

// 		var soundcloud SoundCloudTrack
// 		err := json.Unmarshal(body, &soundcloud)

// 		if err != nil {
// 			color.Red(err.Error())
// 			return
// 		}

// 		if soundcloud.Title != "" {
// 			media.Title = soundcloud.Title
// 		}

// 	case "Youtube":
// 		// Get title
// 		_, body, errs := gorequest.New().Get("https://www.googleapis.com/youtube/v3/videos?part=snippet&id=" + media.ServiceID + "&key=" + APIKeys.GoogleAPI.Key).EndBytes()

// 		if len(errs) > 0 {
// 			color.Red(errs[0].Error())
// 			return
// 		}

// 		var response youtube.VideoListResponse
// 		json.Unmarshal(body, &response)

// 		if len(response.Items) == 0 {
// 			color.Red("Youtube: Items field is empty")
// 			return
// 		}

// 		title := response.Items[0].Snippet.Title

// 		if title != "" {
// 			media.Title = title
// 		}
// 	}
// }
