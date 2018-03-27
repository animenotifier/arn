package arn

// CharacterName ...
type CharacterName struct {
	Canonical string   `json:"canonical"`
	English   string   `json:"english"`
	Japanese  string   `json:"japanese"`
	Synonyms  []string `json:"synonyms"`
}
