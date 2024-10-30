package domain

import (
	"github.com/chehsunliu/poker"
	log "github.com/sirupsen/logrus"
	"sort"
)

type Winner struct {
	Position int
	Chips int64
	HandRank string
}

type Hand []Card

type BestHand struct {
	rank int32
	rankString string
}

func (bh *BestHand) GetRankStr() string {
	return bh.rankString
}

func (t *Table) ComputeWinners() {
	gamingPlayers := t.GamingPlayers()
	if len(gamingPlayers) == 0 {
		if t.IsMultiType() && t.IsPreFlop() {
			// case when all sitting out and on playing, but made fold
			bbPlayer := t.GetBigBlind()
			if bbPlayer != nil {
				t.Winners = []Winner{{Position: bbPlayer.Position, Chips: t.TotalPot}}
				t.Pots.setAllWinnerPos(bbPlayer.Position)
				return
			}
		}
	}
	if len(gamingPlayers) == 1 {
		player := t.lastGamingPlayer()
		t.Winners = []Winner{{Position: player.Position, Chips: t.TotalPot}}
		t.Pots.setAllWinnerPos(player.Position)
		return
	}

	if t.IsSitngoType() || t.IsMultiType() {
		if t.AreAllSittingOut() {
			winnerPos := t.DecidingPosition
			if winnerPos < 0 {
				p, err := t.GetRandomPlayer()
				if err != nil {
					log.Errorf("Table %s: %s", t.ID.Hex(), err)
				} else {
					winnerPos = p.Position
				}
			}
			t.Winners = []Winner{{Position: winnerPos, Chips: t.TotalPot}}
			t.Pots.setAllWinnerPos(winnerPos)
			return
		}
	}

	t.computeRanks()

	players := t.GamingPlayers()
	rankPlayers := make(map[int32][]*Player)
	for _, player := range players {
		rankPlayers[player.HandRank] = append(rankPlayers[player.HandRank], player)
	}
	// make([]T, 0, capacity) https://stackoverflow.com/a/16907455/10160865
	ranks := make([]int32, 0, len(rankPlayers))
	for rank := range rankPlayers {
		ranks = append(ranks, rank)
	}

	// https://stackoverflow.com/a/48568680/10160865
	sort.Slice(ranks, func(i, j int) bool {return ranks[i] < ranks[j]})

	t.lastTakenPotIdx = -1
	t.Pots.setIndexes()
	t.removeLastPotIfEmpty()
	for _, rank := range ranks {
		if t.lastTakenPotIdx == len(t.Pots)-1 {
			break
		}
		playersSameRank := rankPlayers[rank]
		for _, p := range playersSameRank {
			p.potIndex = t.Pots.getPlayerPotsIdx(p.UserId)
		}

		for _, pot := range t.Pots.slice(t.lastTakenPotIdx+1) {
			playersSharing := playersSharingPot(playersSameRank, t.Pots, pot.idx)
			if len(playersSharing) == 0 {
				break
			}
			var potWinPositions []int
			for _, player := range playersSharing {
				potWinPositions = append(potWinPositions, player.Position)
			}
			pot.WinnerPositions = potWinPositions
			t.lastTakenPotIdx = pot.idx
		}
	}
	t.Winners = t.buildWinners()
}

func playersSharingPot(sameRankPlayers []*Player, pots Pots, potIdx int) (players []*Player) {
	for _, p := range sameRankPlayers {
		playerPotIdx := pots.getPlayerPotsIdx(p.UserId)
		if playerPotIdx >= potIdx {
			players = append(players, p)
		}
	}
	return
}

func (t *Table) computeRanks() {
	players := t.GamingPlayers()
	for _, p := range players {
		bestHand := t.ComputeBestHand(p)
		p.HandRank = bestHand.rank
		p.HandRankString = bestHand.rankString
	}
}

func (t *Table) ComputeBestHand(p *Player) BestHand {
	communityCards := t.CommunityCards.AvailableCards()
	if len(communityCards) < 3 {
		log.Warnf("Table.ComputeBestHand invoked with zero community cards")
		return BestHand{}
	}
	return chooseBest5(append(communityCards, p.Cards.Get()...))
}

// Returns highest hand
func chooseBest5(cards []Card) BestHand {
	rank := poker.Evaluate(mapCardsTo(cards))
	return BestHand{rank: rank, rankString: poker.RankString(rank)}
}

func mapCardsTo(cards []Card) []poker.Card {
	var mapped []poker.Card
	for _, card := range cards {
		mapped = append(mapped, mapCardTo(card))
	}
	return mapped
}

func mapCardTo(card Card) poker.Card {
	return poker.NewCard(card.String())
}

