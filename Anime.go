package arn

import "encoding/json"

// Anime ...
type Anime struct {
	ID            string         `json:"id"`
	Type          string         `json:"type"`
	Title         AnimeTitle     `json:"title"`
	Image         ImageTypes     `json:"image"`
	StartDate     string         `json:"startDate"`
	EndDate       string         `json:"endDate"`
	EpisodeCount  int            `json:"episodeCount"`
	EpisodeLength int            `json:"episodeLength"`
	Summary       string         `json:"summary"`
	Trailers      []AnimeTrailer `json:"trailers"`

	// AiringStatus  string          `json:"airingStatus"`
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

// ImageTypes ...
type ImageTypes struct {
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
	anime := new(Anime)
	err := GetObject("Anime", id, anime)
	return anime, err
}

// Save ...
func (anime *Anime) Save() error {
	return SetObject("Anime", anime.ID, anime)
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

// FilterAnime filters all anime by a custom function.
func FilterAnime(filter func(*Anime) bool) ([]*Anime, error) {
	var filtered []*Anime

	channel := make(chan *Anime)
	err := Scan("Anime", channel)

	if err != nil {
		return filtered, err
	}

	for post := range channel {
		if filter(post) {
			filtered = append(filtered, post)
		}
	}

	return filtered, nil
}

// // GetAiringAnime ...
// func GetAiringAnime() ([]*Anime, error) {
// 	return FilterAnime(func(anime *Anime) bool {
// 		return anime.AiringStatus == "currently airing" && !anime.Adult
// 	})
// }
