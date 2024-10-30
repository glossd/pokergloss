package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNextSmallBlind(t *testing.T) {
	table := map[int64]int64{
		1:2, 2:3, 3:5, 5:10, 4:5, 6:10, 9:10, 10:20, 20:30, 25:30, 30:50, 50:100, 73:100, 666:1000, 1000:2000, 3000:5000, 5*1e5:1e6,
	}
	for current, result := range table {
		got := nextSmallBlind(current)
		assert.EqualValues(t, result, got, fmt.Sprintf("sb=%d, expected=%d, got=%d", current, result, got))
	}
}
