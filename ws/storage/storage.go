package storage

import (
	"sync"
)

// userId -> UserConn
var userConns = &sync.Map{}
var userConnsCount int64 = 0

func ReInitStorage() {
	userConns = &sync.Map{}
	userConnsCount = 0
}

func GetOrCreateUserHub(userID string) *UserHub {
	newHub := NewEmptyUserHub()
	v, loaded := userConns.LoadOrStore(userID, newHub)
	hub := v.(*UserHub)
	if !loaded {
		newHub.InitEmptyHub()
		go newHub.Run()
	}
	return hub
}

func GetUserHub(userID string) (*UserHub, bool) {
	load, ok := userConns.Load(userID)
	if ok {
		return load.(*UserHub), true
	}
	return nil, false
}

func DeleteUserClient(userID string) {
	userConns.Delete(userID)
}

func GetUserConnections() int64 {
	return userConnsCount
}