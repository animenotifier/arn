package arn

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/animenotifier/kitsu"
	"github.com/animenotifier/shoboi"
	"github.com/animenotifier/twist"
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
	Rating        *AnimeRating     `json:"rating"`
	Summary       string           `json:"summary"`
	Trailers      []*ExternalMedia `json:"trailers"`
	Mappings      []*Mapping       `json:"mappings"`

	// Episodes      []*AnimeEpisode  `json:"episodes"`

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
	episodes        *AnimeEpisodes
	upcomingEpisode *UpcomingEpisode
	characters      *AnimeCharacters
}

// AnimeImageTypes ...
type AnimeImageTypes struct {
	Tiny     string `json:"tiny"`
	Small    string `json:"small"`
	Large    string `json:"large"`
	Original string `json:"original"`
}

// GetAnime ...
func GetAnime(id string) (*Anime, error) {
	obj, err := DB.Get("Anime", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Anime), nil
}

// Characters ...
func (anime *Anime) Characters() *AnimeCharacters {
	if anime.characters != nil {
		return anime.characters
	}

	anime.characters, _ = GetAnimeCharacters(anime.ID)

	return anime.characters
}

// Link returns the URI to the anime page.
func (anime *Anime) Link() string {
	return "/anime/" + anime.ID
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
func (anime *Anime) AddMapping(serviceName string, serviceID string, userID string) {
	// Is the ID valid?
	if serviceID == "" {
		return
	}

	// If it already exists we don't need to add it
	for _, external := range anime.Mappings {
		if external.Service == serviceName && external.ServiceID == serviceID {
			return
		}
	}

	// Add the mapping
	anime.Mappings = append(anime.Mappings, &Mapping{
		Service:   serviceName,
		ServiceID: serviceID,
		Created:   DateTimeUTC(),
		CreatedBy: userID,
	})

	// Add the references
	switch serviceName {
	case "shoboi/anime":
		go anime.RefreshEpisodes()

	case "anilist/anime":
		DB.Set("AniListToAnime", serviceID, &AniListToAnime{
			AnimeID:   anime.ID,
			ServiceID: serviceID,
			Edited:    DateTimeUTC(),
			EditedBy:  userID,
		})

	case "myanimelist/anime":
		DB.Set("MyAnimeListToAnime", serviceID, &MyAnimeListToAnime{
			AnimeID:   anime.ID,
			ServiceID: serviceID,
			Edited:    DateTimeUTC(),
			EditedBy:  userID,
		})
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
	switch name {
	case "shoboi/anime":
		eps := anime.Episodes()

		if eps != nil {
			eps.Items = eps.Items[:0]
			eps.Save()
		}
	case "anilist/anime":
		DB.Delete("AniListToAnime", id)
	case "myanimelist/anime":
		DB.Delete("MyAnimeListToAnime", id)
	}

	for index, external := range anime.Mappings {
		if external.Service == name && external.ServiceID == id {
			anime.Mappings = append(anime.Mappings[:index], anime.Mappings[index+1:]...)
			return true
		}
	}

	return false
}

// Episodes returns the anime episodes wrapper.
func (anime *Anime) Episodes() *AnimeEpisodes {
	if anime.episodes == nil {
		record, err := DB.Get("AnimeEpisodes", anime.ID)

		if err != nil {
			return nil
		}

		anime.episodes = record.(*AnimeEpisodes)
	}

	return anime.episodes
}

// RefreshEpisodes will refresh the episode data.
func (anime *Anime) RefreshEpisodes() error {
	// Create blank episode templates
	episodes := anime.Episodes()
	episodes.Items = make([]*AnimeEpisode, anime.EpisodeCount, anime.EpisodeCount)

	for i := 0; i < len(episodes.Items); i++ {
		episodes.Items[i] = NewAnimeEpisode()
	}

	// Shoboi
	episodes.Merge(anime.ShoboiEpisodes())

	// AnimeTwist
	episodes.Merge(anime.TwistEpisodes())

	// Number remaining episodes
	startNumber := 0

	for _, episode := range episodes.Items {
		if episode.Number != -1 {
			startNumber = episode.Number
			continue
		}

		startNumber++
		episode.Number = startNumber
	}

	// Guess airing dates
	oneWeek := 7 * 24 * time.Hour
	lastAiringDate := ""
	timeDifference := oneWeek

	for _, episode := range episodes.Items {
		if episode.AiringDate.Start != "" {
			if lastAiringDate != "" {
				a, _ := time.Parse(time.RFC3339, lastAiringDate)
				b, _ := time.Parse(time.RFC3339, episode.AiringDate.Start)
				timeDifference = b.Sub(a)
			}

			lastAiringDate = episode.AiringDate.Start
			continue
		}

		// Add 1 week to the last known airing date
		nextAiringDate, _ := time.Parse(time.RFC3339, lastAiringDate)
		nextAiringDate = nextAiringDate.Add(timeDifference)

		// Guess start and end time
		episode.AiringDate.Start = nextAiringDate.Format(time.RFC3339)
		episode.AiringDate.End = nextAiringDate.Add(30 * time.Minute).Format(time.RFC3339)

		// Set this date as the new last known airing date
		lastAiringDate = episode.AiringDate.Start
	}

	return episodes.Save()
}

// ShoboiEpisodes returns a slice of episode info from cal.syoboi.jp.
func (anime *Anime) ShoboiEpisodes() []*AnimeEpisode {
	shoboiID := anime.GetMapping("shoboi/anime")

	if shoboiID == "" {
		return nil
	}

	shoboiAnime, err := shoboi.GetAnime(shoboiID)

	if err != nil {
		return nil
	}

	arnEpisodes := []*AnimeEpisode{}
	shoboiEpisodes := shoboiAnime.Episodes()

	for _, shoboiEpisode := range shoboiEpisodes {
		airingDate := shoboiEpisode.AiringDate()

		episode := NewAnimeEpisode()
		episode.Number = shoboiEpisode.Number
		episode.Title = &EpisodeTitle{
			Japanese: shoboiEpisode.TitleJapanese,
		}
		episode.AiringDate = &AnimeAiringDate{
			Start: airingDate.Start,
			End:   airingDate.End,
		}

		arnEpisodes = append(arnEpisodes, episode)
	}

	return arnEpisodes
}

// TwistEpisodes returns a slice of episode info from twist.moe.
func (anime *Anime) TwistEpisodes() []*AnimeEpisode {
	var cache ListOfIDs
	err := DB.GetObject("Cache", "animetwist index", &cache)

	if err != nil {
		return nil
	}

	// Does the index contain the ID?
	found := false

	for _, id := range cache.IDList {
		if id == anime.ID {
			found = true
			break
		}
	}

	// If the ID is not the index we don't need to query the feed
	if !found {
		return nil
	}

	// Get twist.moe feed
	feed, err := twist.GetFeedByKitsuID(anime.ID)

	if err != nil {
		return nil
	}

	episodes := feed.Episodes

	// Sort by episode number
	sort.Slice(episodes, func(a, b int) bool {
		return episodes[a].Number < episodes[b].Number
	})

	arnEpisodes := []*AnimeEpisode{}

	for _, episode := range episodes {
		arnEpisode := NewAnimeEpisode()
		arnEpisode.Number = episode.Number
		arnEpisode.Links = map[string]string{
			"twist.moe": strings.Replace(episode.Link, "https://test.twist.moe/", "https://twist.moe/", 1),
		}

		arnEpisodes = append(arnEpisodes, arnEpisode)
	}

	return arnEpisodes
}

// UpcomingEpisodes ...
func (anime *Anime) UpcomingEpisodes() []*UpcomingEpisode {
	// Special hack for K-On with ID 4240 because the episodes are way too far in the future
	if anime.ID == "4240" {
		return nil
	}

	var upcomingEpisodes []*UpcomingEpisode

	now := time.Now().UTC().Format(time.RFC3339)

	for _, episode := range anime.Episodes().Items {
		if episode.AiringDate.Start > now && episode.AiringDate.Start != invalidDate {
			upcomingEpisodes = append(upcomingEpisodes, &UpcomingEpisode{
				Anime:   anime,
				Episode: episode,
			})
		}
	}

	return upcomingEpisodes
}

// UpcomingEpisode ...
func (anime *Anime) UpcomingEpisode() *UpcomingEpisode {
	// Special hack for K-On with ID 4240 because the episodes are way too far in the future
	if anime.ID == "4240" {
		return nil
	}

	if anime.upcomingEpisode != nil {
		return anime.upcomingEpisode
	}

	now := time.Now().UTC().Format(time.RFC3339)

	for _, episode := range anime.Episodes().Items {
		if episode.AiringDate.Start > now && episode.AiringDate.Start != invalidDate {
			anime.upcomingEpisode = &UpcomingEpisode{
				Anime:   anime,
				Episode: episode,
			}

			return anime.upcomingEpisode
		}
	}

	return nil
}

// EpisodeCountString ...
func (anime *Anime) EpisodeCountString() string {
	if anime.EpisodeCount == 0 {
		return "?"
	}

	return strconv.Itoa(anime.EpisodeCount)
}

// EpisodeByNumber returns the episode with the given number.
func (anime *Anime) EpisodeByNumber(number int) *AnimeEpisode {
	for _, episode := range anime.Episodes().Items {
		if number == episode.Number {
			return episode
		}
	}

	return nil
}

// RefreshAnimeCharacters ...
func (anime *Anime) RefreshAnimeCharacters() error {
	resp, err := kitsu.GetAnimeCharactersForAnime(anime.ID)

	if err != nil {
		return err
	}

	animeCharacters := &AnimeCharacters{
		AnimeID: anime.ID,
		Items:   []*AnimeCharacter{},
	}

	for _, incl := range resp.Included {
		if incl.Type != "animeCharacters" {
			continue
		}

		role := incl.Attributes["role"].(string)
		characterID := incl.Relationships.Character.Data.ID

		fmt.Println(role, characterID)

		animeCharacter := &AnimeCharacter{
			CharacterID: characterID,
			Role:        role,
		}

		animeCharacters.Items = append(animeCharacters.Items, animeCharacter)
	}

	PrettyPrint(animeCharacters)

	return animeCharacters.Save()
}

// StreamAnime returns a stream of all anime.
func StreamAnime() (chan *Anime, error) {
	objects, err := DB.All("Anime")
	return objects.(chan *Anime), err
}

// MustStreamAnime returns a stream of all anime.
func MustStreamAnime() chan *Anime {
	stream, err := StreamAnime()
	PanicOnError(err)
	return stream
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
	beforeTime := time.Now().Add(-6 * 30 * 24 * time.Hour)
	beforeTimeString := beforeTime.Format(time.RFC3339)

	return FilterAnime(func(anime *Anime) bool {
		if (anime.Type != "tv" && anime.Type != "movie") || anime.NSFW == 1 || anime.StartDate < beforeTimeString {
			return false
		}

		// return anime.UpcomingEpisode() != nil || anime.Status == "upcoming"
		return anime.Status == "current" || anime.Status == "upcoming"
	})
}

// MustSave saves the anime in the database.
func (anime *Anime) MustSave() {
	PanicOnError(anime.Save())
}
