package service

import (
	"github.com/glossd/pokergloss/survival/bot-squad/conf"
	"github.com/glossd/pokergloss/survival/bot-squad/domain"
	"github.com/glossd/pokergloss/survival/bot-squad/web/client/rest"
	"github.com/glossd/pokergloss/survival/bot-squad/web/client/rpc"
	"time"
)

func Init(c conf.Config) error {
	t, err := initTable(c)
	if err != nil {
		return err
	}
	bots := initBots(c, t)
	StoreUserData(c.UserID, &UserData{Table: t, Bots: bots})
	return nil
}

func initTable(c conf.Config) (*domain.Table, error) {
	var t *domain.Table
	var err error
	if c.Protocol == "ws" {
		t, err = rest.GetTable(c)
	} else {
		t, err = rpc.GetTable(c.TableID)
	}
	if err != nil {
		return nil, err
	}
	t.UserPosition = c.UserPosition
	return t, nil
}

func initBots(c conf.Config, t *domain.Table) []*domain.Bot {
	weak := c.GetWeakPositionsSet()
	looseness := c.GetLooseness()
	aggression := c.GetAggression()
	var bots []*domain.Bot
	for i, position := range c.Squad.GetPositions() {
		_, isWeak := weak[position]
		bots = append(bots, &domain.Bot{
			Position:   position,
			Looseness:  looseness[i],
			Aggression: aggression[i],
			IsWeak:     isWeak,
		})
	}

	if t.Status != "waiting" {
		t.RankHoleCards(c.Squad.GetPositions())
		if t.DecidingPosition != c.UserPosition {
			bot := bots[t.DecidingPosition-1]
			time.AfterFunc(time.Second, func() { // give some time for subscribe for table events
				if c.Protocol == "ws" {
					rest.MakeAction(c, bot.Position, bot.GetAction(t), t)
				} else {
					rpc.MakeAction(c, bot.Position, bot.GetAction(t), t)
				}
			})
		}
	}
	return bots
}
