package domain

import "github.com/glossd/pokergloss/gomq/mqtable"

type GameEnd struct {
	*mqtable.GameEnd
	WinnersMap     map[string]struct{}
	PlayersMap     map[string]*mqtable.Player
	WageredPlayers []*mqtable.Player
}

func NewGameEnd(end *mqtable.GameEnd) *GameEnd {
	winnersMap := make(map[string]struct{})
	for _, winner := range end.Winners {
		winnersMap[winner.UserId] = struct{}{}
	}
	playersMap := make(map[string]*mqtable.Player)
	for _, p := range end.Players {
		playersMap[p.UserId] = p
	}
	var wageredPlayers []*mqtable.Player
	for _, p := range end.Players {
		if p.WageredChips > 0 {
			wageredPlayers = append(wageredPlayers, p)
		}
	}

	return &GameEnd{
		GameEnd:        end,
		WinnersMap:     winnersMap,
		PlayersMap:     playersMap,
		WageredPlayers: wageredPlayers,
	}
}

func (ge *GameEnd) AllInWageredPlayersExceptWinners() []*mqtable.Player {
	var players []*mqtable.Player
	for _, player := range ge.WageredPlayers {
		if player.LastAction == "allIn" {
			if _, ok := ge.WinnersMap[player.UserId]; !ok {
				players = append(players, player)
			}
		}
	}
	return players
}

func (ge *GameEnd) TillTheEndPlayersExceptWinners() []*mqtable.Player {
	var players []*mqtable.Player
	for _, player := range ge.WageredPlayers {
		if player.LastAction != "fold" && player.LastAction != "" {
			if _, ok := ge.WinnersMap[player.UserId]; !ok {
				players = append(players, player)
			}
		}
	}
	return players
}
