package arn

// Save saves the anime relations object in the database.
func (relations *AnimeRelations) Save() {
	DB.Set("AnimeRelations", relations.AnimeID, relations)
}
