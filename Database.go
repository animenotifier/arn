package arn

import (
	"github.com/aerogo/aerospike"
	"github.com/aerogo/api"
)

// DB is the main database client.
var DB = aerospike.NewDatabase(
	"arn-db",
	3000,
	"arn",
	DBTypes,
)

// DBTypes ...
var DBTypes = []interface{}{
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
	(*Inventory)(nil),
	(*NickToUser)(nil),
	(*PayPalPayment)(nil),
	(*Post)(nil),
	(*Purchase)(nil),
	(*PushSubscriptions)(nil),
	(*SearchIndex)(nil),
	(*Settings)(nil),
	(*SoundTrack)(nil),
	(*Thread)(nil),
	(*TwitterToUser)(nil),
	(*User)(nil),
	(*UserFollows)(nil),
}

// API ...
var API = api.New("/api/", DB)
