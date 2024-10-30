package e2e

import (
	"fmt"
	"github.com/glossd/pokergloss/table-chat/web/clients/wsclient"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestPostEmoji(t *testing.T) {
	emoji := "joy"
	tableId := "123456abcd"
	postEmoji(t, tableId, emoji, http.StatusCreated)

	msg := <-wsclient.TestMQ
	assert.EqualValues(t, tableId, msg.ToEntityIds[0])
	assert.Len(t, msg.Events, 1)
	assert.EqualValues(t, "emojiMessage", msg.Events[0].Type)
	payload := msg.Events[0].Payload
	assert.EqualValues(t, emoji, gjson.Get(payload, "emoji").String())
	assert.EqualValues(t, defaultIdentity.UserId, gjson.Get(payload, "user.userId").String())
	assert.EqualValues(t, defaultIdentity.Username, gjson.Get(payload, "user.username").String())
}

func TestFailToPostEmoji(t *testing.T) {
	emoji := "blah"
	tableId := "123456abcd"
	postEmoji(t, tableId, emoji, http.StatusBadRequest)
}

func postEmoji(t *testing.T, tableId string, emoji string, code int) {
	body := fmt.Sprintf(`{"emoji":"%s"}`, emoji)
	rr := testRouter.Request(t, http.MethodPost, fmt.Sprintf("/tables/%s/emojis", tableId), &body, nil)
	assert.EqualValues(t, rr.Code, code)
}
