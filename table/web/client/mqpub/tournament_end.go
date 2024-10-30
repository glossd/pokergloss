package mqpub

import (
	conf "github.com/glossd/pokergloss/goconf"
	mqtable2 "github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
)

func PublishTournamentEndEvent(t *domain.Table) {
	if !t.IsTournament() {
		return
	}
	if conf.IsLocalOnly() {
		return
	}
	if t.IsSurvival {
		return
	}
	var playersWithPrize []*domain.Player
	for _, p := range t.NullifiedLeavingPlayers() {
		if p.GetTournamentInfo().Prize > 0 {
			playersWithPrize = append(playersWithPrize, p)
		}
	}

	winners := make([]*mqtable2.TournamentWinner, 0, len(playersWithPrize))
	for _, p := range playersWithPrize {
		winners = append(winners, &mqtable2.TournamentWinner{
			UserId: p.UserId,
			Place:  int64(p.GetTournamentInfo().Place),
		})
	}

	var msgType mqtable2.TournamentEnd_Type
	switch t.Type {
	case domain.SitngoType:
		msgType = mqtable2.TournamentEnd_SITNGO
	case domain.MultiType:
		msgType = mqtable2.TournamentEnd_MULTI
	default:
		log.Errorf("Failed to publish end: no such tournament end type mapper for %s", t.Type)
		return
	}

	var leftPlayers []*mqtable2.TournamentPlayer
	for _, p := range t.NullifiedLeavingPlayers() {
		leftPlayers = append(leftPlayers, &mqtable2.TournamentPlayer{
			UserId:   p.UserId,
			Place:    int64(p.GetTournamentInfo().Place),
			WonChips: p.GetTournamentInfo().Prize,
		})
	}

	end := &mqtable2.TournamentEnd{
		Type:              msgType,
		BuyIn:             t.TournamentAttributes.BuyIn,
		Fee:               t.TournamentAttributes.Fee(),
		Players:           leftPlayers,
		TournamentWinners: winners,
	}

	if conf.IsProd() {
		err := mqtable2.PublishTournamentEnd(end)
		if err != nil {
			log.Errorf("Publish tournament end: publish error: %s", err)
		}
	}
}
