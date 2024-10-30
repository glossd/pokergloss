package service

import (
	"github.com/glossd/pokergloss/survival/bot-squad/domain"
	"sync"
)

var botMap sync.Map

type UserData struct {
	Bots  []*domain.Bot
	Table *domain.Table
}

func StoreUserData(userID string, bots *UserData) {
	botMap.Store(userID, bots)
}

func GetUserData(userID string) *UserData {
	v, ok := botMap.Load(userID)
	if !ok {
		return nil
	}
	ud, ok := v.(*UserData)
	if !ok {
		return nil
	}
	return ud
}

func DeleteUserData(userID string) {
	botMap.Delete(userID)
}
