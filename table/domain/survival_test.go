package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStartSurvivalLater(t *testing.T) {
	table, err := NewTableSurvival(NewSurvivalTableParams{
		User: User{Iden: firstIdentity, Stack: 100},
		Name: "Blah",
		DecisionTimeSec: 10,
		BigBlind: 2,
		Bots: []Bot{{Name:"bot1", Stack: 60}},
		ThemeID: Hell,
		LevelIncreaseTime: time.Minute,
		SurvivalLevel: 1,
	})

	assert.Nil(t, err)
	assert.True(t, table.IsWaiting())
	assert.EqualValues(t, -1, table.DecidingPosition)

	Algo = Algo_2P_MockSecondPlayerLoses()

	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsPlaying())
	assert.NotEqualValues(t, -1, table.DecidingPosition)

	table.tallIn(t)
	table.tallIn(t)
	assert.True(t, table.IsGameEnd())

	assert.Nil(t, table.StartNextGame())
	assert.True(t, table.IsWaiting())
	assert.EqualValues(t, 1, len(table.AllPlayers()))

}
