package arn

import (
	"github.com/aerogo/api"
	"github.com/aerogo/nano"
	"github.com/animenotifier/kitsu"
	"github.com/animenotifier/mal"
)

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
	(*Character)(nil),
	(*Company)(nil),
	(*DraftIndex)(nil),
	(*EditLogEntry)(nil),
	(*EmailToUser)(nil),
	(*FacebookToUser)(nil),
	(*GoogleToUser)(nil),
	(*Group)(nil),
	(*GroupPost)(nil),
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
	(*ShopItem)(nil),
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
	(*mal.Character)(nil),
)

// Kitsu is the client for the Kitsu database.
var Kitsu = Node.Namespace("kitsu").RegisterTypes(
	(*kitsu.Anime)(nil),
	(*kitsu.Mapping)(nil),
	(*kitsu.Character)(nil),
)

// API ...
var API = api.New("/api/", DB)
