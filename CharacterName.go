package arn

// CharacterName ...
type CharacterName struct {
	Canonical string   `json:"canonical" editable:"true"`
	English   string   `json:"english" editable:"true"`
	Japanese  string   `json:"japanese" editable:"true"`
	Synonyms  []string `json:"synonyms" editable:"true"`
}
