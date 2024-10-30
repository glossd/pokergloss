package domain

import "github.com/glossd/pokergloss/gomq/mqtable"

// todo for future use polymorphism and unmarshalBSON
// https://levelup.gitconnected.com/golang-mongodb-with-polymorphism-and-bson-unmarshall-84779eb364c3

type Unit interface {
	getMaxCount() int64
	getPrize() int64
	match(p *mqtable.Player, ge *mqtable.GameEnd) bool
}

type WinWithHandUnit struct {
	Hand
}
