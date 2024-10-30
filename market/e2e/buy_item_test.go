package e2e

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/market/db"
	"github.com/glossd/pokergloss/market/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestBuyItem(t *testing.T) {
	t.Cleanup(cleanUp)
	buyGlassOfWine(t)
}

func TestBuyItem_NotForSale(t *testing.T) {
	buyItemStatus(t, domain.CrownID, domain.Day, http.StatusBadRequest)
}

func TestBuyItem_NotExisting(t *testing.T) {
	buyItemStatus(t, "blah", domain.Day, http.StatusBadRequest)
}

func TestBuyItem_Month(t *testing.T) {
	ctx := context.Background()
	buyItemMonth(t, domain.MartiniID)
	inv, err := db.FindInventory(ctx, defaultIdentity.UserId)
	assert.Nil(t, err)
	martini := inv.Inventory[domain.MartiniID]
	expAt := time.Unix(martini.ExpiresAt, 0)
	now := time.Now()
	assert.EqualValues(t, expAt.Day(), now.AddDate(0, 0, 30).Day())
}

func buyGlassOfWine(t *testing.T) {
	buyItem(t, domain.GlassOfWineID)
}

func buyItem(t *testing.T, itemID domain.ItemID) {
	buyItemStatus(t, itemID, domain.Day, http.StatusOK)
}

func buyItemMonth(t *testing.T, itemID domain.ItemID) {
	buyItemStatus(t, itemID, domain.Month, http.StatusOK)
}

func buyItemStatus(t *testing.T, itemID domain.ItemID, tf domain.TimeFrame, status int) {
	postBody := fmt.Sprintf(`{"itemId":"%s","units":1,"timeFrame":"%s"}`, itemID, tf)
	rr := testRouter.Request(t, http.MethodPost, "/items", &postBody, nil)
	assert.EqualValues(t, status, rr.Code, rr.Body.String())
}
