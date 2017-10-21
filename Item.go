package arn

const (
	// ItemRarityCommon ...
	ItemRarityCommon = "common"

	// ItemRaritySuperior ...
	ItemRaritySuperior = "superior"

	// ItemRarityRare ...
	ItemRarityRare = "rare"

	// ItemRarityUnique ...
	ItemRarityUnique = "unique"

	// ItemRarityLegendary ...
	ItemRarityLegendary = "legendary"
)

// Item ...
type Item struct {
	ID          string     `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	Icon        string `json:"icon"`
	Rarity      string `json:"rarity"`
	Order       int    `json:"order"`
	Consumable  bool   `json:"consumable"`
}

// GetItem ...
func GetItem(id string) (*Item, error) {
	obj, err := DB.Get("Item", id)

	if err != nil {
		return nil, err
	}

	return obj.(*Item), nil
}

// StreamItems returns a stream of all items.
func StreamItems() (chan *Item, error) {
	objects, err := DB.All("Item")
	return objects.(chan *Item), err
}

// MustStreamItems returns a stream of all items.
func MustStreamItems() chan *Item {
	stream, err := StreamItems()
	PanicOnError(err)
	return stream
}

// AllItems returns a slice of all items.
func AllItems() ([]*Item, error) {
	var all []*Item

	stream, err := StreamItems()

	if err != nil {
		return nil, err
	}

	for obj := range stream {
		all = append(all, obj)
	}

	return all, nil
}
