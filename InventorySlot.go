package arn

// InventorySlot ...
type InventorySlot struct {
	ItemID   string `json:"itemId"`
	Quantity uint   `json:"quantity"`

	item *Item
}

// IsEmpty ...
func (slot *InventorySlot) IsEmpty() bool {
	return slot.ItemID == ""
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
