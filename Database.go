package arn

import "github.com/aerogo/aerospike"

// DB is the main database client.
var DB = aerospike.NewDatabase(
	"arn-db",
	3000,
	"arn",
	[]interface{}{
		(*Anime)(nil),
		(*AnimeList)(nil),
		(*Post)(nil),
		(*Settings)(nil),
		(*Thread)(nil),
		(*User)(nil),
		(*NickToUser)(nil),
		(*EmailToUser)(nil),
		(*GoogleToUser)(nil),
		(*SearchIndex)(nil),
	},
)
