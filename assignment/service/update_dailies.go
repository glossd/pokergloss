package service

import (
	"context"
	"github.com/glossd/pokergloss/assignment/db"
	"github.com/glossd/pokergloss/assignment/domain"
	"github.com/glossd/pokergloss/assignment/web/mq/bank"
	"github.com/glossd/pokergloss/assignment/web/mq/ws"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/gomq/mqtable"
	"time"
)

func UpdateDailiesWithTime(ctx context.Context, ofTime time.Time, ge *mqtable.GameEnd) error {
	d, err := db.FindDaily(ctx, domain.NewDailyID(ofTime))
	if err != nil {
		return err
	}

	var ackableErr error
	update := func(p *mqtable.Player) {
		err := updateDaily(ctx, d, p, ge, nil)
		if err != nil {
			ackableErr = err
		}
	}

	if len(ge.Winners) > 0 {
		if d.ContainsLosing() {
			for _, player := range ge.Players {
				update(player)
			}
		} else {
			for _, player := range ge.WonPlayers() {
				update(player)
			}
		}
	}

	if ackableErr != nil {
		return gomq.WrapInAckableError(ackableErr)
	}

	return nil
}

func UpdateDailiesTournament(ctx context.Context, te *mqtable.TournamentEnd) error {
	d, err := db.FindDaily(ctx, domain.NewDailyID(time.Now()))
	if err != nil {
		return err
	}

	var ackableErr error
	update := func(p *mqtable.Player) {
		err := updateDaily(ctx, d, p, nil, te)
		if err != nil {
			ackableErr = err
		}
	}

	if len(te.TournamentWinners) > 0 {
		for _, winner := range te.TournamentWinners {
			update(&mqtable.Player{UserId: winner.UserId})
		}
	}

	if ackableErr != nil {
		return gomq.WrapInAckableError(ackableErr)
	}

	return nil
}

func UpdateDailies(ctx context.Context, ge *mqtable.GameEnd) error {
	return UpdateDailiesWithTime(ctx, time.Now(), ge)
}

func updateDaily(ctx context.Context, daily *domain.Daily, p *mqtable.Player, ge *mqtable.GameEnd, te *mqtable.TournamentEnd) error {
	userDaily, err := FindOrBuildUserDaily(ctx, p.UserId, daily)
	if err != nil {
		return err
	}
	var changed bool
	var timeToWait time.Duration
	if ge != nil {
		changed = userDaily.Update(p, ge)
		timeToWait = time.Duration(ge.GameStartAt-time.Now().Unix()-1) * time.Second
	} else if te != nil {
		changed = userDaily.UpdateForTournament(p, te)
		timeToWait = 0
	}

	if changed {
		err := db.UpsertUserDaily(ctx, userDaily)
		if err == nil {
			for _, a := range userDaily.GetChanged() {
				if a.IsDone() {
					bank.DepositPrize(userDaily.UserID, &a.Assignment)
				}

				if timeToWait > 0 {
					time.AfterFunc(timeToWait, func() {
						ws.PublishAssignmentEvent(userDaily.UserID, a)
					})
				} else {
					ws.PublishAssignmentEvent(userDaily.UserID, a)
				}

			}
		}
		return err
	}
	return nil
}
