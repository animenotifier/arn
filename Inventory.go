package arn

import "errors"

// DefaultInventorySlotCount tells you how many slots are available by default in an inventory.
const DefaultInventorySlotCount = 24

// Inventory ...
type Inventory struct {
	UserID string           `json:"userId"`
	Slots  []*InventorySlot `json:"slots"`
}

// AddItem adds a given item to the inventory.
func (inventory *Inventory) AddItem(itemID string, quantity uint) error {
	if itemID == "" {
		return nil
	}

	// Find the slot with the item
	for _, slot := range inventory.Slots {
		if slot.ItemID == itemID {
			slot.Quantity += quantity
			return nil
		}
	}

	// If the item doesn't exist in the inventory yet, add it to the first free slot
	for _, slot := range inventory.Slots {
		if slot.ItemID == "" {
			slot.ItemID = itemID
			slot.Quantity = quantity
			return nil
		}
	}

	// If there is no free slot, return an error
	return errors.New("Inventory is full")
}

// NewInventory creates a new inventory with the default number of slots.
func NewInventory(userID string) *Inventory {
	return &Inventory{
		UserID: userID,
		Slots:  make([]*InventorySlot, DefaultInventorySlotCount, DefaultInventorySlotCount),
	}
}

// GetInventory ...
func GetInventory(userID string) (*Inventory, error) {
	obj, err := DB.Get("Inventory", userID)

	if err != nil {
		return nil, err
	}

	return obj.(*Inventory), nil
}
