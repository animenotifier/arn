package autocorrect

import (
	"regexp"
	"strings"
)

const maxNickLength = 25

var fixNickRegex = regexp.MustCompile(`[\W\s\d]`)

var accountNickRegexes = []*regexp.Regexp{
	regexp.MustCompile(`anilist.co/user/(.*)`),
	regexp.MustCompile(`anilist.co/animelist/(.*)`),
	regexp.MustCompile(`kitsu.io/users/(.*?)/library`),
	regexp.MustCompile(`kitsu.io/users/(.*)`),
	regexp.MustCompile(`anime-planet.com/users/(.*?)/anime`),
	regexp.MustCompile(`anime-planet.com/users/(.*)`),
	regexp.MustCompile(`myanimelist.net/profile/(.*)`),
	regexp.MustCompile(`myanimelist.net/animelist/(.*?)\?`),
	regexp.MustCompile(`myanimelist.net/animelist/(.*)`),
	regexp.MustCompile(`myanimelist.net/(.*)`),
	regexp.MustCompile(`myanimelist.com/(.*)`),
	regexp.MustCompile(`twitter.com/(.*)`),
	regexp.MustCompile(`osu.ppy.sh/u/(.*)`),
}

var animeLinkRegex = regexp.MustCompile(`notify.moe/anime/(\d+)`)
var osuBeatmapRegex = regexp.MustCompile(`osu.ppy.sh/s/(\d+)`)

// FixTag converts links to correct tags automatically.
func FixTag(tag string) string {
	tag = strings.TrimSpace(tag)
	tag = strings.TrimSuffix(tag, "/")

	// Anime
	matches := animeLinkRegex.FindStringSubmatch(tag)

	if len(matches) > 1 {
		return "anime:" + matches[1]
	}

	// Osu beatmap
	matches = osuBeatmapRegex.FindStringSubmatch(tag)

	if len(matches) > 1 {
		return "osu-beatmap:" + matches[1]
	}

	return tag
}

// FixUserNick automatically corrects a username.
func FixUserNick(nick string) string {
	nick = fixNickRegex.ReplaceAllString(nick, "")

	if nick == "" {
		return nick
	}

	nick = strings.Trim(nick, "_")

	if nick == "" {
		return ""
	}

	if len(nick) > maxNickLength {
		nick = nick[:maxNickLength]
	}

	return strings.ToUpper(string(nick[0])) + nick[1:]
}

// FixAccountNick automatically corrects the username/nick of an account.
func FixAccountNick(nick string) string {
	for _, regex := range accountNickRegexes {
		matches := regex.FindStringSubmatch(nick)

		if len(matches) > 1 {
			nick = matches[1]
			return nick
		}
	}

	return nick
}

// FixPostText fixes common mistakes in post texts.
func FixPostText(text string) string {
	text = strings.Replace(text, "http://", "https://", -1)
	text = strings.TrimSpace(text)
	return text
}

// FixThreadTitle ...
func FixThreadTitle(title string) string {
	return strings.TrimSpace(title)
}
