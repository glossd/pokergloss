package timeout

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Type string
const (
	SeatReservation Type = "seatReservation"
	Decision        Type = "decisionTimeout"
	StartGame       Type = "startGameTimeout"
)

type Event struct {
	Type Type  `json:"type"`
	// millis
	At   int64 `json:"at"`
	Key  Key   `json:"key"`
}

type Key struct {
	TableID  primitive.ObjectID
	Position int
	Version int64
}

func (k Key) String() string {
	return fmt.Sprintf("{tableID: %s, position: %d", k.TableID.Hex(), k.Position)
}
