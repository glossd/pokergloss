package db

import (
	"github.com/glossd/pokergloss/table/domain"
	"go.mongodb.org/mongo-driver/bson"
)

func TableUpdate(t *domain.Table) []bson.E {
	updates := CommunityCards(t.CommunityCards)
	updates = append(updates, bson.E{Key: "status", Value: t.Status})
	updates = append(updates, bson.E{Key: "waitingat", Value: t.WaitingAt})
	updates = append(updates, bson.E{Key: "decidingposition", Value: t.DecidingPosition})
	updates = append(updates, bson.E{Key: "decisiontimeoutat", Value: t.DecisionTimeoutAt})
	updates = append(updates, bson.E{Key: "pots", Value: t.Pots})
	updates = append(updates, bson.E{Key: "totalpot", Value: t.TotalPot})
	updates = append(updates, bson.E{Key: "roundpot", Value: t.RoundPot})
	updates = append(updates, bson.E{Key: "winners", Value: t.Winners})
	updates = append(updates, bson.E{Key: "lastaggressorposition", Value: t.LastAggressorPosition})
	updates = append(updates, PlayersCount(t))

	if t.IsSitngoType() || t.IsMultiType() {
		updates = append(updates, TableUpdateForTournament(t)...)
	}
	return updates
}

func AllTableUpdatesGameFlow(t *domain.Table) []bson.E {
	updates := SeatUpdates(t)
	updates = append(updates, TableUpdatesGameFlow(t)...)
	return updates
}

func TableUpdatesGameFlow(t *domain.Table) []bson.E {
	return append(TableUpdate(t), bson.E{Key: "gameflowversion", Value: t.GameFlowVersion + 1})
}

func TableUpdateForTournament(t *domain.Table) []bson.E {
	var updates []bson.E
	updates = append(updates, bson.E{Key: "bigblind", Value: t.BigBlind})
	updates = append(updates, bson.E{Key: "smallblind", Value: t.SmallBlind})
	updates = append(updates, bson.E{Key: "tournamentattributes", Value: t.TournamentAttributes})
	return updates
}

func PlayersCount(t *domain.Table) bson.E {
	return bson.E{Key: "playerscount", Value: len(t.AllPlayers())}
}

func CommunityCards(cards *domain.CommunityCards) []bson.E {
	var updates []bson.E
	updates = append(updates, bson.E{Key: "communitycards.flop", Value: cards.Flop})
	updates = append(updates, bson.E{Key: "communitycards.turn", Value: cards.Turn})
	updates = append(updates, bson.E{Key: "communitycards.river", Value: cards.River})
	return updates
}

func PlayerAutoConfig(pos int, config *domain.PlayerAutoConfig) []bson.E {
	PlayerDbPath(pos)
	return []bson.E{{Key: PlayerDbPath(pos) + ".autoconfig", Value: config}}
}

func TableMultiIsLast() bson.E {
	return bson.E{Key: "multiattrs.islast", Value: true}
}

func PlayerEnrichment(p *domain.Player) []bson.E {
	var updates []bson.E
	updates = append(updates, bson.E{Key: PlayerDbPath(p.Position) + ".level", Value: p.Level})
	updates = append(updates, bson.E{Key: PlayerDbPath(p.Position) + ".bankbalance", Value: p.BankBalance})
	updates = append(updates, bson.E{Key: PlayerDbPath(p.Position) + ".bankrank", Value: p.BankRank})
	updates = append(updates, bson.E{Key: PlayerDbPath(p.Position) + ".marketitemid", Value: p.MarketItemID})
	updates = append(updates, bson.E{Key: PlayerDbPath(p.Position) + ".marketitemcoinsdayprice", Value: p.MarketItemCoinsDayPrice})
	updates = append(updates, bson.E{Key: PlayerDbPath(p.Position) + ".autoconfig", Value: p.AutoConfig})
	return updates
}
