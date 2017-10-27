package arn

// AnimeCharacters ...
type AnimeCharacters struct {
	AnimeID string            `json:"animeId"`
	Items   []*AnimeCharacter `json:"items"`
}

// AnimeCharacter ...
type AnimeCharacter struct {
	CharacterID string `json:"characterId"`
	Role        string `json:"role"`

	character *Character
}

// Character ...
func (char *AnimeCharacter) Character() *Character {
	if char.character != nil {
		return char.character
	}

	char.character, _ = GetCharacter(char.CharacterID)

	return char.character
}

// Save saves the character in the database.
func (chars *AnimeCharacters) Save() {
	DB.Set("AnimeCharacters", chars.AnimeID, chars)
}

// GetAnimeCharacters ...
func GetAnimeCharacters(animeID string) (*AnimeCharacters, error) {
	obj, err := DB.Get("AnimeCharacters", animeID)

	if err != nil {
		return nil, err
	}

	return obj.(*AnimeCharacters), nil
}
