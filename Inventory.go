package arn

const defaultSlotCount = 25

// Inventory ...
type Inventory struct {
	UserID string          `json:"userId"`
	Slots  []InventorySlot `json:"slots"`
}

// NewInventory creates a new inventory with the default number of slots.
func NewInventory(userID string) *Inventory {
	return &Inventory{
		UserID: userID,
		Slots:  make([]InventorySlot, defaultSlotCount, defaultSlotCount),
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
