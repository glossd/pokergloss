package model

import "github.com/glossd/pokergloss/table/domain"

type Intent struct {
	Type  domain.IntentType `json:"type" enums:"fold,all-in,call,call-fold,call-any,raise,check-fold,check,check-call-any,bet"`
	Chips int64             `json:"chips"`
}

func ToIntent(i *domain.Intent) *Intent {
	if i == nil {
		return nil
	}
	return &Intent{
		Type:  i.Type,
		Chips: i.Chips,
	}
}
