package arn_test

import (
	"testing"

	"github.com/animenotifier/arn"
)

func TestNormalizeRatings(t *testing.T) {
	user, _ := arn.GetUser("4J6qpK1ve")
	animeList := user.AnimeList()
	animeList.NormalizeRatings()
}
