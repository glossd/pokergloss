package e2e

import (
	"context"
	"github.com/glossd/pokergloss/market/db"
	"github.com/glossd/pokergloss/market/domain"
	"github.com/glossd/pokergloss/market/service"
	"github.com/glossd/pokergloss/market/web/rest/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReselect(t *testing.T) {
	t.Cleanup(cleanUp)
	ctx := context.Background()
	err := service.BuyItem(ctx, defaultIdentity, model.BuyItemParams{
		ItemID:    string(domain.Ship.ID),
		Units:     1,
		TimeFrame: domain.Day,
	})
	assert.Nil(t, err)

	domain.Now = func() time.Time {
		return time.Now().Add(25 * time.Hour)
	}

	item, err := service.GetSelectedItem(ctx, defaultIdentity.UserId)
	assert.Nil(t, err)
	assert.EqualValues(t, domain.Invisible.ID, item.ItemID)
	inventory, err := db.FindInventory(ctx, defaultIdentity.UserId)
	assert.Nil(t, err)
	assert.EqualValues(t, domain.Invisible.ID, inventory.SelectedItemID)
	assert.EqualValues(t, 1, len(inventory.Inventory))
}
