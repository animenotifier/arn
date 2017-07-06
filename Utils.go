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
	shortid "github.com/ventu-io/go-shortid"
	"github.com/xrash/smetrics"
)

var stripTagsRegex = regexp.MustCompile(`<[^>]*>`)
var sourceRegex = regexp.MustCompile(`\(Source: (.*?)\)`)
var writtenByRegex = regexp.MustCompile(`\[Written by (.*?)\]`)

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

// SetObjectProperties updates the object with the given map[string]interface{}
func SetObjectProperties(rootObj interface{}, updates map[string]interface{}, skip func(fullKeyName string, field *reflect.StructField, property *reflect.Value, newValue reflect.Value) bool) error {
	var t reflect.Type
	var v reflect.Value
	var field reflect.StructField
	var found bool

	for key, value := range updates {
		t = reflect.TypeOf(rootObj).Elem()
		v = reflect.ValueOf(rootObj).Elem()

		if strings.HasPrefix(key, "Custom:") {
			skip(key, nil, nil, reflect.ValueOf(value))
			continue
		}

		// Nested properties
		parts := strings.Split(key, ".")

		for _, part := range parts {
			field, found = t.FieldByName(part)

			if !found {
				return errors.New("Field '" + part + "' does not exist in type " + t.Name())
			}

			t = field.Type
			v = reflect.Indirect(v.FieldByName(field.Name))

			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
		}

		newValue := reflect.ValueOf(value)

		// Is somebody attempting to edit fields that aren't editable?
		if field.Tag.Get("editable") != "true" {
			return errors.New("Field " + key + " is not editable")
		}

		// Is this manually handled by the class so we can skip it?
		// Also make sure to pass full "key" value here instead of "fieldName".
		if skip != nil && skip(key, &field, &v, newValue) {
			continue
		}

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

// OverallRatingName returns Overall in general, but Hype when episodes watched is zero.
func OverallRatingName(episodes int) string {
	if episodes == 0 {
		return "Hype"
	}

	return "Overall"
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
