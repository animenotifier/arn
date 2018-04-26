package search_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/animenotifier/arn/search"
)

// Run these search terms and expect the
// anime ID on the right as first result.
var tests = map[string]string{
	"lucky star":     "Pg9BcFmig", // Lucky☆Star
	"dragn bll":      "hbih5KmmR", // Dragon Ball
	"dragon ball":    "hbih5KmmR", // Dragon Ball
	"dragon ball z":  "ir-05Fmmg", // Dragon Ball Z
	"masotan":        "grdNhFiiR", // Hisone to Maso-tan
	"akame":          "iEaTpFiig", // Akame ga Kill!
	"kimi":           "7VjCpFiiR", // Kimi no Na wa.
	"working":        "0iIgtFimg", // Working!!
	"k on":           "LP8j5Kmig", // K-On!
	"kon":            "LP8j5Kmig", // K-On!
	"danmachi":       "LTTPtKmiR", // Dungeon ni Deai wo Motomeru no wa Machigatteiru Darou ka
	"sword oratoria": "ifGetFmig", // Dungeon ni Deai wo Motomeru no wa Machigatteiru Darou ka Gaiden: Sword Oratoria
	"gint":           "QAZ1cKmig", // Gintama
	"k":              "EDSOtKmig", // K
	"champloo":       "0ER25Fiig", // Samurai Champloo
	"one peace":      "jdZp5KmiR", // One Piece
}

func TestAnimeSearch(t *testing.T) {
	for term, expectedAnimeID := range tests {
		results := search.Anime(term, 1)
		assert.Len(t, results, 1)
		assert.Equal(t, expectedAnimeID, results[0].ID)
	}
}
