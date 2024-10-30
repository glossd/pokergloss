package domain

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeRakeFrom(t *testing.T) {
	type data struct {
		chips int64
		rakePercent float64
		rake int64
	}

	onePercent := 0.01
	var cases = []data{
		{0, onePercent, 0},
		{1, onePercent, 0},
		{99, onePercent, 0},
		{100, onePercent, 1},
		{101, onePercent, 1},
		{999, onePercent, 9},
		{654, onePercent, 6},
		{1e8, onePercent, 1e6},
	}
	for _, d := range cases {
		got := computeRakeFrom(d.chips, d.rakePercent)
		assert.EqualValues(t, d.rake, got, fmt.Sprintf("chips %d, expected %d, got %d", d.chips, d.rake, got))
	}
}
