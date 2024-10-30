package e2e

import (
	"encoding/json"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/multi"
	"github.com/glossd/pokergloss/table/web/client/mq"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"log"
	"testing"
)

func TestPlayersUpdates_StackAfterEachGame(t *testing.T) {
	multiSetUp(t)
	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{numOfUsers: 4})

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))

	lobby, err := db.FindLobbyMultiNoCtx(lobby.ID)
	assert.Nil(t, err)
	assert.Len(t, lobby.TableIDs, 2)

	hex0 := lobby.TableIDs[0].Hex()
	restMakeBetAction(t, hex0, domain.Raise, 10, 0, getToken(0))
	restMakeAction(t, hex0, domain.Call, 1, getToken(1))

	mq.ResetTestMQ()

	restMakeAction(t, hex0, domain.Fold, 1, getToken(1))
	assertActionGameEndWithOneShowdown(t, 1, 0)
	readMessage() // holeCards...
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertType(events.MultiPlayersUpdateType)
		as[0].assertPayload("tableId", hex0)
		var players []*model.UserWithStack
		playersJSON := []byte(gjson.Get(as[0].Event.Payload, "players").String())
		assert.Nil(t, json.Unmarshal(playersJSON, &players))
		assert.EqualValues(t, 2, len(players))
		wholeBet := defaultSmallBlind + 10
		assert.EqualValues(t, defaultBuyIn+wholeBet, getUser(0, players).Stack)
		assert.EqualValues(t, defaultBuyIn-wholeBet, getUser(1, players).Stack)
	})
}

func TestPlayersUpdates_ShouldPlusPlayersOnRemove(t *testing.T) {
	multiSetUp(t)

	lobby := insertFullLobbyMulti(t, NewLobbyMultiParams{tableSize: 4, numOfUsers: 5})

	algo, err := domain.NewMockAlgoMultiGame(domain.CardsStr("2s", "7d", "Ad", "As"))
	assert.Nil(t, err)
	domain.Algo = algo

	assert.Nil(t, multi.LaunchMultiLobbies(lobby.StartAt))
	lobby, err = db.FindLobbyMultiNoCtx(lobby.ID)
	assert.Nil(t, err)
	assert.Len(t, lobby.TableIDs, 2)

	communityCardsMock := domain.NewMockCards("Ac", "Ah", "Kd", "Ks", "Qh")

	domain.Algo = communityCardsMock
	hex0 := lobby.TableIDs[0].Hex()
	hex1 := lobby.TableIDs[1].Hex()
	restMakeAction(t, hex0, domain.AllIn, 0, getToken(0))
	restMakeAction(t, hex0, domain.AllIn, 1, getToken(1))
	restMakeAction(t, hex0, domain.Fold, 2, fifthPlayerToken)
	multi.Rebalance(lobby.ID)

	restMakeAction(t, hex1, domain.Call, 0, thirdPlayerToken)
	mq.ResetTestMQ()
	restMakeAction(t, hex1, domain.Fold, 1, fourthPlayerToken)
	assertActionGameEndWithOneShowdown(t, 1, 0)
	assertMultiPlayerMove(t, []string{thirdIdentity.UserId, fourthIdentity.UserId}, hex0)
	assertMessage(t, 2, func(as []*Asserter) {
		as[0].assertBankroll(0) // previous broke and left
		as[1].assertBankroll(3)
	})
	assertMessage(t, 1, func(as []*Asserter) {
		as[0].assertReset(4)
	})
	assertMessage(t, 2, func(as []*Asserter) {
		as[0].assertMultiPlayersUpdate(hex1, 0)
		as[1].assertMultiPlusPlayersUpdate(hex0, 2)
	})
}

func getUser(pos int, users []*model.UserWithStack) *model.UserWithStack {
	userID := getIden(pos).UserId
	for _, user := range users {
		if user.UserId == userID {
			return user
		}
	}
	log.Panicf("no user found in the user list")
	return nil
}
