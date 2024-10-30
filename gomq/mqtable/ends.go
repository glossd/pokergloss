package mqtable

import (
	"context"
	"github.com/glossd/memmq"
	"github.com/glossd/pokergloss/gomq"
	log "github.com/sirupsen/logrus"
)

const TopicID = "pg.table.game-end"
const SurvivalEndTopicID = "pg.table.survival-end"
const TournamentEndTopicID = "pg.table.tournament-end"

// Deprecated, use PublishGameEnd.
func Publish(msg *GameEnd) error {
	return PublishGameEnd(msg)
}

func PublishGameEnd(msg *GameEnd) error {
	return memmq.Publish(TopicID, msg)
}

func PublishSurvivalEnd(msg *SurvivalEnd) error {
	return memmq.Publish(SurvivalEndTopicID, msg)
}

func PublishTournamentEnd(msg *TournamentEnd) error {
	return memmq.Publish(TournamentEndTopicID, msg)
}

// Deprecated, use SubscribeGameEnd.
func Subscribe(subID string, receiver func(ctx context.Context, msg *GameEnd) error) error {
	return SubscribeGameEnd(subID, receiver)
}

func SubscribeGameEnd(subID string, receiver func(ctx context.Context, msg *GameEnd) error) error {
	return memmq.Subscribe(TopicID, func(msg interface{}) bool {
		v, ok := msg.(*GameEnd)
		if !ok {
			log.Errorf("memmq: expected *GameEnd, got: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return gomq.IsAckableError(err)
	})
}

func SubscribeTournamentEnd(subID string, receiver func(ctx context.Context, msg *TournamentEnd) error) error {
	return memmq.Subscribe(TournamentEndTopicID, func(msg interface{}) bool {
		v, ok := msg.(*TournamentEnd)
		if !ok {
			log.Errorf("memmq: expected *TournamentEnd, got: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return gomq.IsAckableError(err)
	})
}

func SubscribeSurvivalEnd(subID string, receiver func(ctx context.Context, msg *SurvivalEnd) error) error {
	return memmq.Subscribe(SurvivalEndTopicID, func(msg interface{}) bool {
		v, ok := msg.(*SurvivalEnd)
		if !ok {
			log.Errorf("memmq: expected *SurvivalEnd, got: %T", v)
			return true
		}
		err := receiver(context.Background(), v)
		return gomq.IsAckableError(err)
	})
}

func (x *GameEnd) LostEndPlayers() (players []*Player) {
	if x == nil {
		return nil
	}
	for _, p := range x.Players {
		if !p.IsWinner && (p.LastAction != "" && p.LastAction != "fold") {
			players = append(players, p)
		}
	}
	return players
}

func (x *GameEnd) WonPlayers() (players []*Player) {
	if x == nil {
		return nil
	}
	for _, p := range x.Players {
		if p.IsWinner {
			players = append(players, p)
		}
	}
	return players
}

func (x *GameEnd) LostPlayers() (players []*Player) {
	if x == nil {
		return nil
	}
	for _, p := range x.Players {
		if !p.IsWinner {
			players = append(players, p)
		}
	}
	return players
}

func (x *GameEnd) LostChipsPlayers() (players []*Player) {
	if x == nil {
		return nil
	}
	for _, p := range x.Players {
		if !p.IsWinner && p.WageredChips > 0 {
			players = append(players, p)
		}
	}
	return players
}
