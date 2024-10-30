package e2e

import (
	"fmt"
	"github.com/glossd/pokergloss/table-chat/web/clients/wsclient"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestPostMessageTest(t *testing.T) {
	text := "hi, guys!"
	body := fmt.Sprintf(`{"text":"%s"}`, text)

	tableId := "123456abcd"
	rr := testRouter.Request(t, http.MethodPost, fmt.Sprintf("/tables/%s/messages", tableId), &body, nil)
	assert.EqualValues(t, rr.Code, http.StatusCreated)

	msg := <-wsclient.TestMQ
	assert.EqualValues(t, tableId, msg.ToEntityIds[0])
	assert.Len(t, msg.Events, 1)
	assert.EqualValues(t, "chatMessage", msg.Events[0].Type)
	payload := msg.Events[0].Payload
	assert.EqualValues(t, text, gjson.Get(payload, "text").String())
	assert.EqualValues(t, defaultIdentity.UserId, gjson.Get(payload, "user.userId").String())
	assert.EqualValues(t, defaultIdentity.Username, gjson.Get(payload, "user.username").String())
}
