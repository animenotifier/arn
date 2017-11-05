package arn

// AnimeCharacter ...
type AnimeCharacter struct {
	CharacterID string `json:"characterId"`
	Role        string `json:"role"`

	character *Character
}
