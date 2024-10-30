package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogarithm(t *testing.T) {
	assert.EqualValues(t, 2, Logarithm(25, 5))
	assert.EqualValues(t, 2.4306765580733933, Logarithm(50, 5))

	assert.EqualValues(t, 3, Logarithm(1000, 10))
}
