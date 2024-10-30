package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/glossd/pokergloss/messenger/web/model"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestSendSecondMessage(t *testing.T) {
	t.Cleanup(cleanUp)
	chat := createChat(t)
	sendChatMessage(t, chat.ID, "hi")
	sendChatMessage(t, chat.ID, "hi, again")
	msgs := getChatMessages(t, chat.ID)
	assert.EqualValues(t, 2, len(msgs))
	assert.EqualValues(t, "hi, again", *msgs[0].Text)
	assert.EqualValues(t, "hi", *msgs[1].Text)
}

func TestGetChatMessages(t *testing.T) {
	t.Cleanup(cleanUp)
	chat := createChat(t)
	sendChatMessage(t, chat.ID, "hi 1")
	sendChatMessage(t, chat.ID, "hi 2")
	msg := sendChatMessage(t, chat.ID, "hi 3")
	sendChatMessage(t, chat.ID, "hi 4")
	msgs := getChatMessagesLastId(t, chat.ID, msg.ID)
	assert.EqualValues(t, 2, len(msgs))
	assert.EqualValues(t, "hi 2", *msgs[0].Text)
	assert.EqualValues(t, "hi 1", *msgs[1].Text)
}

func getChatMessages(t *testing.T, chatID string) []model.Message {
	path := fmt.Sprintf("/chats/%s/messages", chatID)
	rr := testRouter.GET(t, path, authHeaders(firstToken))
	assert.EqualValues(t, http.StatusOK, rr.Code)
	var result []model.Message
	assert.Nil(t, json.Unmarshal([]byte(gjson.Get(rr.Body.String(), "messages").String()), &result))
	return result
}

func getChatMessagesLastId(t *testing.T, chatID, lastID string) []model.Message {
	path := fmt.Sprintf("/chats/%s/messages?lastId=%s", chatID, lastID)
	rr := testRouter.GET(t, path, authHeaders(firstToken))
	assert.EqualValues(t, http.StatusOK, rr.Code)
	var result []model.Message
	assert.Nil(t, json.Unmarshal([]byte(gjson.Get(rr.Body.String(), "messages").String()), &result))
	return result
}

func sendChatMessage(t *testing.T, chatID string, text string) model.Message {
	body := fmt.Sprintf(`{"toUserId": "%s", "text":"%s"}`, secondIdentity.UserId, text)
	path := fmt.Sprintf("/chats/%s/messages", chatID)
	rr := testRouter.POST(t, path, body, authHeaders(firstToken))
	assert.EqualValues(t, http.StatusCreated, rr.Code, rr.Body.String())
	var result model.Message
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &result))
	return result
}

func sendFirstMessage(t *testing.T, text string) model.Message {
	body := fmt.Sprintf(`{"toUserId": "%s", "text":"%s"}`, secondIdentity.UserId, text)
	rr := testRouter.POST(t, "/u2u/messages", body, authHeaders(firstToken))
	assert.EqualValues(t, http.StatusCreated, rr.Code, rr.Body.String())
	var result model.Message
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &result))
	return result
}
