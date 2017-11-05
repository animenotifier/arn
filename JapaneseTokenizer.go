package arn

import "github.com/animenotifier/japanese/client"

// JapaneseTokenizer tokenizes a sentence via the HTTP API.
var JapaneseTokenizer = &client.Tokenizer{
	Endpoint: "http://arn-jp:1234/",
}
