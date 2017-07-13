package arn

// Save saves the character in the database.
func (char *Character) Save() error {
	return DB.Set("Character", char.ID, char)
}
