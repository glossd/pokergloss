package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildPriceList(t *testing.T) {
	pl := buildPriceList(5000)
	assert.EqualValues(t, 5000, pl.Day)
	assert.EqualValues(t, 29750, pl.Week)
	assert.EqualValues(t, 105000, pl.Month)

	pl = buildPriceList(3000)
	assert.EqualValues(t, 3000, pl.Day)
	assert.EqualValues(t, 17850, pl.Week)
	assert.EqualValues(t, 63000, pl.Month)

	pl = buildPriceList(1)
	assert.EqualValues(t, 1, pl.Day)
	assert.EqualValues(t, 6, pl.Week)
	assert.EqualValues(t, 21, pl.Month)
}
