package arn

// Save saves the index in the database.
func (index *DraftIndex) Save() error {
	return DB.Set("DraftIndex", index.UserID, index)
}
