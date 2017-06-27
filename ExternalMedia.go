package arn

// ExternalMedia ...
type ExternalMedia struct {
	Service   string `json:"service"`
	ServiceID string `json:"serviceId"`
	Title     string `json:"title"`
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
