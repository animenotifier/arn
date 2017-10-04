package arn

// InventorySlot ...
type InventorySlot struct {
	ItemID   string `json:"itemId"`
	Quantity uint   `json:"quantity"`

	item *Item
}

// Item ...
func (slot *InventorySlot) Item() *Item {
	if slot.item != nil {
		return slot.item
	}

	if slot.ItemID == "" {
		return nil
	}

	slot.item, _ = GetItem(slot.ItemID)
	return slot.item
}
