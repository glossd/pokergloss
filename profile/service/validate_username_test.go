package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateUsername(t *testing.T) {
	assert.Nil(t, validateUsername("denis_glotov"))
	assert.Nil(t, validateUsername("Denis"))

	assert.EqualValues(t, ErrSpaceUsername, validateUsername("as sa"))
	assert.EqualValues(t, ErrMinLengthUsername, validateUsername("as"))
	assert.EqualValues(t, ErrMaxLengthUsername, validateUsername("asdhfjhsdjkjfhadsjlfdsahjadfshlkjsdfasds"))
	assert.EqualValues(t, ErrUsernameWrong, validateUsername("denis@"))
	assert.NotNil(t, validateUsername("_denis"))
	assert.NotNil(t, validateUsername("denis_"))
	assert.NotNil(t, validateUsername("denis__glotov"))
}