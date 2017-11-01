package arn

import (
	"github.com/aerogo/api"
	"github.com/aerogo/nano"
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
	(*DraftIndex)(nil),
	(*MyAnimeListToAnime)(nil),
	(*EmailToUser)(nil),
	(*FacebookToUser)(nil),
	(*GoogleToUser)(nil),
	(*Group)(nil),
	(*Item)(nil),
	(*IDList)(nil),
	(*Inventory)(nil),
	(*NickToUser)(nil),
	(*PayPalPayment)(nil),
	(*Post)(nil),
	(*Purchase)(nil),
	(*PushSubscriptions)(nil),
	(*SearchIndex)(nil),
	(*Session)(nil),
	(*Settings)(nil),
	(*SoundTrack)(nil),
	(*Thread)(nil),
	(*TwitterToUser)(nil),
	(*User)(nil),
	(*UserFollows)(nil),
)

// API ...
var API = api.New("/api/", DB)
