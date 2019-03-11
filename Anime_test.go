package arn_test

import (
	"testing"

	"github.com/animenotifier/arn"
	"github.com/stretchr/testify/assert"
)

func TestNewAnime(t *testing.T) {
	anime := arn.NewAnime()
	assert.NotNil(t, anime)
	assert.NotEmpty(t, anime.ID)
	assert.NotEmpty(t, anime.Created)
}

func TestGetAnime(t *testing.T) {
	// Existing anime
	anime, err := arn.GetAnime("74y2cFiiR")
	assert.NoError(t, err)
	assert.NotNil(t, anime)
	assert.NotEmpty(t, anime.ID)
	assert.NotEmpty(t, anime.Title.Canonical)

	// Not existing anime
	anime, err = arn.GetAnime("does not exist")
	assert.Error(t, err)
	assert.Nil(t, anime)
}

func TestAllAnime(t *testing.T) {
	validAnimeStatus := []string{
		"finished",
		"current",
		"upcoming",
		"tba",
	}

	validAnimeType := []string{
		"tv",
		"movie",
		"ova",
		"ona",
		"special",
		"music",
	}

	allAnime := arn.AllAnime()

	for _, anime := range allAnime {
		assert.NotEmpty(t, anime.ID)
		assert.Contains(t, validAnimeStatus, anime.Status)
		assert.Contains(t, validAnimeType, anime.Type)
		assert.Contains(t, validAnimeStatus, anime.CalculatedStatus())
		assert.NotEmpty(t, anime.StatusHumanReadable())
		assert.NotEmpty(t, anime.TypeHumanReadable())
		assert.NotEmpty(t, anime.Link())
		assert.NotEmpty(t, anime.EpisodeCountString())

		anime.Episodes()
		anime.Characters()
		anime.StartDateTime()
		anime.EndDateTime()
		anime.HasImage()
		anime.GetMapping("shoboi/anime")
		anime.Studios()
		anime.Producers()
		anime.Licensors()
		anime.Prequels()
	}
}
