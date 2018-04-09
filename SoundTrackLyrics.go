package arn

// SoundTrackLyrics represents song lyrics.
type SoundTrackLyrics struct {
	Native string `json:"native" editable:"true" type:"textarea"`
	Romaji string `json:"romaji" editable:"true" type:"textarea"`
	// English string `json:"english" editable:"true" type:"textarea"`
}
