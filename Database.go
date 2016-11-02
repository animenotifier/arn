package arn

import as "github.com/aerospike/aerospike-client-go"

// Client ...
var client *as.Client

// GetAnime ...
func GetAnime(id int) (*Anime, error) {
	key, _ := as.NewKey("arn", "Anime", id)
	anime := new(Anime)
	err := client.GetObject(nil, key, anime)
	return anime, err
}

func init() {
	as.SetAerospikeTag("json")

	var err error
	client, err = as.NewClient("127.0.0.1", 3000)

	if err != nil {
		panic(err)
	}
}
