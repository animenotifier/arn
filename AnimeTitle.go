package arn

// AnimeTitle ...
type AnimeTitle struct {
	Romaji    string   `json:"romaji"`
	English   string   `json:"english"`
	Japanese  string   `json:"japanese"`
	Hiragana  string   `json:"hiragana"`
	Canonical string   `json:"canonical"`
	Synonyms  []string `json:"synonyms"`
}
