package arn

import (
	"net/http"
	"strings"

	"github.com/aerogo/http/client"
	"github.com/aerogo/nano"
	"github.com/animenotifier/kitsu"
	"github.com/fatih/color"
)

// NewAnimeFromKitsuAnime ...
func NewAnimeFromKitsuAnime(kitsuAnime *kitsu.Anime) (*Anime, *AnimeCharacters, *AnimeRelations, *AnimeEpisodes) {
	anime := NewAnime()
	attr := kitsuAnime.Attributes

	// General data
	anime.ID = kitsuAnime.ID
	anime.Type = strings.ToLower(attr.ShowType)
	anime.Title.Canonical = attr.CanonicalTitle
	anime.Title.English = attr.Titles.En
	anime.Title.Romaji = attr.Titles.EnJp
	anime.Title.Japanese = attr.Titles.JaJp
	anime.Title.Synonyms = attr.AbbreviatedTitles
	anime.StartDate = attr.StartDate
	anime.EndDate = attr.EndDate
	anime.EpisodeCount = attr.EpisodeCount
	anime.EpisodeLength = attr.EpisodeLength
	anime.Status = attr.Status
	anime.Summary = FixAnimeDescription(attr.Synopsis)

	// Status "unreleased" means the same as "upcoming" so we should normalize it
	if anime.Status == "unreleased" {
		anime.Status = "upcoming"
	}

	// // Import mappings
	// for _, mapping := range kitsuAnime.Mappings {
	// 	switch mapping.Attributes.ExternalSite {
	// 	case "myanimelist/anime":
	// 		anime.SetMapping("myanimelist/anime", mapping.Attributes.ExternalID, "")
	// 	case "anidb":
	// 		anime.SetMapping("anidb/anime", mapping.Attributes.ExternalID, "")
	// 	case "thetvdb", "thetvdb/series":
	// 		fmt.Println(mapping.Attributes.ExternalSite, mapping.Attributes.ExternalID)
	// 		anime.SetMapping("thetvdb/anime", mapping.Attributes.ExternalID, "")
	// 	case "thetvdb/season":
	// 		// Ignore
	// 	default:
	// 		color.Yellow("Unknown mapping: %s %s", mapping.Attributes.ExternalSite, mapping.Attributes.ExternalID)
	// 	}
	// }

	// Download image
	response, err := client.Get(attr.PosterImage.Original).End()

	if err == nil && response.StatusCode() == http.StatusOK {
		anime.SetImageBytes(response.Bytes())
	} else {
		color.Red("No image for [%s] %s (%d)", anime.ID, anime, response.StatusCode())
	}

	// Rating
	if anime.Rating.IsNotRated() {
		anime.Rating.Reset()
	}

	// Trailers
	if attr.YoutubeVideoID != "" {
		anime.Trailers = append(anime.Trailers, &ExternalMedia{
			Service:   "Youtube",
			ServiceID: attr.YoutubeVideoID,
		})
	}

	// Characters
	characters, _ := GetAnimeCharacters(anime.ID)

	if characters == nil {
		characters = &AnimeCharacters{
			AnimeID: anime.ID,
			Items:   []*AnimeCharacter{},
		}
	}

	// Episodes
	episodes, _ := GetAnimeEpisodes(anime.ID)

	if episodes == nil {
		episodes = &AnimeEpisodes{
			AnimeID: anime.ID,
			Items:   []*AnimeEpisode{},
		}
	}

	// Relations
	relations, _ := GetAnimeRelations(anime.ID)

	if relations == nil {
		relations = &AnimeRelations{
			AnimeID: anime.ID,
			Items:   []*AnimeRelation{},
		}
	}

	return anime, characters, relations, episodes
}

// StreamKitsuAnime returns a stream of all Kitsu anime.
func StreamKitsuAnime() chan *kitsu.Anime {
	channel := make(chan *kitsu.Anime, nano.ChannelBufferSize)

	go func() {
		for obj := range Kitsu.All("Anime") {
			channel <- obj.(*kitsu.Anime)
		}

		close(channel)
	}()

	return channel
}

// FilterKitsuAnime filters all Kitsu anime by a custom function.
func FilterKitsuAnime(filter func(*kitsu.Anime) bool) []*kitsu.Anime {
	var filtered []*kitsu.Anime

	channel := Kitsu.All("Anime")

	for obj := range channel {
		realObject := obj.(*kitsu.Anime)

		if filter(realObject) {
			filtered = append(filtered, realObject)
		}
	}

	return filtered
}

// AllKitsuAnime returns a slice of all Kitsu anime.
func AllKitsuAnime() []*kitsu.Anime {
	var all []*kitsu.Anime

	stream := StreamKitsuAnime()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}
