package arn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamAnime(t *testing.T) {
	allAnime, err := StreamAnime()

	assert.NoError(t, err)
	assert.NotNil(t, allAnime)

	validAnimeStatus := []string{
		"finished",
		"current",
		"tba",
		"upcoming",
		"unreleased",
	}

	for anime := range allAnime {
		assert.NotEmpty(t, anime.ID)
		assert.Contains(t, validAnimeStatus, anime.Status)
	}
}
