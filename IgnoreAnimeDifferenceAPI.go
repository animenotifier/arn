package arn

// Save saves the object in the database.
func (ignore *IgnoreAnimeDifference) Save() {
	DB.Set("IgnoreAnimeDifference", ignore.ID, ignore)
}
