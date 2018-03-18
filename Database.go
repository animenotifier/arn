package arn

import (
	"github.com/aerogo/api"
	"github.com/aerogo/nano"
	"github.com/animenotifier/jikan"
	"github.com/animenotifier/kitsu"
	"github.com/animenotifier/mal"
)

// Session ...
type Session map[string]interface{}

// Node represents the database node.
var Node = nano.New(5000)

// DB is the main database client.
var DB = Node.Namespace("arn").RegisterTypes(
	(*Analytics)(nil),
	(*Anime)(nil),
	(*AnimeCharacters)(nil),
	(*AnimeEpisodes)(nil),
	(*AnimeRelations)(nil),
	(*AnimeList)(nil),
	(*AniListToAnime)(nil),
	(*Character)(nil),
	(*Company)(nil),
	(*DraftIndex)(nil),
	(*MyAnimeListToAnime)(nil),
	(*EditLogEntry)(nil),
	(*EmailToUser)(nil),
	(*FacebookToUser)(nil),
	(*GoogleToUser)(nil),
	(*Group)(nil),
	(*GroupPost)(nil),
	(*Item)(nil),
	(*IDList)(nil),
	(*IgnoreAnimeDifference)(nil),
	(*Inventory)(nil),
	(*NickToUser)(nil),
	(*Notification)(nil),
	(*PayPalPayment)(nil),
	(*Post)(nil),
	(*Purchase)(nil),
	(*PushSubscriptions)(nil),
	(*Quote)(nil),
	(*Session)(nil),
	(*Settings)(nil),
	(*SoundTrack)(nil),
	(*Thread)(nil),
	(*TwitterToUser)(nil),
	(*User)(nil),
	(*UserFollows)(nil),
	(*UserNotifications)(nil),
)

// MAL is the client for the MyAnimeList database.
var MAL = Node.Namespace("mal").RegisterTypes(
	(*mal.Anime)(nil),
)

// Kitsu is the client for the Kitsu database.
var Kitsu = Node.Namespace("kitsu").RegisterTypes(
	(*kitsu.Anime)(nil),
	(*kitsu.Mapping)(nil),
)

// API ...
var API = api.New("/api/", DB)

// init ...
func init() {
	Node.Namespace("jikan").RegisterTypes(
		(*jikan.Anime)(nil),
		(*jikan.Character)(nil),
	)
}
