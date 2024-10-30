package model

import (
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LobbyMulti struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	StartAt int64 `json:"startAt"`
	BuyIn   int64 `json:"buyIn"`

	TableSize            int                 `json:"tableSize"`
	BettingLimit         domain.BettingLimit `json:"bettingLimit"`
	StartingChips        int64               `json:"startingChips"`
	StartingBigBlind     int64               `json:"startingBigBlind"`
	StartingSmallBlind   int64               `json:"startingSmallBlind"`
	LevelIncreaseTimeMin int64               `json:"levelIncreaseTimeMin"`

	Status    domain.LobbyStatus `json:"status"`
	PrizePool int64              `json:"prizePool"`

	Players []Identity                `json:"players"`
	Prizes  *[]domain.TournamentPrize `json:"prizes"`

	RealPrizeDescription *string      `json:"realPrizeDescription,omitempty"`
	LastVideoID          *string      `json:"lastVideoID,omitempty"`
	MarketPrize          *MarketPrize `json:"marketPrize,omitempty"`

	TableIDs *[]primitive.ObjectID `json:"tableIds,omitempty"`
	Tables   *[]*Table             `json:"tables,omitempty"`
}

type Identity struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Picture  string `json:"picture"`
}

type MarketPrize struct {
	ItemID       string `json:"itemId"`
	NumberOfDays int    `json:"numberOfDays"`
}

func ToMultiLobbies(list []*domain.LobbyMulti) []*LobbyMulti {
	results := make([]*LobbyMulti, 0, len(list))
	for _, l := range list {
		results = append(results, ToLobbyMulti(l))
	}
	return results
}

func ToLobbyMulti(l *domain.LobbyMulti) *LobbyMulti {
	var tableIds *[]primitive.ObjectID
	if len(l.TableIDs) > 0 {
		tableIds = &l.TableIDs
	}
	var tables *[]*Table
	if len(l.GetTables()) > 0 {
		mtables := ToModelTables(l.GetTables())
		tables = &mtables
	}
	return &LobbyMulti{
		ID:                   l.ID.Hex(),
		Name:                 l.Name,
		StartAt:              l.StartAt,
		BuyIn:                l.BuyIn,
		PrizePool:            l.PrizePool(),
		RealPrizeDescription: l.RealPrizeDescription,
		LastVideoID:          l.LastVideoID,
		TableSize:            l.TableSize,
		BettingLimit:         l.BettingLimit,
		Status:               l.Status,
		StartingChips:        l.BuyIn,
		StartingBigBlind:     l.BigBlind,
		StartingSmallBlind:   l.BigBlind / 2,
		LevelIncreaseTimeMin: int64(l.LevelIncreaseTime.Minutes()),
		Players:              ToIdentities(l.Users),
		Prizes:               &l.Prizes,
		MarketPrize:          toMarketPrize(l.MarketPrize),
		TableIDs:             tableIds,
		Tables:               tables,
	}
}

func toMarketPrize(mp *domain.MarketPrize) *MarketPrize {
	var marketPrize *MarketPrize
	if mp != nil {
		marketPrize = &MarketPrize{
			ItemID:       mp.ItemID,
			NumberOfDays: mp.NumberOfDays,
		}
	}
	return marketPrize
}

func ToIdentities(idens []authid.Identity) []Identity {
	results := make([]Identity, 0, len(idens))
	for _, l := range idens {
		results = append(results, ToIdentity(l))
	}
	return results
}

func ToIdentity(iden authid.Identity) Identity {
	return Identity{
		UserID:   iden.UserId,
		Username: iden.Username,
		Picture:  iden.Picture,
	}
}

func (l *LobbyMulti) FillTables() {
	oid, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Errorf("LobbyMulti.FillTables failed to oid: %v", err)
		return
	}
	tables, err := db.FindTablesByLobbyID(oid)
	if err != nil {
		log.Errorf("LobbyMulti.FillTables failed to find tables: %v", err)
		return
	}

	result := ToModelTables(tables)
	l.Tables = &result
}
