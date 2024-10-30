package events

import (
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
)

const Enrichment TET = "enrichPlayers"

func BuildEnrichment(players []*domain.Player) *TableEvent {
	return &TableEvent{Type: Enrichment, Payload: TEP{Table: model.TablePlayers(players, model.ToPlayerEnrichment)}}
}
