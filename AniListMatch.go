package arn

import "encoding/json"

// AniListMatch ...
type AniListMatch struct {
	AniListItem *AniListAnimeListItem `json:"anilistItem"`
	ARNAnime    *Anime                `json:"arnAnime"`
}

// JSON ...
func (match *AniListMatch) JSON() string {
	b, err := json.Marshal(match)
	PanicOnError(err)
	return string(b)
}
