package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const defaultUserID = "userId"

func TestUser_BuyItem_ShouldIncreaseExpiration(t *testing.T) {
	now := time.Now()
	u := NewInventory(defaultUserID)
	cmd, err := NewPurchaseItemCommand(defaultUserID, GlassOfWineID, 1, Day)
	assert.Nil(t, err)
	u.AddItem(cmd)
	assert.EqualValues(t, 2, len(u.AvailableItems()))
	assert.LessOrEqual(t, now.AddDate(0, 0, 1).Unix(), u.GetSelectedItemUnsafe().ExpiresAt)

	cmd, err = NewPurchaseItemCommand(defaultUserID, GlassOfWineID, 1, Week)
	assert.Nil(t, err)
	u.AddItem(cmd)
	assert.LessOrEqual(t, now.AddDate(0, 0, 8).Unix(), u.GetSelectedItemUnsafe().ExpiresAt)
}

func TestInventory_Reselect(t *testing.T) {
	t.Cleanup(func() {
		Now = time.Now
	})
	inv := NewInventory(defaultUserID)
	assert.EqualValues(t, Invisible.ID, inv.SelectedItemID)

	gift, err := NewGiftItemCommand(defaultUserID, HellAmulet.ID, 3, Day)
	inv.AddItem(gift)
	assert.EqualValues(t, HellAmulet.ID, inv.SelectedItemID)

	purchase, err := NewPurchaseItemCommand(defaultUserID, Ship.ID, 1, Day)
	assert.Nil(t, err)
	inv.AddItem(purchase)
	assert.EqualValues(t, Ship.ID, inv.SelectedItemID)
	assert.EqualValues(t, 3, len(inv.Inventory))

	Now = func() time.Time {
		return time.Now().Add(25*time.Hour)
	}
	assert.True(t, inv.Reselect())
	assert.EqualValues(t, HellAmulet.ID, inv.SelectedItemID)
	assert.EqualValues(t, 2, len(inv.Inventory))
}
