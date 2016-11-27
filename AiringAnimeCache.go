package arn

// AiringAnimeCache ...
type AiringAnimeCache struct {
	Anime []*Anime `json:"anime"`
}

// AiringAnimeCacheSmall ...
type AiringAnimeCacheSmall struct {
	Anime []*AnimeSmall `json:"anime"`
}
