package service

import (
	"github.com/glossd/pokergloss/survival/bot-squad/conf"
	"github.com/glossd/pokergloss/survival/bot-squad/domain"
	"github.com/glossd/pokergloss/survival/bot-squad/web/client/rest"
	"github.com/glossd/pokergloss/survival/bot-squad/web/client/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"math/rand"
	"time"
)

func HandleEvents(config conf.Config, events []*domain.Event) {
	var filteredEvents []*domain.Event
	for _, event := range events {
		if gjson.Get(event.Payload, "table").Exists() {
			filteredEvents = append(filteredEvents, event)
		}
	}
	userData := GetUserData(config.UserID)
	if userData == nil {
		return
	}
	theTable := userData.Table
	theBots := userData.Bots
	var prevEvent *domain.Event
	for _, event := range filteredEvents {
		log.Tracef("Event: %s", event.Type)
		if event.Type == "playerAction" {
			tableUpdate, err := domain.NewEvent(event.Type, event.Payload).GetTable()
			if err != nil {
				log.Errorf("playerAction handling failed: %s", err)
				continue
			}
			if tableUpdate.Seats[0].Position == config.UserPosition {
				domain.UpdateUserTracker(theTable, tableUpdate.Seats[0].Player)
			}
		}

		err := event.Merge(config, theTable)
		if err != nil {
			log.Errorf("Failed to merge event %s: %s", event.Type, err)
		}
		if event.Type == "holeCards" {
			for _, bot := range theBots {
				bot.GameReset()
			}
		}
		if event.Type == "newBettingRound" {
			for _, bot := range theBots {
				bot.RoundMadeActionCount = 0
			}
		}
		if event.Type == "timeToDecide" && theTable.DecidingPosition != config.UserPosition {
			log.Tracef("Bot deciding at %d", theTable.DecidingPosition)
			bot := theBots[theTable.DecidingPosition-1]
			action := bot.GetAction(theTable)
			timeAcc := simulateThinking(action, bot, theTable)
			if prevEvent != nil && prevEvent.Type == "newBettingRound" {
				timeAcc += 1000 * time.Millisecond
			}
			if prevEvent != nil && prevEvent.Type == "holeCards" {
				timeAcc += 1000 * time.Millisecond
			}
			time.AfterFunc(timeAcc, func() {
				if config.Protocol == "ws" {
					rest.MakeAction(config, bot.Position, action, theTable)
				} else {
					rpc.MakeAction(config, bot.Position, action, theTable)
				}
			})
		}

		if event.Type == "timeToDecideTimeout" {
			outPosition := int(gjson.Get(event.Payload, "table.seats.0.position").Int())
			log.Tracef("Timeout at position %d", outPosition)
			if config.UserPosition != outPosition {
				if config.Protocol == "ws" {
					rest.SitBack(config, outPosition)
				} else {
					rpc.SitBack(config, outPosition)
				}
			}
		}

		prevEvent = event
	}
}

func simulateThinking(a domain.Action, bot *domain.Bot, t *domain.Table) time.Duration {
	min := 500 * time.Millisecond

	if bot.Aggression > 0.8 {
		return min
	}

	switch a.Type {
	case domain.CheckType:
		return min
	case domain.FoldType:
		if t.MaxRoundBet > t.BigBlind {
			return min
		}
		if bot.GetConfidence() > 0.3 {
			if rand.Float64() < 0.2 {
				return time.Second
			}
			return min + confidenceTime(bot, time.Second)
		}
		return min
	case domain.CallType:
		if t.MaxRoundBet > t.BigBlind {
			if rand.Float64() < 0.2 {
				return time.Second
			}
			return min + confidenceTime(bot, time.Second)
		}
		return min
	case domain.BetType, domain.RaiseType:
		if rand.Float64() < 0.2 {
			return 3 * time.Second
		}
		return min + confidenceTime(bot, 3*time.Second)
	}
	return min
}

func confidenceTime(bot *domain.Bot, t time.Duration) time.Duration {
	return time.Duration((1 - bot.GetConfidence()) * float64(t))
}
