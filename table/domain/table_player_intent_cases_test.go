package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntent_CallFold(t *testing.T) {
	t.Run("no upgrade, call", func(t *testing.T) {
		table := table3Players_startedGame(t)
		table.tintent(t, table.SmallBlindPosition(), CallFoldIntent)
		table.tcall(t)
		assert.EqualValues(t, Call, table.SmallBlindPlayer().LastGameAction)
		assert.EqualValues(t, table.BigBlindPosition(), table.DecidingPosition)
	})
	t.Run("upgrade to fold", func(t *testing.T) {
		table := table3Players_startedGame(t)
		table.tintent(t, table.SmallBlindPosition(), CallFoldIntent)
		table.traise(t, 20)
		assert.EqualValues(t, Fold, table.SmallBlindPlayer().LastGameAction)
		assert.EqualValues(t, table.BigBlindPosition(), table.DecidingPosition)
	})

}

func (t *Table) tintent(test *testing.T, pos int, intent Intent) {
	assert.Nil(test, t.SetIntent(pos, getIden(pos), intent))
}
