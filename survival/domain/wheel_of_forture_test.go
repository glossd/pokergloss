package domain

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestCreateWoF(t *testing.T) {
	rand.Seed(time.Now().Unix())
	t.Run("onlyChips", func(t *testing.T) {
		wof := createWoF(1)
		assert.EqualValues(t, 12, len(wof.Slots))
	})

	t.Run("withItem", func(t *testing.T) {
		wof := createWoF(12)
		assert.EqualValues(t, 12, len(wof.Slots))
	})
}

func TestWoFChances(t *testing.T) {
	var sum float64
	for _, chance := range slotChances {
		sum += chance
	}
	assert.EqualValues(t, 1000, math.Round(sum*1000))
}