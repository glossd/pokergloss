package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/glossd/pokergloss/market/domain"
	"github.com/glossd/pokergloss/market/web/rest/model"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"testing"
)

func TestSelectItemOnBuy(t *testing.T) {
	t.Cleanup(cleanUp)
	buyGlassOfWine(t)
	body := listUserItems(t)
	for _, result := range gjson.Get(body, "@this").Array() {
		if gjson.Get(result.String(), "itemId").String() == string(domain.GlassOfWine.ID) {
			assert.True(t, gjson.Get(result.String(), "selected").Bool())
			continue
		}
		if gjson.Get(result.String(), "itemId").String() == string(domain.Invisible.ID) {
			assert.False(t, gjson.Get(result.String(), "selected").Bool())
			continue
		}
		log.Panicf("no such item")
	}
}

func TestSelectItem(t *testing.T) {
	t.Cleanup(cleanUp)
	buyGlassOfWine(t)
	buyItem(t, domain.CloverID)

	rr := testRouter.Request(t, http.MethodPut, fmt.Sprintf("/items/%s/select", domain.GlassOfWineID), nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())

	body := listUserItems(t)
	var items []*model.UserItem
	assert.Nil(t, json.Unmarshal([]byte(body), &items))
	for _, item := range items {
		switch item.ItemID {
		case domain.GlassOfWineID:
			assert.True(t, item.Selected)
		case domain.CloverID:
			assert.False(t, item.Selected)
		}
	}
}
