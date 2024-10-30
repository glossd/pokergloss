package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHoleCards_GetPattern(t *testing.T) {
	assert.EqualValues(t, "ATo", HoleCards{"Ad", "Tc"}.GetPattern())
	assert.EqualValues(t, "ATs", HoleCards{"Ad", "Td"}.GetPattern())
}

func TestHoleCards_GetPatternRank(t *testing.T) {
	assert.EqualValues(t, 41, HoleCards{"Ad", "Tc"}.GetPatternRank())
}
