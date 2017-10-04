package arn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/animenotifier/arn"
)

func TestInventory(t *testing.T) {
	inventory := arn.NewInventory("4J6qpK1ve")
	assert.Len(t, inventory.Slots, arn.DefaultInventorySlotCount)
}
