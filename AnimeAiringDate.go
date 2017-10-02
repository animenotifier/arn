package arn

import (
	"time"
)

// AnimeAiringDate ...
type AnimeAiringDate struct {
	Start string `json:"start"`
	End   string `json:"end"`

	startHumanReadable string
	endHumanReadable   string
}

// StartDateHuman ...
func (airing *AnimeAiringDate) StartDateHuman() string {
	if airing.startHumanReadable == "" {
		t, _ := time.Parse(time.RFC3339, airing.Start)
		airing.startHumanReadable = t.Format(time.RFC1123)
	}

	return airing.startHumanReadable[:len("Thu, 25 May 2017")]
}

// EndDateHuman ...
func (airing *AnimeAiringDate) EndDateHuman() string {
	if airing.endHumanReadable == "" {
		t, _ := time.Parse(time.RFC3339, airing.End)
		airing.endHumanReadable = t.Format(time.RFC1123)
	}

	return airing.endHumanReadable[:len("Thu, 25 May 2017")]
}

// StartTimeHuman ...
func (airing *AnimeAiringDate) StartTimeHuman() string {
	if airing.startHumanReadable == "" {
		t, _ := time.Parse(time.RFC3339, airing.Start)
		airing.startHumanReadable = t.Format(time.RFC1123)
	}

	return airing.startHumanReadable[len("Thu, 25 May 2017 "):]
}

// EndTimeHuman ...
func (airing *AnimeAiringDate) EndTimeHuman() string {
	if airing.endHumanReadable == "" {
		t, _ := time.Parse(time.RFC3339, airing.End)
		airing.endHumanReadable = t.Format(time.RFC1123)
	}

	return airing.endHumanReadable[len("Thu, 25 May 2017 "):]
}
