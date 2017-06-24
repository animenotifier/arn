package arn

import (
	"encoding/json"
	"strconv"
)

// NotFoundAnime is the dummy object representing
var NotFoundAnime = &Anime{
	ID:   "error",
	Type: "error",
	Title: AnimeTitle{
		Canonical: "Error",
		Romaji:    "Error",
		Japanese:  "Error",
	},
	Summary: "Error fetching anime data",
}

// Anime ...
type Anime struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	Title         AnimeTitle      `json:"title"`
	Image         AnimeImageTypes `json:"image"`
	StartDate     string          `json:"startDate"`
	EndDate       string          `json:"endDate"`
	EpisodeCount  int             `json:"episodeCount"`
	EpisodeLength int             `json:"episodeLength"`
	Status        string          `json:"status"`
	NSFW          bool            `json:"nsfw"`
	Rating        AnimeRating     `json:"rating"`
	Summary       string          `json:"summary"`
	Trailers      []AnimeTrailer  `json:"trailers"`

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
	Overall float64 `json:"overall" editable:"true"`
	Story   float64 `json:"story" editable:"true"`
	Visuals float64 `json:"visuals" editable:"true"`
	Music   float64 `json:"music" editable:"true"`
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
	Canonical string   `json:"canonical"`
	Synonyms  []string `json:"synonyms"`
}

// AnimeTrailer ...
type AnimeTrailer struct {
	Service string `json:"service"`
	VideoID string `json:"videoId"`
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

// EpisodeCountString ...
func (anime *Anime) EpisodeCountString() string {
	if anime.EpisodeCount == 0 {
		return "?"
	}

	return strconv.Itoa(anime.EpisodeCount)
}

// AllAnime returns a stream of all anime.
func AllAnime() (chan *Anime, error) {
	channel := make(chan *Anime)
	err := DB.Scan("Anime", channel)

	return channel, err
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
		return anime.Status == "current" && anime.Type == "tv" && !anime.NSFW && anime.Rating.Overall > 50
	})
}
