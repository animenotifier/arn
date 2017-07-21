package arn

import "github.com/aerogo/aerospike"

// DB is the main database client.
var DB = aerospike.NewDatabase(
	"arn-db",
	3000,
	"arn",
	[]interface{}{
		(*Analytics)(nil),
		(*Anime)(nil),
		(*AnimeCharacters)(nil),
		(*AnimeEpisodes)(nil),
		(*AnimeList)(nil),
		(*AniListToAnime)(nil),
		(*Character)(nil),
		(*MyAnimeListToAnime)(nil),
		(*EmailToUser)(nil),
		(*FacebookToUser)(nil),
		(*GoogleToUser)(nil),
		(*NickToUser)(nil),
		(*PayPalPayment)(nil),
		(*Post)(nil),
		(*PushSubscriptions)(nil),
		(*SearchIndex)(nil),
		(*Settings)(nil),
		(*SoundCloudToSoundTrack)(nil),
		(*SoundTrack)(nil),
		(*Thread)(nil),
		(*TwitterToUser)(nil),
		(*User)(nil),
		(*UserFollows)(nil),
		(*YoutubeToSoundTrack)(nil),
	},
)
