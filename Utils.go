package arn

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aerogo/aero"
	"github.com/aerogo/mirror"
	"github.com/animenotifier/kitsu"
	"github.com/animenotifier/mal"
	shortid "github.com/ventu-io/go-shortid"
	"github.com/xrash/smetrics"
)

var (
	// MediaHost is the host we use to link image files.
	MediaHost = "media.notify.moe"

	// Regular expressions
	stripTagsRegex = regexp.MustCompile(`<[^>]*>`)
	sourceRegex    = regexp.MustCompile(`\(Source: (.*?)\)`)
	writtenByRegex = regexp.MustCompile(`\[Written by (.*?)\]`)
)

// GenerateID generates a unique ID for a given table.
func GenerateID(table string) string {
	id, _ := shortid.Generate()

	// Retry until we find an unused ID
	retry := 0

	for {
		_, err := DB.Get(table, id)

		if err != nil && strings.Index(err.Error(), "not found") != -1 {
			return id
		}

		retry++

		if retry > 10 {
			panic(errors.New("Can't generate unique ID"))
		}

		id, _ = shortid.Generate()
	}
}

// GetUserFromContext returns the logged in user for the given context.
func GetUserFromContext(ctx *aero.Context) *User {
	if !ctx.HasSession() {
		return nil
	}

	userID := ctx.Session().GetString("userId")

	if userID == "" {
		return nil
	}

	user, err := GetUser(userID)

	if err != nil {
		return nil
	}

	return user
}

// GetObjectTitle ...
func GetObjectTitle(typeName string, id string) string {
	obj, err := DB.Get(typeName, id)

	if err != nil {
		return fmt.Sprintf("<not found: %s>", id)
	}

	return fmt.Sprint(obj)
}

// FilterIDTags returns all IDs of the given type in the tag list.
func FilterIDTags(tags []string, idType string) []string {
	var idList []string
	prefix := idType + ":"

	for _, tag := range tags {
		if strings.HasPrefix(tag, prefix) {
			id := strings.TrimPrefix(tag, prefix)
			idList = append(idList, id)
		}
	}

	return idList
}

// JSON turns the object into a JSON string.
func JSON(obj interface{}) string {
	data, err := json.Marshal(obj)

	if err == nil {
		return string(data)
	}

	return err.Error()
}

// SetObjectProperties updates the object with the given map[string]interface{}
func SetObjectProperties(rootObj interface{}, updates map[string]interface{}) error {
	for key, value := range updates {
		field, _, v, err := mirror.GetField(rootObj, key)

		if err != nil {
			return err
		}

		// Is somebody attempting to edit fields that aren't editable?
		if field.Tag.Get("editable") != "true" {
			return errors.New("Field " + key + " is not editable")
		}

		newValue := reflect.ValueOf(value)

		// Implement special data type cases here
		if v.Kind() == reflect.Int {
			x := int64(newValue.Float())

			if !v.OverflowInt(x) {
				v.SetInt(x)
			}
		} else {
			v.Set(newValue)
		}
	}

	return nil
}

// GetGenreIDByName ...
func GetGenreIDByName(genre string) string {
	genre = strings.Replace(genre, "-", "", -1)
	genre = strings.Replace(genre, " ", "", -1)
	genre = strings.ToLower(genre)
	return genre
}

// FixAnimeDescription ...
func FixAnimeDescription(description string) string {
	description = stripTagsRegex.ReplaceAllString(description, "")
	description = sourceRegex.ReplaceAllString(description, "")
	description = writtenByRegex.ReplaceAllString(description, "")
	return strings.TrimSpace(description)
}

// FixGender ...
func FixGender(gender string) string {
	if gender != "male" && gender != "female" {
		return ""
	}

	return gender
}

// AnimeRatingStars displays the rating in Unicode stars.
func AnimeRatingStars(rating float64) string {
	stars := int(rating/20 + 0.5)
	return strings.Repeat("★", stars) + strings.Repeat("☆", 5-stars)
}

// EpisodesToString shows a question mark if the episode count is zero.
func EpisodesToString(episodes int) string {
	if episodes == 0 {
		return "?"
	}

	return ToString(episodes)
}

// EpisodeCountMax is used for the max value of number input on episodes.
func EpisodeCountMax(episodes int) string {
	if episodes == 0 {
		return ""
	}

	return strconv.Itoa(episodes)
}

// DateTimeUTC returns the current UTC time in RFC3339 format.
func DateTimeUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// StringSimilarity returns 1.0 if the strings are equal and goes closer to 0 when they are different.
func StringSimilarity(a string, b string) float64 {
	return smetrics.JaroWinkler(a, b, 0.7, 4)
}

// AdvancedStringSimilarity is like StringSimilarity but boosts the value if a appears directly in b.
func AdvancedStringSimilarity(a string, b string) float64 {
	s := StringSimilarity(a, b)

	if strings.Contains(b, a) {
		s += 0.5

		if strings.HasPrefix(b, a) {
			s += 5.0
		}
	}

	return s
}

// OverallRatingName returns Overall in general, but Hype when episodes watched is zero.
func OverallRatingName(episodes int) string {
	if episodes == 0 {
		return "Hype"
	}

	return "Overall"
}

// IsIPv6 tells you whether the given address is IPv6 encoded.
func IsIPv6(ip string) bool {
	for i := 0; i < len(ip); i++ {
		if ip[i] == ':' {
			return true
		}
	}

	return false
}

// MyAnimeListStatusToARNStatus ...
func MyAnimeListStatusToARNStatus(status string) string {
	switch status {
	case mal.AnimeListStatusCompleted:
		return AnimeListStatusCompleted
	case mal.AnimeListStatusWatching:
		return AnimeListStatusWatching
	case mal.AnimeListStatusPlanned:
		return AnimeListStatusPlanned
	case mal.AnimeListStatusHold:
		return AnimeListStatusHold
	case mal.AnimeListStatusDropped:
		return AnimeListStatusDropped
	default:
		return ""
	}
}

// KitsuStatusToARNStatus ...
func KitsuStatusToARNStatus(status string) string {
	switch status {
	case kitsu.AnimeListStatusCompleted:
		return AnimeListStatusCompleted
	case kitsu.AnimeListStatusWatching:
		return AnimeListStatusWatching
	case kitsu.AnimeListStatusPlanned:
		return AnimeListStatusPlanned
	case kitsu.AnimeListStatusHold:
		return AnimeListStatusHold
	case kitsu.AnimeListStatusDropped:
		return AnimeListStatusDropped
	default:
		return ""
	}
}

// ListItemStatusName ...
func ListItemStatusName(status string) string {
	switch status {
	case AnimeListStatusWatching:
		return "Watching"
	case AnimeListStatusCompleted:
		return "Completed"
	case AnimeListStatusPlanned:
		return "Planned"
	case AnimeListStatusHold:
		return "On hold"
	case AnimeListStatusDropped:
		return "Dropped"
	default:
		return ""
	}
}

// PanicOnError will panic if the error is not nil.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// PrettyPrint prints the object as indented JSON data on the console.
func PrettyPrint(obj interface{}) {
	pretty, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Println(string(pretty))
}
