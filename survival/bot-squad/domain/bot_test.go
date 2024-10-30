package domain

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBot_PreFlopConfidence(t *testing.T) {
	b := Bot{Position: 1, Aggression: 1, Looseness: 0.5}
	table := tableBotDealer(botData{cards: holeCards("8d", "8c"), stack: 65})
	// reRaise from user
	table.MaxRoundBet = 24
	table.DecidingPlayer().TotalRoundBet = 12
	passConfidence := b.preFlopAggressionConfidence(table)
	log.Println(passConfidence)
	assert.Greater(t, passConfidence, 0.4)
	assert.Less(t, passConfidence, 1.0)
}
