package arn

import "errors"

// InventorySlot ...
type InventorySlot struct {
	ItemID   string `json:"itemId"`
	Quantity uint   `json:"quantity"`

	item *ShopItem
}

// IsEmpty ...
func (slot *InventorySlot) IsEmpty() bool {
	return slot.ItemID == ""
}

// Item ...
func (slot *InventorySlot) Item() *ShopItem {
	if slot.item != nil {
		return slot.item
	}

	if slot.ItemID == "" {
		return nil
	}

	slot.item, _ = GetShopItem(slot.ItemID)
	return slot.item
}

// Decrease reduces the quantity by the given number.
func (slot *InventorySlot) Decrease(count uint) error {
	if slot.Quantity < count {
		return errors.New("Not enough items")
	}

	slot.Quantity -= count

	if slot.Quantity == 0 {
		slot.ItemID = ""
	}

	return nil
}

// Increase increases the quantity by the given number.
func (slot *InventorySlot) Increase(count uint) {
	slot.Quantity += count
}
