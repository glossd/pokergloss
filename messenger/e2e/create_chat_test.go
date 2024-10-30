package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/glossd/pokergloss/messenger/db"
	"github.com/glossd/pokergloss/messenger/web/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSendFirstMessage_ShouldCreateTwoUserChatLi(t *testing.T) {
	t.Cleanup(cleanUp)
	ctx := context.Background()
	chat := createChat(t)
	fmt.Println(*chat)

	_, err := db.FindUserChatList(ctx, firstIdentity.UserId)
	assert.Nil(t, err)
	_, err = db.FindUserChatList(ctx, secondIdentity.UserId)
	assert.Nil(t, err)
}

func createChat(t *testing.T) *model.Chat {
	rr := testRouter.POST(t, "/u2u/chats", fmt.Sprintf(`{"userId":"%s"}`, secondIdentity.UserId), authHeaders(firstToken))
	assert.EqualValues(t, http.StatusCreated, rr.Code, rr.Body.String())
	var chat model.Chat
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &chat))
	return &chat
}
