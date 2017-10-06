package arn

// Save saves the purchase in the database.
func (purchase *Purchase) Save() error {
	return DB.Set("Purchase", purchase.ID, purchase)
}
