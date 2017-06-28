package arn

import "time"

// AnimeAiringDate ...
type AnimeAiringDate struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// StartDateHuman ...
func (airing *AnimeAiringDate) StartDateHuman() string {
	t, _ := time.Parse(time.RFC3339, airing.Start)
	return t.Format(time.RFC1123)
}

// EndDateHuman ...
func (airing *AnimeAiringDate) EndDateHuman() string {
	t, _ := time.Parse(time.RFC3339, airing.End)
	return t.Format(time.RFC1123)
}
