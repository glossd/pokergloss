package ws

import "sync"

// tableId -> TableHub
var tableConns = &sync.Map{}
var tableConnsCount int64 = 0

func getOrCreateTableHub(tableID string) *Hub {
	newHub := &Hub{}
	v, loaded := tableConns.LoadOrStore(tableID, newHub)
	if !loaded {
		newHub.initHub(tableID)
		go newHub.run()
	}
	return v.(*Hub)
}

func GetTableConnsCount() int64 {
	return tableConnsCount
}
