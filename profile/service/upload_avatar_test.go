package service

import (
	"fmt"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestGetAvatarURL(t *testing.T) {
	result := getAvatarURL(authid.Identity{UserId: "userId"})
	fmt.Println(result)
	matched, err := regexp.MatchString("userId-.{4}", result)
	assert.Nil(t, err)
	assert.True(t, matched)
}
