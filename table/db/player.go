package db

import (
	"github.com/glossd/pokergloss/table/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PlayerNullify(position int) primitive.E {
	return primitive.E{Key: PlayerDbPath(position), Value: nil}
}

func PlayerLeaving(position int) primitive.E {
	return primitive.E{Key: PlayerDbPath(position) + ".isleaving", Value: true}
}

func PlayerAddChips(p *domain.Player) []bson.E {
	var es []bson.E
	if p.ChipsToAddOnReset > 0 {
		es = append(es, bson.E{Key: PlayerDbPath(p.Position) + ".chipstoaddonreset", Value: p.ChipsToAddOnReset})
	} else {
		es = append(es, bson.E{Key: PlayerDbPath(p.Position) + ".status", Value: p.Status})
	}
	return es
}
