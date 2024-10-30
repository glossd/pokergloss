package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCards_IsStraightBackdoorDraw(t *testing.T) {
	trueSBD(t, 0, "Ks", "Qc", "Js", "Td", "9c")

	trueSBD(t, 0, "Ks", "Qc", "Js", "Td")

	trueSBD(t, 0, "As", "3d", "2c")
	trueSBD(t, 0, "As", "Kd", "Qc")
	trueSBD(t, 0, "Kd", "Qc", "Js")
	trueSBD(t, 1, "Kd", "Qc", "Ts")
	trueSBD(t, 2, "Kd", "Qc", "9s")
	trueSBD(t, 2, "Kd", "Qc", "9s")
	falseSBD(t, "Kd", "Qc", "8s")
}

func trueSBD(t *testing.T, gaps int, cards ...string) {
	is, resGaps := cardsStr(cards...).IsStraightBackdoorDraw()
	assert.True(t, is)
	assert.EqualValues(t, gaps, resGaps)
}

func falseSBD(t *testing.T, cards ...string) {
	is, resGaps := cardsStr(cards...).IsStraightBackdoorDraw()
	assert.False(t, is)
	assert.Zero(t, resGaps)
}
