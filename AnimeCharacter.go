package arn

// AnimeCharacter ...
type AnimeCharacter struct {
	CharacterID string `json:"characterId"`
	Role        string `json:"role"`
}

// Character ...
func (char *AnimeCharacter) Character() *Character {
	character, _ := GetCharacter(char.CharacterID)
	return character
}
