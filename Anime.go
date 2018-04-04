package arn

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aerogo/nano"
	"github.com/animenotifier/arn/validator"
	"github.com/animenotifier/twist"

	"github.com/animenotifier/kitsu"
	"github.com/animenotifier/shoboi"
	"github.com/fatih/color"
)

// AnimeDateFormat describes the anime date format for the date conversion.
const AnimeDateFormat = "2006-01-02"

// AnimeSourceHumanReadable maps the anime source to a human readable version.
var AnimeSourceHumanReadable = map[string]string{}

// Register a list of supported anime status and source types.
func init() {
	DataLists["anime-types"] = []*Option{
		&Option{"tv", "TV"},
		&Option{"movie", "Movie"},
		&Option{"ova", "OVA"},
		&Option{"ona", "ONA"},
		&Option{"special", "Special"},
		&Option{"music", "Music"},
	}

	DataLists["anime-status"] = []*Option{
		&Option{"current", "Current"},
		&Option{"finished", "Finished"},
		&Option{"upcoming", "Upcoming"},
		&Option{"tba", "To be announced"},
	}

	DataLists["anime-sources"] = []*Option{
		&Option{"", "Unknown"},
		&Option{"original", "Original"},
		&Option{"manga", "Manga"},
		&Option{"novel", "Novel"},
		&Option{"light novel", "Light novel"},
		&Option{"visual novel", "Visual novel"},
		&Option{"game", "Game"},
		&Option{"book", "Book"},
		&Option{"4-koma manga", "4-koma Manga"},
		&Option{"music", "Music"},
		&Option{"picture book", "Picture book"},
		&Option{"web manga", "Web manga"},
		&Option{"other", "Other"},
	}

	for _, option := range DataLists["anime-source"] {
		AnimeSourceHumanReadable[option.Value] = option.Label
	}
}

// Anime represents an anime.
type Anime struct {
	ID            string           `json:"id"`
	Type          string           `json:"type" editable:"true" datalist:"anime-types"`
	Title         *AnimeTitle      `json:"title" editable:"true"`
	Summary       string           `json:"summary" editable:"true" type:"textarea"`
	Status        string           `json:"status" editable:"true" datalist:"anime-status"`
	Genres        []string         `json:"genres" editable:"true"`
	StartDate     string           `json:"startDate" editable:"true"`
	EndDate       string           `json:"endDate" editable:"true"`
	EpisodeCount  int              `json:"episodeCount" editable:"true"`
	EpisodeLength int              `json:"episodeLength" editable:"true"`
	Source        string           `json:"source" editable:"true" datalist:"anime-sources"`
	Image         AnimeImage       `json:"image"`
	FirstChannel  string           `json:"firstChannel"`
	Rating        *AnimeRating     `json:"rating"`
	Popularity    *AnimePopularity `json:"popularity"`
	Trailers      []*ExternalMedia `json:"trailers" editable:"true"`

	// Mixins
	HasMappings

	// Company IDs
	StudioIDs   []string `json:"studios" editable:"true"`
	ProducerIDs []string `json:"producers" editable:"true"`
	LicensorIDs []string `json:"licensors" editable:"true"`

	// Links to external websites
	Links []*Link `json:"links" editable:"true"`

	// Editing dates
	Created   string `json:"created"`
	CreatedBy string `json:"createdBy"`
	Edited    string `json:"edited"`
	EditedBy  string `json:"editedBy"`

	// SynopsisSource string        `json:"synopsisSource" editable:"true"`
	// Hashtag        string        `json:"hashtag"`
	// Created        string        `json:"created"`
	// CreatedBy      string        `json:"createdBy"`
}

// NewAnime creates a new anime.
func NewAnime() *Anime {
	return &Anime{
		ID:         GenerateID("Anime"),
		Title:      &AnimeTitle{},
		Rating:     &AnimeRating{},
		Popularity: &AnimePopularity{},
		Trailers:   []*ExternalMedia{},
		Created:    DateTimeUTC(),
		HasMappings: HasMappings{
			Mappings: []*Mapping{},
		},
	}
}

// GetAnime ...
func GetAnime(id string) (*Anime, error) {
	obj, err := DB.Get("Anime", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Anime), nil
}

// AddStudio adds the company ID to the studio ID list if it doesn't exist already.
func (anime *Anime) AddStudio(companyID string) {
	// Is the ID valid?
	if companyID == "" {
		return
	}

	// If it already exists we don't need to add it
	for _, id := range anime.StudioIDs {
		if id == companyID {
			return
		}
	}

	anime.StudioIDs = append(anime.StudioIDs, companyID)
}

