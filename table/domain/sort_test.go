package domain

import (
	"github.com/stretchr/testify/assert"
"github.com/glossd/pokergloss/auth/authid"
"testing"
)

func TestSortPlayersByTotalRoundBet(t *testing.T) {
	p1 := &Player{Identity: authid.Identity{UserId: "rich"}, TotalRoundBet: 50}
	p2 := &Player{Identity: authid.Identity{UserId: "poor"}, TotalRoundBet: 10}
	p3 := &Player{Identity: authid.Identity{UserId: "avg"}, TotalRoundBet: 25}
	slice := []*Player{p1, p2, p3}
	SortPlayersByTotalRoundBet(slice)
	assert.EqualValues(t, "poor", slice[0].UserId)
	assert.EqualValues(t, "avg", slice[1].UserId)
	assert.EqualValues(t, "rich", slice[2].UserId)
}

func TestSortPlayersByHandRank(t *testing.T) {
	p1 := &Player{Identity: authid.Identity{UserId: "middle"}, HandRank: 1742}
	p2 := &Player{Identity: authid.Identity{UserId: "winner"}, HandRank: 1}
	p3 := &Player{Identity: authid.Identity{UserId: "loser"}, HandRank: 6695}
	slice := []*Player{p1, p2, p3}
	SortPlayersByHandRank(slice)
	assert.EqualValues(t, "winner", slice[0].UserId)
	assert.EqualValues(t, "middle", slice[1].UserId)
	assert.EqualValues(t, "loser", slice[2].UserId)
}

func TestSortPlayersByStartGameStack(t *testing.T) {
	p1 := &Player{Identity: authid.Identity{UserId: "middle"}, StartGameStack: 100}
	p2 := &Player{Identity: authid.Identity{UserId: "poor"}, StartGameStack: 5}
	p3 := &Player{Identity: authid.Identity{UserId: "rich"}, StartGameStack: 250}
	slice := []*Player{p1, p2, p3}
	SortPlayersByStartGameStack(slice)
	assert.EqualValues(t, "poor", slice[0].UserId)
	assert.EqualValues(t, "middle", slice[1].UserId)
	assert.EqualValues(t, "rich", slice[2].UserId)
}

func TestSortPlayersByStack(t *testing.T) {
	p1 := &Player{Identity: authid.Identity{UserId: "middle"}, Stack: 100}
	p2 := &Player{Identity: authid.Identity{UserId: "poor"}, Stack: 5}
	p3 := &Player{Identity: authid.Identity{UserId: "rich"}, Stack: 250}
	slice := []*Player{p1, p2, p3}
	SortPlayersByStack(slice)
	assert.EqualValues(t, "poor", slice[0].UserId)
	assert.EqualValues(t, "middle", slice[1].UserId)
	assert.EqualValues(t, "rich", slice[2].UserId)
}
