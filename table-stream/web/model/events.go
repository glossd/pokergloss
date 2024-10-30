package model

import (
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/gomq/mqws"
)

const (
	NewTableIdentity   = "newTableIdentity"
	NewTableAnonymous  = "newTableAnonymous"
	LeftTableIdentity  = "leftTableIdentity"
	LeftTableAnonymous = "leftTableAnonymous"
)

type WsEvent struct {
	Type    string    `json:"type" enums:"newTableIdentity,newTableAnonymous,leftTableIdentity,leftTableAnonymous"`
	Payload WsPayload `json:"payload"`
}
type WsPayload struct {
	User User `json:"user"`
}

// @ID events from websocket
// @Success 200 {array} WsEvent
// @Router /ws [get]
func BuildNewConnection(iden *authid.Identity) *mqws.Event {
	if iden == nil {
		return &mqws.Event{
			Type:    NewTableAnonymous,
			Payload: "{}",
		}
	} else {
		return &mqws.Event{
			Type:    NewTableIdentity,
			Payload: gomq.M{"user": iden}.JSON(),
		}
	}
}

func BuildLeftConnection(iden *authid.Identity) *mqws.Event {
	if iden == nil {
		return &mqws.Event{
			Type:    LeftTableAnonymous,
			Payload: "{}",
		}
	} else {
		return &mqws.Event{
			Type:    LeftTableIdentity,
			Payload: gomq.M{"user": iden}.JSON(),
		}
	}
}
