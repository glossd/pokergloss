package e2e

import (
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/services/multi"
	"github.com/glossd/pokergloss/table/services/recovery"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultiRecovery_NoDocs(t *testing.T) {
	t.Cleanup(cleanUp)

	recovery.InitMultiSchedulerRecovery()

	lobbies, err := db.FindAllMultiLobbiesNoCtx()
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(lobbies), multi.NumberOfTournaments)
}
