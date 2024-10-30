package domain

import (
	"github.com/glossd/pokergloss/auth/authid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var defaultIden = authid.Identity{UserId: "2fR5MYyjqMSMgzLtdkdrKZkc18t1", Username: "denis", Picture: "https://storage.googleapis.com/pokerblow-avatars/2fR5MYyjqMSMgzLtdkdrKZkc18t1"}

func TestSurvivalArchData(t *testing.T) {
	data := [][]int{
		{9, 1, 120},
		{10, 2, 100},
		{11, 3, 80},
		{12, 4, 60},
		{13, 1, 140},
		{14, 2, 120},
		{15, 3, 100},
		{16, 4, 80},
		{17, 5, 60},
		{18, 1, 160},
	}
	for _, datum := range data {
		assertArchData(t, datum[0], datum[1], datum[2])
	}
}

func assertArchData(t *testing.T, lvl, num, stack int) {
	resNum, resStack := survivalLevel(lvl).archLevelData()
	assert.EqualValues(t, num, resNum)
	assert.EqualValues(t, stack, resStack)
}

func survivalLevel(lvl int) *Survival {
	s := New(defaultIden, Params{})
	for i := 1; i < lvl; i++ {
		s.NewLevel()
	}
	if s.Level != lvl {
		log.Panicf("level not equal")
	}
	return s
}
