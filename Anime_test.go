package arn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamAnime(t *testing.T) {
	validAnimeStatus := []string{
		"finished",
		"current",
		"upcoming",
		"tba",
	}

	for anime := range StreamAnime() {
		assert.NotEmpty(t, anime.ID)
		assert.Contains(t, validAnimeStatus, anime.Status)
		assert.NotEmpty(t, anime.Link())

		anime.Episodes()
		anime.Characters()
		anime.GetMapping("shoboi/anime")
	}
}
