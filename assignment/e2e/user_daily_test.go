package e2e

import (
	"context"
	"github.com/glossd/pokergloss/assignment/db"
	"github.com/glossd/pokergloss/assignment/domain"
	"github.com/glossd/pokergloss/assignment/service"
	"github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
	"time"
)

func TestFirstUserDaily(t *testing.T) {
	t.Cleanup(cleanUp)
	_, err := service.CreateDaily(time.Now())
	assert.Nil(t, err)
	rr := testRouter.GET(t, "/my/daily/assignments", nil)
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
	assert.EqualValues(t, 3, gjson.Get(rr.Body.String(), "#").Int())
}

func TestUserDaily(t *testing.T) {
	t.Cleanup(cleanUp)
	now := time.Now()
	insertDailyWithStraight(t, now)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err := service.UpdateDailies(ctx, gameEndWithStraight())
	assert.Nil(t, err)
	ud := findUserDaily(t, "1")
	assert.EqualValues(t, 3, len(ud.Assignments))
	assert.True(t, ud.Assignments[0].IsDone())
	assert.False(t, ud.Assignments[1].IsDone())
	assert.False(t, ud.Assignments[2].IsDone())
}

func TestUpdateUserDailyOfYesterday(t *testing.T) {
	t.Cleanup(cleanUp)
	now := time.Now()
	yesterday := now.Add(24 * time.Hour)
	insertDailyWithStraight(t, yesterday)

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	err := service.UpdateDailiesWithTime(ctx, yesterday, gameEndWithStraight())
	assert.Nil(t, err)
	ud := findUserDaily(t, "1")
	assert.True(t, ud.Assignments[0].IsDone())
	assert.EqualValues(t, domain.Straight, ud.Assignments[0].Hand)

	insertDailyWithFlush(t, now)
	ctx, _ = context.WithTimeout(context.Background(), 3*time.Second)
	err = service.UpdateDailiesWithTime(ctx, now, gameEndWithFlush())
	assert.Nil(t, err)
	ud = findUserDaily(t, "1")
	assert.True(t, ud.Assignments[0].IsDone())
	assert.EqualValues(t, domain.Flush, ud.Assignments[0].Hand)
}

func gameEndWithStraight() *mqtable.GameEnd {
	return &mqtable.GameEnd{
		Type:    mqtable.GameEnd_GAME_END,
		Winners: []*mqtable.Winner{{UserId: "1", Chips: 10, Hand: string(domain.Straight)}},
		Players: []*mqtable.Player{
			{UserId: "1", WageredChips: 5, LastAction: "check", Hand: string(domain.Straight), Cards: []string{"Ks", "Qd"}, IsWinner: true},
			{UserId: "2", WageredChips: 5, LastAction: "check", Hand: string(domain.HighCard), Cards: []string{"5s", "10d"}, IsWinner: false},
		},
	}
}

func gameEndWithFlush() *mqtable.GameEnd {
	return &mqtable.GameEnd{
		Type:    mqtable.GameEnd_GAME_END,
		Winners: []*mqtable.Winner{{UserId: "1", Chips: 10, Hand: string(domain.Flush)}},
		Players: []*mqtable.Player{
			{UserId: "1", WageredChips: 5, LastAction: "check", Hand: string(domain.Flush), Cards: []string{"Ks", "Qs"}, IsWinner: true},
			{UserId: "2", WageredChips: 5, LastAction: "check", Hand: string(domain.HighCard), Cards: []string{"5s", "10d"}, IsWinner: false},
		},
	}
}

func insertDailyWithStraight(t *testing.T, now time.Time) {
	d := domain.NewDaily(now)
	assignment := domain.NewAssignment(*domain.NewWithHand(domain.WinWithHand, domain.Straight))
	assignment.Count = 1
	d.Assignments[0] = assignment
	assert.Nil(t, db.InsertDaily(d))
}

func insertDailyWithFlush(t *testing.T, now time.Time) {
	d := domain.NewDaily(now)
	assignment := domain.NewAssignment(*domain.NewWithHand(domain.WinWithHand, domain.Flush))
	assignment.Count = 1
	d.Assignments[0] = assignment
	assert.Nil(t, db.InsertDaily(d))
}

func findUserDaily(t *testing.T, userID string) *domain.UserDaily {
	ud, err := db.FindUserDaily(context.Background(), userID)
	assert.Nil(t, err)
	return ud
}
