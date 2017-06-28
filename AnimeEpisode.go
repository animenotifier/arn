package arn

// AnimeEpisode ...
type AnimeEpisode struct {
	Number     int              `json:"number"`
	Title      *EpisodeTitle    `json:"title"`
	AiringDate *AnimeAiringDate `json:"airingDate"`
}

// EpisodeTitle ...
type EpisodeTitle struct {
	Romaji   string `json:"romaji"`
	English  string `json:"english"`
	Japanese string `json:"japanese"`
}
