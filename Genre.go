package arn

// Genre ...
type Genre struct {
	Genre     string   `json:"genre"`
	AnimeList []*Anime `json:"animeList"`
}
