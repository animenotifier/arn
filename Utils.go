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
	"unicode"
	"unicode/utf8"

	"github.com/aerogo/aero"
	shortid "github.com/ventu-io/go-shortid"
)

var stripTagsRegex = regexp.MustCompile(`<[^>]*>`)
var sourceRegex = regexp.MustCompile(`\(Source: (.*?)\)`)
var writtenByRegex = regexp.MustCompile(`\[Written by (.*?)\]`)

// GenerateUserID generates a unique user ID.
func GenerateUserID() string {
	id, _ := shortid.Generate()

	// Retry until we find an unused ID
	retry := 0

	for {
		_, err := GetUser(id)

		if err != nil && strings.Index(err.Error(), "not found") != -1 {
			return id
		}

		retry++

		if retry > 10 {
			panic(errors.New("Can't generate unique user ID"))
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
func SetObjectProperties(item interface{}, updates map[string]interface{}) error {
	t := reflect.TypeOf(item).Elem()
	v := reflect.ValueOf(item).Elem()

	for key, value := range updates {
		_, found := t.FieldByName(key)

		if !found {
			return errors.New("Field '" + key + "' does not exist in type " + t.Name())
		}

		valueInfo := reflect.ValueOf(value)
		fieldValue := v.FieldByName(key)

		if fieldValue.Kind() == reflect.Int {
			x := int64(valueInfo.Float())

			if !fieldValue.OverflowInt(x) {
				fieldValue.SetInt(x)
			}
		} else {
			fieldValue.Set(valueInfo)
		}
	}

	return nil
}

// AuthorizeIfLoggedInAndOwnData authorizes the given request if a user is logged in
// and the user ID matches the ID in the request.
func AuthorizeIfLoggedInAndOwnData(ctx *aero.Context, userIDParameterName string) error {
	if !ctx.HasSession() {
		return errors.New("Neither logged in nor in session")
	}

	userID, ok := ctx.Session().Get("userId").(string)

	if !ok || userID == "" {
		return errors.New("Not logged in")
	}

	if userID != ctx.Get(userIDParameterName) {
		return errors.New("Can not modify data from other users")
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

// Capitalize returns the string with the first letter capitalized.
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
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

// ToString converts anything into a string.
func ToString(v interface{}) string {
	return fmt.Sprint(v)
}

// Plural returns the number concatenated to the proper pluralization of the word.
func Plural(count int, singular string) string {
	if count == 1 || count == -1 {
		return ToString(count) + " " + singular
	}

	return ToString(count) + " " + singular + "s"
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
