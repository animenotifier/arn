package arn

import (
	"testing"

	"github.com/animenotifier/anilist"
	"github.com/stretchr/testify/assert"
)

func TestStreamAnime(t *testing.T) {
	assert.NoError(t, anilist.Authorize())

	count := 0
	stream := anilist.StreamAnime()

	for anime := range stream {
		assert.NotNil(t, anime)
		assert.NotEmpty(t, anime.TitleRomaji)
		count++

		if count >= 80 {
			break
		}
	}
}
