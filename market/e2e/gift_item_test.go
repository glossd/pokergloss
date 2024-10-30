package e2e

import (
	"context"
	"github.com/glossd/pokergloss/market/db"
	"github.com/glossd/pokergloss/market/domain"
	"github.com/glossd/pokergloss/market/service"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGiftItem(t *testing.T) {
	t.Cleanup(cleanUp)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	assert.Nil(t, service.GiftItem(ctx, "1", domain.BurgerID, 1, domain.Day))

	inv, err := db.FindInventoryNoCtx("1")
	assert.Nil(t, err)
	inv.SelectedItemID = "burger"
}
