package mqpub

import (
	conf "github.com/glossd/pokergloss/goconf"
	mqtable2 "github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

func PublishGameEnd(t *domain.Table) {
	if conf.IsLocalOnly() {
		return
	}
	if t.IsSurvival {
		return
	}
	winners := make([]*mqtable2.Winner, 0, len(t.Winners))
	for _, winner := range t.Winners {
		p, err := t.GetPlayer(winner.Position)
		if err != nil {
			log.Errorf("PublishGameEnd, winner position no player: %s", err)
			continue
		}
		winnerHandRank := winner.HandRank
		if winnerHandRank == "" && len(t.CommunityCards.AvailableCards()) >= 3 {
			// case where everyone did fold on winner's bet
			pos := winner.Position
			p, err := t.GetPlayer(pos)
			if err == nil {
				bestHand := t.ComputeBestHand(p)
				winnerHandRank = bestHand.GetRankStr()
			}
		}
		winners = append(winners, &mqtable2.Winner{
			UserId: p.UserId,
			Chips:  winner.Chips,
			Hand:   winnerHandRank,
		})
	}

	tablePlayers := t.PlayingPlayersByGameType()
	players := make([]*mqtable2.Player, 0, len(tablePlayers))
	for _, p := range tablePlayers {
		stackNoWinnings := p.Stack - p.GetWonChips()
		players = append(players, &mqtable2.Player{
			UserId:       p.UserId,
			WageredChips: p.StartGameStack - stackNoWinnings,
			LastAction:   string(p.LastGameAction),
			Hand:         p.HandRankString,
			Cards:        p.Cards.AsStringArray(),
			IsWinner:     p.GetWonChips() > 0,
			WonChips:     p.GetWonChips(),
			Rake:         t.GetRake().Of(p.Position),
		})
	}

	gameEnd := &mqtable2.GameEnd{
		TableType:      toTableType(t.Type),
		TableRound:     toTableRound(t.RoundType()),
		GameStartAt:    t.DecisionTimeoutAt / 1000,
		Winners:        winners,
		Players:        players,
		CommunityCards: t.CommunityCards.AsStringArray(),
		GotToShowDown:  t.WasShowDown(),
	}

	if conf.IsE2E() {
		mq.TestGameEndMQ <- gameEnd
		return
	}

	if conf.IsProd() {
		err := mqtable2.PublishGameEnd(gameEnd)

		if err != nil {
			log.Errorf("SendGameEndTableEvent: publish error: %s", err)
		}
	}
}

func toTableType(t domain.TableType) mqtable2.GameEnd_TableType {
	switch t {
	case domain.CashType:
		return mqtable2.GameEnd_LIVE
	case domain.SitngoType:
		return mqtable2.GameEnd_SITNGO
	case domain.MultiType:
		return mqtable2.GameEnd_MULTI
	default:
		return mqtable2.GameEnd_UNKNOWN_TABLE_TYPE
	}
}

func toTableRound(t domain.RoundType) mqtable2.GameEnd_TableRound {
	switch t {
	case domain.PreFlopRound:
		return mqtable2.GameEnd_PRE_FLOP
	case domain.FlopRound:
		return mqtable2.GameEnd_FLOP
	case domain.TurnRound:
		return mqtable2.GameEnd_TURN
	case domain.RiverRound:
		return mqtable2.GameEnd_RIVER
	default:
		return mqtable2.GameEnd_UNKNOWN_TABLE_ROUND
	}
}
