package playerbank

import (
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq/mqmarket"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	log "github.com/sirupsen/logrus"
)

func HandleNullifiedPlayersLeft(table *domain.Table) []*events.TableEvent {
	var tableEvents []*events.TableEvent
	nullifiedPlayers := table.NullifiedLeavingPlayers()
	if len(nullifiedPlayers) > 0 {
		if table.IsMultiType() {
			handleMultiNullifiedPlayers(table)
		}
		for _, leftPlayer := range nullifiedPlayers {
			tableEvents = append(tableEvents, events.BuildPlayerLeft(leftPlayer))
		}

		for _, leftPlayer := range nullifiedPlayers {
			SendPlayerChipsToBank(leftPlayer, table)
		}

		mqpub.PublishTournamentEndEvent(table)
	}
	return tableEvents
}

func handleMultiNullifiedPlayers(table *domain.Table) {
	tables, err := db.FindTablesByLobbyID(table.LobbyID)
	if err != nil {
		log.Errorf("Failed to set multi tournament info for left player: %s", err)
		return
	}
	var allPlayersCount int
	for _, t := range tables {
		allPlayersCount += len(t.AllPlayers())
	}

	sortedNullifiedPlayers := table.MultiSortedNullifiedPlayers()
	for i := len(sortedNullifiedPlayers) - 1; i >= 0; i-- {
		allPlayersCount++
		table.MultiSetTournamentInfo(sortedNullifiedPlayers[i], allPlayersCount)
	}
	if table.TournamentAttributes.MarketPrize != nil {
		for _, p := range sortedNullifiedPlayers {
			if p.GetTournamentInfo().Place == 1 {
				if conf.IsProd() {
					err := mqmarket.PublishGift(&mqmarket.Gift{
						ToUserId:  p.UserId,
						ItemId:    table.MarketPrize.ItemID,
						Units:     int64(table.MarketPrize.NumberOfDays),
						TimeFrame: mqmarket.Gift_DAY,
					})
					if err != nil {
						log.Errorf("Failed to market prize: %s", err)
					}
					p.SetTournamentMarketPrize(table.TournamentAttributes.MarketPrize)
				}
			}
		}
	}
}
