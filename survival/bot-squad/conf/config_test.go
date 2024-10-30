package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringsToInts(t *testing.T) {
	t.Run("split", func(t *testing.T) {
		res := stringsToInts([]string{"1 2 3"}, "blah")
		assert.EqualValues(t, 3, len(res))
		assert.EqualValues(t, 1, res[0])
		assert.EqualValues(t, 2, res[1])
		assert.EqualValues(t, 3, res[2])
	})
	t.Run("no split", func(t *testing.T) {
		res := stringsToInts([]string{"1", "2", "3"}, "blah")
		assert.EqualValues(t, 3, len(res))
		assert.EqualValues(t, 1, res[0])
		assert.EqualValues(t, 2, res[1])
		assert.EqualValues(t, 3, res[2])
	})
}

func TestStringsToFloats(t *testing.T) {
	t.Run("split", func(t *testing.T) {
		res := stringsToFloats([]string{"0.65 0.6 1"}, "blah")
		assert.EqualValues(t, 3, len(res))
		assert.EqualValues(t, 0.65, res[0])
		assert.EqualValues(t, 0.6, res[1])
		assert.EqualValues(t, 1, res[2])
	})
	t.Run("no split", func(t *testing.T) {
		res := stringsToFloats([]string{"0.65", "0.6", "1"}, "blah")
		assert.EqualValues(t, 3, len(res))
		assert.EqualValues(t, 0.65, res[0])
		assert.EqualValues(t, 0.6, res[1])
		assert.EqualValues(t, 1, res[2])
	})
}
