package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWinCounterComputePrize(t *testing.T) {
	wc := NewWinCounter()
	wc.Inc()
	assert.Zero(t, wc.GetPrize().Chips)
}
