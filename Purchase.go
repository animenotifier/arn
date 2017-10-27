package arn

import "github.com/aerogo/database"

// Purchase ...
type Purchase struct {
	ID       string `json:"id"`
	UserID   string `json:"userId"`
	ItemID   string `json:"itemId"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Currency string `json:"currency"`
	Date     string `json:"date"`
}

// Item returns the item the user bought.
func (purchase *Purchase) Item() *Item {
	item, _ := GetItem(purchase.ItemID)
	return item
}

// User returns the user who made the purchase.
func (purchase *Purchase) User() *User {
	user, _ := GetUser(purchase.UserID)
	return user
}

// NewPurchase creates a new Purchase object with a generated ID.
func NewPurchase(userID string, itemID string, quantity int, price int, currency string) *Purchase {
	return &Purchase{
		ID:       GenerateID("Purchase"),
		UserID:   userID,
		ItemID:   itemID,
		Quantity: quantity,
		Price:    price,
		Currency: currency,
		Date:     DateTimeUTC(),
	}
}

// StreamPurchases returns a stream of all purchases.
func StreamPurchases() chan *Purchase {
	channel := make(chan *Purchase, database.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Purchase") {
			channel <- obj.(*Purchase)
		}

		close(channel)
	}()

	return channel
}

// AllPurchases returns a slice of all anime.
func AllPurchases() ([]*Purchase, error) {
	var all []*Purchase

	for obj := range StreamPurchases() {
		all = append(all, obj)
	}

	return all, nil
}

// FilterPurchases filters all purchases by a custom function.
func FilterPurchases(filter func(*Purchase) bool) ([]*Purchase, error) {
	var filtered []*Purchase

	for obj := range StreamPurchases() {
		if filter(obj) {
			filtered = append(filtered, obj)
		}
	}

	return filtered, nil
}
