package search_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/animenotifier/arn/search"
)

// Run these search terms and expect the
// anime ID on the right as first result.
var tests = map[string]string{
	"lucky star":    "Pg9BcFmig",
	"dragon ball":   "hbih5KmmR",
	"dragon ball z": "ir-05Fmmg",
	"masotan":       "grdNhFiiR",
	"akame ga":      "iEaTpFiig",
}

func TestAnimeSearch(t *testing.T) {
	for term, expectedAnimeID := range tests {
		results := search.Anime(term, 1)
		assert.Len(t, results, 1)
		assert.Equal(t, results[0].ID, expectedAnimeID)
	}
}
