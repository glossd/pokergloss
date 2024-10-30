package e2e

import (
	"fmt"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services/multi"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMultiLobby_GameStart(t *testing.T) {
	t.Cleanup(cleanUp)

	lobby := defaultLobbyMulti()
	assert.Nil(t, db.InsertOneLobbyMulti(lobby))

	restMultiRegister(t, lobby.ID.Hex(), defaultToken)
	assertMultiRegistered(t)

	restMultiRegister(t, lobby.ID.Hex(), secondPlayerToken)
	assertMultiRegistered(t)

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))

	lobby, err := db.FindLobbyMultiNoCtx(lobby.ID)
	assert.Nil(t, err)
	assert.Len(t, lobby.TableIDs, 1)

	assertMultiGameStart(t, lobby.TableIDs[0].Hex())
}

func assertMultiRegistered(t *testing.T) {
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertMultiRegister()
	})
}

func restMultiRegister(t *testing.T, id string, token ...string) {
	rr := testRouter.Request(t, "PUT", fmt.Sprintf("/multi/lobbies/%s/register", id), nil, extractAuthHeaders(token))
	assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
}