// AddProducer adds the company ID to the producer ID list if it doesn't exist already.
func (anime *Anime) AddProducer(companyID string) {
	// Is the ID valid?
	if companyID == "" {
		return
	}

	// If it already exists we don't need to add it
	for _, id := range anime.ProducerIDs {
		if id == companyID {
			return
		}
	}

	anime.ProducerIDs = append(anime.ProducerIDs, companyID)
}

// AddLicensor adds the company ID to the licensor ID list if it doesn't exist already.
func (anime *Anime) AddLicensor(companyID string) {
	// Is the ID valid?
	if companyID == "" {
		return
	}

	// If it already exists we don't need to add it
	for _, id := range anime.LicensorIDs {
		if id == companyID {
			return
		}
	}

	anime.LicensorIDs = append(anime.LicensorIDs, companyID)
}

// Studios returns the list of studios for this anime.
func (anime *Anime) Studios() []*Company {
	companies := []*Company{}

	for _, obj := range DB.GetMany("Company", anime.StudioIDs) {
		if obj == nil {
			continue
		}

		companies = append(companies, obj.(*Company))
	}

	return companies
}

// Producers returns the list of producers for this anime.
func (anime *Anime) Producers() []*Company {
	companies := []*Company{}

	for _, obj := range DB.GetMany("Company", anime.ProducerIDs) {
		if obj == nil {
			continue
		}

		companies = append(companies, obj.(*Company))
	}

	return companies
}

// Licensors returns the list of licensors for this anime.
func (anime *Anime) Licensors() []*Company {
	companies := []*Company{}

	for _, obj := range DB.GetMany("Company", anime.LicensorIDs) {
		if obj == nil {
			continue
		}

		companies = append(companies, obj.(*Company))
	}

	return companies
}

// Prequels returns the list of prequels for that anime.
func (anime *Anime) Prequels() []*Anime {
	prequels := []*Anime{}
	relations := anime.Relations()

	relations.Lock()
	defer relations.Unlock()

	for _, relation := range relations.Items {
		if relation.Type != "prequel" {
			continue
		}

		prequel := relation.Anime()

		if prequel == nil {
			color.Red("Anime %s has invalid anime relation ID %s", anime.ID, relation.AnimeID)
			continue
		}

		prequels = append(prequels, prequel)
	}

	return prequels
}

// ImageLink ...
func (anime *Anime) ImageLink(size string) string {
	extension := ".jpg"

	if size == "original" {
		extension = anime.Image.Extension
	}

	return fmt.Sprintf("//%s/images/anime/%s/%s%s?%v", MediaHost, size, anime.ID, extension, anime.Image.LastModified)
}

// AverageColor returns the average color of the image.
func (anime *Anime) AverageColor() string {
	color := anime.Image.AverageColor

	if color.Hue == 0 && color.Saturation == 0 && color.Lightness == 0 {
		return ""
	}

	return color.String()
}

// Characters ...
func (anime *Anime) Characters() *AnimeCharacters {
	characters, _ := GetAnimeCharacters(anime.ID)

	if characters != nil {
		// TODO: Sort by role in sync-characters job
		// Sort by role
		sort.Slice(characters.Items, func(i, j int) bool {
			// A little trick because "main" < "supporting"
			return characters.Items[i].Role < characters.Items[j].Role
		})
	}

	return characters
}

// Relations ...
func (anime *Anime) Relations() *AnimeRelations {
	relations, _ := GetAnimeRelations(anime.ID)
	return relations
}

// Link returns the URI to the anime page.
func (anime *Anime) Link() string {
	return "/anime/" + anime.ID
}

// StartDateTime ...
func (anime *Anime) StartDateTime() time.Time {
	t, _ := time.Parse(AnimeDateFormat, anime.StartDate)
	return t
}

// Episodes returns the anime episodes wrapper.
func (anime *Anime) Episodes() *AnimeEpisodes {
	record, err := DB.Get("AnimeEpisodes", anime.ID)

	if err != nil {
		return nil
	}

	return record.(*AnimeEpisodes)
}

