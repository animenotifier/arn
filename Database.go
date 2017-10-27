package arn

import (
	"github.com/aerogo/api"
	"github.com/aerogo/database"
)

// DB is the main database client.
var DB = database.New("arn", DBTypes)

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
