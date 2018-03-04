package arn

import (
	"time"
)

// AnimeAiringDate represents the airing date of an anime.
type AnimeAiringDate struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// StartDateHuman returns the start date of the anime in human readable form.
func (airing *AnimeAiringDate) StartDateHuman() string {
	t, _ := time.Parse(time.RFC3339, airing.Start)
	humanReadable := t.Format(time.RFC1123)

	return humanReadable[:len("Thu, 25 May 2017")]
}

// EndDateHuman returns the end date of the anime in human readable form.
func (airing *AnimeAiringDate) EndDateHuman() string {
	t, _ := time.Parse(time.RFC3339, airing.End)
	humanReadable := t.Format(time.RFC1123)

	return humanReadable[:len("Thu, 25 May 2017")]
}

// StartTimeHuman returns the start time of the anime in human readable form.
func (airing *AnimeAiringDate) StartTimeHuman() string {
	t, _ := time.Parse(time.RFC3339, airing.Start)
	humanReadable := t.Format(time.RFC1123)

	return humanReadable[len("Thu, 25 May 2017 "):]
}

// EndTimeHuman returns the end time of the anime in human readable form.
func (airing *AnimeAiringDate) EndTimeHuman() string {
	t, _ := time.Parse(time.RFC3339, airing.End)
	humanReadable := t.Format(time.RFC1123)

	return humanReadable[len("Thu, 25 May 2017 "):]
}
