package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDailyBonus_CalculateBonus(t *testing.T) {
	b := NewDailyBonus("uid")
	b.Visit()
	assert.EqualValues(t, 500, b.CalculateBonus())

	b.Reset()
	b.Visit()
	assert.EqualValues(t, 1207, b.CalculateBonus())

	b.Reset()
	b.Visit()
	assert.EqualValues(t, 1500, b.CalculateBonus())

	b.Reset()
	b.Visit()
	assert.EqualValues(t, 1724, b.CalculateBonus())

	b.Reset()
	b.Visit()
	assert.EqualValues(t, 1914, b.CalculateBonus())
}
