package e2e

import (
	"encoding/json"
	"fmt"
	"github.com/glossd/pokergloss/messenger/domain"
	"github.com/glossd/pokergloss/messenger/web/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFirstChatList_ShouldReturnPhantomSupport(t *testing.T) {
	t.Cleanup(cleanUp)
	chats := getChats(t)
	assert.EqualValues(t, 1, len(chats))
	assert.EqualValues(t, "support", chats[0].Name)
	assert.True(t, chats[0].IsPhantom)
}

func TestCreateChat_ShouldReturnPhantomSupport(t *testing.T) {
	t.Cleanup(cleanUp)

	rr := testRouter.POST(t, "/u2u/chats", fmt.Sprintf(`{"userId":"%s"}`, secondIdentity.UserId), nil)
	assert.EqualValues(t, http.StatusCreated, rr.Code)

	chats := getChats(t)
	assert.EqualValues(t, 2, len(chats))
	assert.EqualValues(t, "support", chats[0].Name)
	assert.True(t, chats[0].IsPhantom)
}

func TestCreateChatWithSupport_ShouldNotReturnPhantomSupport(t *testing.T) {
	t.Cleanup(cleanUp)

	rr := testRouter.POST(t, "/u2u/chats", fmt.Sprintf(`{"userId":"%s"}`, domain.SupportUserID), nil)
	assert.EqualValues(t, http.StatusCreated, rr.Code)

	chats := getChats(t)
	assert.EqualValues(t, 1, len(chats))
	assert.False(t, chats[0].IsPhantom)
}

func getChats(t *testing.T) []*model.Chat {
	rr := testRouter.GET(t, "/chats", nil)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	var chats []*model.Chat
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &chats))
	return chats
}
