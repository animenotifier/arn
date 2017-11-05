package arn

// Save saves the item in the database.
func (item *Item) Save() {
	DB.Set("Item", item.ID, item)
}
