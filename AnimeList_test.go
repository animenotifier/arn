package arn

import (
	"testing"
)

func TestNormalizeRatings(t *testing.T) {
	user, _ := GetUser("4J6qpK1ve")
	animeList := user.AnimeList()
	animeList.NormalizeRatings()
}
