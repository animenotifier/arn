package arn

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
)

func TestStreamAnime(t *testing.T) {
	PanicOnError(AniList.Authorize())
	color.Green(AniList.AccessToken)

	count := 0
	stream := AniList.StreamAnime()

	for anime := range stream {
		fmt.Println(anime.TitleRomaji)
		count++

		if count >= 80 {
			break
		}
	}
}
