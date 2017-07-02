package arn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamAnime(t *testing.T) {
	assert.Nil(t, AniList.Authorize())

	count := 0
	stream := AniList.StreamAnime()

	for anime := range stream {
		assert.NotNil(t, anime)
		assert.NotEmpty(t, anime.TitleRomaji)
		count++

		if count >= 80 {
			break
		}
	}
}
