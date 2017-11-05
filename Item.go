package arn

import "github.com/aerogo/nano"

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
	ID          string `json:"id"`
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
func StreamItems() chan *Item {
	channel := make(chan *Item, nano.ChannelBufferSize)

	go func() {
		for obj := range DB.All("Item") {
			channel <- obj.(*Item)
		}

		close(channel)
	}()

	return channel
}

// AllItems returns a slice of all items.
func AllItems() ([]*Item, error) {
	var all []*Item

	for obj := range StreamItems() {
		all = append(all, obj)
	}

	return all, nil
}
