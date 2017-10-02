package arn

import (
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/animenotifier/arn/validator"

	"github.com/animenotifier/kitsu"
	"github.com/animenotifier/shoboi"
	"github.com/animenotifier/twist"
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
	Rating        *AnimeRating     `json:"rating"`
	Popularity    *AnimePopularity `json:"popularity"`
	Summary       string           `json:"summary"`
	Trailers      []*ExternalMedia `json:"trailers"`
	Mappings      []*Mapping       `json:"mappings"`

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

	if anime.characters != nil {
		// Sort by role
		sort.Slice(anime.characters.Items, func(i, j int) bool {
			// A little trick because "main" < "supporting"
			return anime.characters.Items[i].Role < anime.characters.Items[j].Role
		})
	}

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

// UsersWatchingOrPlanned returns a list of users who are watching the anime right now.
func (anime *Anime) UsersWatchingOrPlanned() []*User {
	users, err := FilterUsers(func(user *User) bool {
		obj, err := user.AnimeList().Get(anime.ID)

		if err != nil || obj == nil {
			return false
		}

		item := obj.(*AnimeListItem)
		return item.Status == AnimeListStatusWatching || item.Status == AnimeListStatusPlanned
	})

	if err != nil {
		return nil
	}

	return users
}

// RefreshEpisodes will refresh the episode data.
func (anime *Anime) RefreshEpisodes() error {
	// Fetch episodes
	episodes := anime.Episodes()

	if episodes == nil {
		episodes = &AnimeEpisodes{
			AnimeID: anime.ID,
			Items:   []*AnimeEpisode{},
		}
	}

	// Save number of available episodes for comparison later
	oldAvailableCount := episodes.AvailableCount()

	// Create blank episode templates
	episodes.Items = make([]*AnimeEpisode, anime.EpisodeCount, anime.EpisodeCount)

	for i := 0; i < len(episodes.Items); i++ {
		episodes.Items[i] = NewAnimeEpisode()
	}

	// Shoboi
	shoboiEpisodes, err := anime.ShoboiEpisodes()

	if err != nil {
		color.Red(err.Error())
	}

	episodes.Merge(shoboiEpisodes)

	// AnimeTwist
	twistEpisodes, err := anime.TwistEpisodes()

	if err != nil {
		color.Red(err.Error())
	}

	episodes.Merge(twistEpisodes)

	// Count number of available episodes
	newAvailableCount := episodes.AvailableCount()

	if anime.Status != "finished" && newAvailableCount > oldAvailableCount {
		notification := &Notification{
			Title:   anime.Title.Canonical,
			Message: "Episode " + strconv.Itoa(newAvailableCount) + " has been released!",
			Icon:    anime.Image.Small,
			Link:    "https://notify.moe" + anime.Link(),
		}

		// New episodes have been released.
		// Notify all users who are watching the anime.
		go func() {
			for _, user := range anime.UsersWatchingOrPlanned() {
				user.SendNotification(notification)
			}
		}()
	}

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
		if validator.IsValidDate(episode.AiringDate.Start) {
			if lastAiringDate != "" {
				a, _ := time.Parse(time.RFC3339, lastAiringDate)
				b, _ := time.Parse(time.RFC3339, episode.AiringDate.Start)
				timeDifference = b.Sub(a)

				// Cap time difference at one week
				if timeDifference > oneWeek {
					timeDifference = oneWeek
				}
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
func (anime *Anime) ShoboiEpisodes() ([]*AnimeEpisode, error) {
	shoboiID := anime.GetMapping("shoboi/anime")

	if shoboiID == "" {
		return nil, errors.New("Missing shoboi/anime mapping")
	}

	shoboiAnime, err := shoboi.GetAnime(shoboiID)

	if err != nil {
		return nil, err
	}

	arnEpisodes := []*AnimeEpisode{}
	shoboiEpisodes := shoboiAnime.Episodes()

	for _, shoboiEpisode := range shoboiEpisodes {
		episode := NewAnimeEpisode()
		episode.Number = shoboiEpisode.Number
		episode.Title = &EpisodeTitle{
			Japanese: shoboiEpisode.TitleJapanese,
		}

		// Try to get airing date
		airingDate := shoboiEpisode.AiringDate

		if airingDate != nil {
			episode.AiringDate = &AnimeAiringDate{
				Start: airingDate.Start,
				End:   airingDate.End,
			}
		} else {
			episode.AiringDate = &AnimeAiringDate{
				Start: "",
				End:   "",
			}
		}

		arnEpisodes = append(arnEpisodes, episode)
	}

	return arnEpisodes, nil
}

// TwistEpisodes returns a slice of episode info from twist.moe.
func (anime *Anime) TwistEpisodes() ([]*AnimeEpisode, error) {
	var cache ListOfIDs
	err := DB.GetObject("Cache", "animetwist index", &cache)

	if err != nil {
		return nil, err
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
		return nil, errors.New("Not available in twist.moe anime index")
	}

	// Get twist.moe feed
	feed, err := twist.GetFeedByKitsuID(anime.ID)

	if err != nil {
		return nil, err
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

	return arnEpisodes, nil
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
		if episode.AiringDate.Start > now && validator.IsValidDate(episode.AiringDate.Start) {
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
		if episode.AiringDate.Start > now && validator.IsValidDate(episode.AiringDate.Start) {
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
func (anime *Anime) RefreshAnimeCharacters() (*AnimeCharacters, error) {
	resp, err := kitsu.GetAnimeCharactersForAnime(anime.ID)

	if err != nil {
		return nil, err
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

		animeCharacter := &AnimeCharacter{
			CharacterID: characterID,
			Role:        role,
		}

		animeCharacters.Items = append(animeCharacters.Items, animeCharacter)
	}

	return animeCharacters, animeCharacters.Save()
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
