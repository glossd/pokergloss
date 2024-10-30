package ws

import "github.com/glossd/pokergloss/gomq/mqws"

func Broadcast(tableID string, payload []byte) {
	hub := getOrCreateTableHub(tableID)
	hub.broadcast <- payload
}

func Direct(tableID string, events *mqws.TableUserEvents) {
	hub := getOrCreateTableHub(tableID)
	hub.direct <- events
}
