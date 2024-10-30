package e2e

import (
	"context"
	"github.com/glossd/pokergloss/achievement/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertGameEnd(t *testing.T) {
	t.Cleanup(cleanUp)
	gameEnd := oneWinnerTwoPlayers()
	assert.Nil(t, db.InsertGameEnd(context.Background(), gameEnd))
}
