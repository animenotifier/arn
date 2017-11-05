package arn

// AnimeCharacters ...
type AnimeCharacters struct {
	AnimeID string            `json:"animeId"`
	Items   []*AnimeCharacter `json:"items"`
}

// Character ...
func (char *AnimeCharacter) Character() *Character {
	character, _ := GetCharacter(char.CharacterID)
	return character
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