// UsersWatchingOrPlanned returns a list of users who are watching the anime right now.
func (anime *Anime) UsersWatchingOrPlanned() []*User {
	users := FilterUsers(func(user *User) bool {
		item := user.AnimeList().Find(anime.ID)

		if item == nil {
			return false
		}

		return item.Status == AnimeListStatusWatching || item.Status == AnimeListStatusPlanned
	})

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
		// New episodes have been released.
		// Notify all users who are watching the anime.
		go func() {
			for _, user := range anime.UsersWatchingOrPlanned() {
				if !user.Settings().Notification.AnimeEpisodeReleases {
					continue
				}

				user.SendNotification(&PushNotification{
					Title:   anime.Title.ByUser(user),
					Message: "Episode " + strconv.Itoa(newAvailableCount) + " has been released!",
					Icon:    anime.ImageLink("medium"),
					Link:    "https://notify.moe" + anime.Link(),
					Type:    NotificationTypeAnimeEpisode,
				})
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

	episodes.Save()

	return nil
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
	idList, err := GetIDList("animetwist index")

	if err != nil {
		return nil, err
	}

	// Does the index contain the ID?
	kitsuID := anime.GetMapping("kitsu/anime")
	found := false

	for _, id := range idList {
		if id == kitsuID {
			found = true
			break
		}
	}

	// If the ID is not the index we don't need to query the feed
	if !found {
		return nil, errors.New("Not available in twist.moe anime index")
	}

	// Get twist.moe feed
	feed, err := twist.GetFeedByKitsuID(kitsuID)

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
	now := time.Now().UTC().Format(time.RFC3339)

	for _, episode := range anime.Episodes().Items {
		if episode.AiringDate.Start > now && validator.IsValidDate(episode.AiringDate.Start) {
			return &UpcomingEpisode{
				Anime:   anime,
				Episode: episode,
			}
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

// ImportKitsuMapping imports the given Kitsu mapping.
func (anime *Anime) ImportKitsuMapping(mapping *kitsu.Mapping) {
	switch mapping.Attributes.ExternalSite {
	case "myanimelist/anime":
		anime.SetMapping("myanimelist/anime", mapping.Attributes.ExternalID)
	case "anidb":
		anime.SetMapping("anidb/anime", mapping.Attributes.ExternalID)
	case "trakt":
		anime.SetMapping("trakt/anime", mapping.Attributes.ExternalID)
	// case "hulu":
	// 	anime.SetMapping("hulu/anime", mapping.Attributes.ExternalID)
	case "anilist":
		externalID := mapping.Attributes.ExternalID

		if strings.HasPrefix(externalID, "anime/") {
			externalID = externalID[len("anime/"):]
		}

		anime.SetMapping("anilist/anime", externalID)
	case "thetvdb", "thetvdb/series":
		externalID := mapping.Attributes.ExternalID
		slashPos := strings.Index(externalID, "/")

		if slashPos != -1 {
			externalID = externalID[:slashPos]
		}

		anime.SetMapping("thetvdb/anime", externalID)
	case "thetvdb/season":
		// Ignore
	default:
		color.Yellow("Unknown mapping: %s %s", mapping.Attributes.ExternalSite, mapping.Attributes.ExternalID)
	}
}

// TypeHumanReadable ...
func (anime *Anime) TypeHumanReadable() string {
	switch anime.Type {
	case "tv":
		return "TV"
	case "movie":
		return "Movie"
	case "ova":
		return "OVA"
	case "ona":
		return "ONA"
	case "special":
		return "Special"
	case "music":
		return "Music"
	default:
		return anime.Type
	}
}

// StatusHumanReadable ...
func (anime *Anime) StatusHumanReadable() string {
	switch anime.Status {
	case "finished":
		return "Finished"
	case "current":
		return "Airing"
	case "upcoming":
		return "Upcoming"
	case "unannounced":
		return "Unannounced"
	case "tba":
		return "To be announced"
	default:
		return anime.Status
	}
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
	resp, err := kitsu.GetAnimeCharactersForAnime(anime.GetMapping("kitsu/anime"))

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

	animeCharacters.Save()

	return animeCharacters, nil
}

// SetID performs a database-wide ID change.
// Calling this will automatically save the anime.
func (anime *Anime) SetID(newID string) {
	oldID := anime.ID

	// Update anime ID in character list
	characters, _ := GetAnimeCharacters(oldID)
	characters.Delete()
	characters.AnimeID = newID
	characters.Save()

	// Update anime ID in relation list
	relations, _ := GetAnimeRelations(oldID)
	relations.Delete()
	relations.AnimeID = newID
	relations.Save()

	// Update anime ID in episode list
	episodes, _ := GetAnimeEpisodes(oldID)
	episodes.Delete()
	episodes.AnimeID = newID
	episodes.Save()

	// Update anime list items
	for animeList := range StreamAnimeLists() {
		item := animeList.Find(oldID)

		if item != nil {
			item.AnimeID = newID
			animeList.Save()
		}
	}

	// Update relations pointing to this anime
	for relations := range StreamAnimeRelations() {
		relation := relations.Find(oldID)

		if relation != nil {
			relation.AnimeID = newID
			relations.Save()
		}
	}

	// Update quotes
	for quote := range StreamQuotes() {
		if quote.AnimeID == oldID {
			quote.AnimeID = newID
			quote.Save()
		}
	}

	// Update log entries
	for entry := range StreamEditLogEntries() {
		switch entry.ObjectType {
		case "Anime", "AnimeRelations", "AnimeCharacters", "AnimeEpisodes":
			if entry.ObjectID == oldID {
				entry.ObjectID = newID
				entry.Save()
			}
		}
	}

	// Update ignored anime differences
	for ignore := range StreamIgnoreAnimeDifferences() {
		// ID example: arn:10052|mal:28701|RomajiTitle
		arnPart := strings.Split(ignore.ID, "|")[0]
		actualID := strings.Split(arnPart, ":")[1]

		if actualID == oldID {
			DB.Delete("IgnoreAnimeDifference", ignore.ID)
			ignore.ID = strings.Replace(ignore.ID, arnPart, "arn:"+newID, 1)
			ignore.Save()
		}
	}

	// Update soundtrack tags
	for track := range StreamSoundTracks() {
		newTags := []string{}
		modified := false

		for _, tag := range track.Tags {
			if strings.HasPrefix(tag, "anime:") {
				parts := strings.Split(tag, ":")
				id := parts[1]

				if id == oldID {
					newTags = append(newTags, "anime:"+newID)
					modified = true
					continue
				}
			}

			newTags = append(newTags, tag)
		}

		if modified {
			track.Tags = newTags
			track.Save()
		}
	}

	// Update images on file system
	if anime.Image.Extension != "" {
		err := os.Rename(
			path.Join(Root, "images/anime/original/", oldID+anime.Image.Extension),
			path.Join(Root, "images/anime/original/", newID+anime.Image.Extension),
		)

		if err != nil {
			// Don't return the error.
			// It's too late to stop the process at this point.
			// Instead, log the error.
			color.Red(err.Error())
		}

		os.Rename(
			path.Join(Root, "images/anime/large/", oldID+".jpg"),
			path.Join(Root, "images/anime/large/", newID+".jpg"),
		)

		os.Rename(
			path.Join(Root, "images/anime/large/", oldID+"@2.jpg"),
			path.Join(Root, "images/anime/large/", newID+"@2.jpg"),
		)

		os.Rename(
			path.Join(Root, "images/anime/large/", oldID+".webp"),
			path.Join(Root, "images/anime/large/", newID+".webp"),
		)

		os.Rename(
			path.Join(Root, "images/anime/large/", oldID+"@2.webp"),
			path.Join(Root, "images/anime/large/", newID+"@2.webp"),
		)

		os.Rename(
			path.Join(Root, "images/anime/medium/", oldID+".jpg"),
			path.Join(Root, "images/anime/medium/", newID+".jpg"),
		)

		os.Rename(
			path.Join(Root, "images/anime/medium/", oldID+"@2.jpg"),
			path.Join(Root, "images/anime/medium/", newID+"@2.jpg"),
		)

		os.Rename(
			path.Join(Root, "images/anime/medium/", oldID+".webp"),
			path.Join(Root, "images/anime/medium/", newID+".webp"),
		)

		os.Rename(
			path.Join(Root, "images/anime/medium/", oldID+"@2.webp"),
			path.Join(Root, "images/anime/medium/", newID+"@2.webp"),
		)

		os.Rename(
			path.Join(Root, "images/anime/small/", oldID+".jpg"),
			path.Join(Root, "images/anime/small/", newID+".jpg"),
		)

		os.Rename(
			path.Join(Root, "images/anime/small/", oldID+"@2.jpg"),
			path.Join(Root, "images/anime/small/", newID+"@2.jpg"),
		)

		os.Rename(
			path.Join(Root, "images/anime/small/", oldID+".webp"),
			path.Join(Root, "images/anime/small/", newID+".webp"),
		)

		os.Rename(
			path.Join(Root, "images/anime/small/", oldID+"@2.webp"),
			path.Join(Root, "images/anime/small/", newID+"@2.webp"),
		)
	}

	// Delete old anime ID
	DB.Delete("Anime", oldID)

	// Change anime ID and save it
	anime.ID = newID
	anime.Save()
}

// String implements the default string serialization.
func (anime *Anime) String() string {
	return anime.Title.Canonical
}

// StreamAnime returns a stream of all anime.
func StreamAnime() chan *Anime {
	channel := make(chan *Anime, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Anime") {
			channel <- obj.(*Anime)
		}

		close(channel)
	}()

	return channel
}

// AllAnime returns a slice of all anime.
func AllAnime() []*Anime {
	var all []*Anime

	stream := StreamAnime()

	for obj := range stream {
		all = append(all, obj)
	}

	return all
}

// FilterAnime filters all anime by a custom function.
func FilterAnime(filter func(*Anime) bool) []*Anime {
	var filtered []*Anime

	channel := DB.All("Anime")

	for obj := range channel {
		realObject := obj.(*Anime)

		if filter(realObject) {
			filtered = append(filtered, realObject)
		}
	}

	return filtered
}
