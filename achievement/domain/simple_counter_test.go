package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleCounter(t *testing.T) {
	t.Run("NextLevelCount", func(t *testing.T) {
		sc := NewSimpleCounter(100, 5, false)
		assert.EqualValues(t, 5, sc.NextLevelCount())
		assert.EqualValues(t, 500, sc.NextLevelPrize())
	})
	t.Run("Inc", func(t *testing.T) {
		t.Run("FirstIncludedFalse", func(t *testing.T) {
			sc := NewSimpleCounter(100, 5, false)
			assert.EqualValues(t, 5, sc.NextLevelCount())
			sc.Inc()
			assert.EqualValues(t, 0, sc.Level)
		})
		t.Run("FirstIncludedTrue", func(t *testing.T) {
			sc := NewSimpleCounter(100, 5, true)
			assert.EqualValues(t, 1, sc.NextLevelCount())
			sc.Inc()
			assert.EqualValues(t, 1, sc.Level)
			assert.EqualValues(t, 5, sc.NextLevelCount())
		})
	})
}
