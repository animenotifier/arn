package arn

// Save saves the character in the database.
func (char *Character) Save() {
	DB.Set("Character", char.ID, char)
}
