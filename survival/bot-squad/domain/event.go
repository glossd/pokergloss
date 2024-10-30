package domain

import (
	"encoding/json"
	"github.com/glossd/pokergloss/survival/bot-squad/conf"
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Event struct {
	Type    string
	Payload string
}

func NewEventBytes(data []byte) *Event {
	return &Event{
		Type:    gjson.GetBytes(data, "type").String(),
		Payload: gjson.GetBytes(data, "payload").String(),
	}
}

func NewEvent(t, payload string) *Event {
	return &Event{
		Type:    t,
		Payload: payload,
	}
}

func (e *Event) Merge(c conf.Config, t *Table) error {
	newTable, err := e.GetTable()
	if err != nil {
		return err
	}
	seats := t.Seats
	err = mergo.Merge(t, newTable, mergo.WithOverride)
	if err != nil {
		return err
	}
	t.Seats = seats
	mergeSeats(t, newTable)

	if e.Type == "holeCards" {
		t.RankHoleCards(c.Squad.GetPositions())
	}

	if e.Type == "timeToDecide" {
		t.DecidingPosition = newTable.DecidingPosition
	}
	if e.Type == "newBettingRound" {
		t.SortCommunityCards()
		t.MaxRoundBet = 0
		for _, seat := range t.Seats {
			if seat.Player != nil {
				seat.Player.TotalRoundBet = 0
			}
		}
	}
	if e.Type == "reset" {
		t.CommCards = nil
		t.ResetPlayers()
	}

	if e.Type == "playerLeft" {
		t.Seats[newTable.Seats[0].Position].Player = nil
		t.Seats[newTable.Seats[0].Position].Blind = ""
	}

	return nil
}

func (e *Event) GetTable() (*Table, error) {
	tj := gjson.Get(e.Payload, "table").String()
	var t Table
	err := json.Unmarshal([]byte(tj), &t)
	if err != nil {
		log.Errorf("Failed to unmarshal table: %s", err)
		return nil, err
	}
	return &t, nil
}

func mergeSeats(t *Table, update *Table) {
	for _, seat := range update.Seats {
		oldSeat := t.Seats[seat.Position]
		err := mergo.Merge(oldSeat, seat, mergo.WithOverride)
		if err != nil {
			log.Errorf("Merge seat failed: %s", err)
		}
	}
}

func MergeTableBytes(c conf.Config, t *Table, data []byte) error {
	return NewEventBytes(data).Merge(c, t)
}
