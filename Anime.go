package arn

import (
	"encoding/json"
	"strconv"

	"github.com/animenotifier/shoboi"
	"github.com/fatih/color"
)

// Anime ...
type Anime struct {
	ID            string           `json:"id"`
	Type          string           `json:"type"`
	Title         *AnimeTitle      `json:"title"`
	Image         *AnimeImageTypes `json:"image"`
	FirstChannel  string           `json:"firstChannel"`
	StartDate     string           `json:"startDate"`
	EndDate       string           `json:"endDate"`
	EpisodeCount  int              `json:"episodeCount"`
	EpisodeLength int              `json:"episodeLength"`
	Status        string           `json:"status"`
	NSFW          int              `json:"nsfw"`
	Rating        AnimeRating      `json:"rating"`
	Summary       string           `json:"summary"`
	Trailers      []*ExternalMedia `json:"trailers"`
	Mappings      []*Mapping       `json:"mappings"`
	Episodes      []*AnimeEpisode  `json:"episodes"`

	// Adult         bool            `json:"adult"`

	// Hashtag       string          `json:"hashtag"`
	// Source        string          `json:"source"`

	// PageGenerated string          `json:"pageGenerated"`
	// AnilistEdited uint64          `json:"anilistEdited"`
	// Genres        []string        `json:"genres"`
	// Tracks        *AnimeTrackList `json:"tracks"`
	// Links         []AnimeLink     `json:"links"`
	// Studios       []AnimeStudio   `json:"studios"`
	// Relations     []AnimeRelation `json:"relations"`
	// Created       string          `json:"created"`
	// CreatedBy     string          `json:"createdBy"`
}

// AnimeRating ...
type AnimeRating struct {
	Overall    float64 `json:"overall" editable:"true"`
	Story      float64 `json:"story" editable:"true"`
	Visuals    float64 `json:"visuals" editable:"true"`
	Soundtrack float64 `json:"soundtrack" editable:"true"`
}

// AnimeImageTypes ...
type AnimeImageTypes struct {
	Tiny     string `json:"tiny"`
	Small    string `json:"small"`
	Large    string `json:"large"`
	Original string `json:"original"`
}

// AnimeTitle ...
type AnimeTitle struct {
	Romaji    string   `json:"romaji"`
	English   string   `json:"english"`
	Japanese  string   `json:"japanese"`
	Hiragana  string   `json:"hiragana"`
	Canonical string   `json:"canonical"`
	Synonyms  []string `json:"synonyms"`
}

// GetAnime ...
func GetAnime(id string) (*Anime, error) {
	obj, err := DB.Get("Anime", id)
	return obj.(*Anime), err
}

// Link returns the URI to the anime page.
func (anime *Anime) Link() string {
	return "/anime/" + anime.ID
}

// Save saves the anime in the database.
func (anime *Anime) Save() error {
	return DB.Set("Anime", anime.ID, anime)
}

// PrettyJSON ...
func (anime *Anime) PrettyJSON() (string, error) {
	data, err := json.MarshalIndent(anime, "", "    ")
	return string(data), err
}

// Watching ...
func (anime *Anime) Watching() int {
	return 0
}

// AddMapping adds the ID of an external site to the anime.
func (anime *Anime) AddMapping(name string, id string, userID string) {
	for _, external := range anime.Mappings {
		// If it already exists we don't need to add it
		if external.Service == name && external.ServiceID == id {
			return
		}
	}

	anime.Mappings = append(anime.Mappings, &Mapping{
		Service:   name,
		ServiceID: id,
		Created:   DateTimeUTC(),
		CreatedBy: userID,
	})

	go anime.RefreshMapping(name, id)
}

// RefreshMapping will refresh all the data we can get from the external site for that anime.
func (anime *Anime) RefreshMapping(serviceName string, serviceID string) {
	switch serviceName {
	case "shoboi/anime":
		shoboiAnime, err := shoboi.GetAnime(serviceID)

		if err != nil {
			color.Red(err.Error())
			return
		}

		anime.Episodes = []*AnimeEpisode{}

		shoboiEpisodes := shoboiAnime.Episodes()
		for _, shoboiEpisode := range shoboiEpisodes {
			airingDate := shoboiEpisode.AiringDate()

			episode := &AnimeEpisode{
				Number: shoboiEpisode.Number,
				Title: &EpisodeTitle{
					Japanese: shoboiEpisode.TitleJapanese,
				},
				AiringDate: &AnimeAiringDate{
					Start: airingDate.Start,
					End:   airingDate.End,
				},
			}

			anime.Episodes = append(anime.Episodes, episode)
		}

		anime.Save()
	}
}

// GetMapping returns the external ID for the given service.
func (anime *Anime) GetMapping(name string) string {
	for _, external := range anime.Mappings {
		if external.Service == name {
			return external.ServiceID
		}
	}

	return ""
}

// RemoveMapping removes all mappings with the given service name and ID.
func (anime *Anime) RemoveMapping(name string, id string) bool {
	for index, external := range anime.Mappings {
		if external.Service == name && external.ServiceID == id {
			anime.Mappings = append(anime.Mappings[:index], anime.Mappings[index+1:]...)
			return true
		}
	}

	return false
}

// EpisodeCountString ...
func (anime *Anime) EpisodeCountString() string {
	if anime.EpisodeCount == 0 {
		return "?"
	}

	return strconv.Itoa(anime.EpisodeCount)
}

// StreamAnime returns a stream of all anime.
func StreamAnime() (chan *Anime, error) {
	objects, err := DB.All("Anime")
	return objects.(chan *Anime), err
}

// AllAnime returns a slice of all anime.
func AllAnime() ([]*Anime, error) {
	var all []*Anime

	stream, err := StreamAnime()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}

// FilterAnime filters all anime by a custom function.
func FilterAnime(filter func(*Anime) bool) ([]*Anime, error) {
	var filtered []*Anime

	channel := make(chan *Anime)
	err := DB.Scan("Anime", channel)

	if err != nil {
		return filtered, err
	}

	for obj := range channel {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered, nil
}

// GetAiringAnime ...
func GetAiringAnime() ([]*Anime, error) {
	return FilterAnime(func(anime *Anime) bool {
		return anime.Status == "current" && anime.Type == "tv" && anime.NSFW == 0 && anime.Rating.Overall > 50
	})
}
