package arn

import "github.com/aerogo/aerospike"

// DB is the main database client.
var DB = aerospike.NewDatabase(
	"arn-db",
	3000,
	"arn",
	[]interface{}{
		new(Anime),
		new(AnimeList),
		new(Post),
		new(Settings),
		new(Thread),
		new(User),
	},
)
