package arn

// AnimeEpisode ...
type AnimeEpisode struct {
	Number     int               `json:"number"`
	Title      *EpisodeTitle     `json:"title"`
	AiringDate *AnimeAiringDate  `json:"airingDate"`
	Links      map[string]string `json:"links"`
}

// EpisodeTitle ...
type EpisodeTitle struct {
	Romaji   string `json:"romaji"`
	English  string `json:"english"`
	Japanese string `json:"japanese"`
}
