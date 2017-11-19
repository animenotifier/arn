package arn

// Save saves the character in the database.
func (chars *AnimeCharacters) Save() {
	DB.Set("AnimeCharacters", chars.AnimeID, chars)
}
