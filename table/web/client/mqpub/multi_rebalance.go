package mqpub

import (
	"github.com/glossd/memmq"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/web/client/mq"
	log "github.com/sirupsen/logrus"
)

func PublishMultiRebalanceStart(lobbyID string) {
	publishMultiRebalance(lobbyID, "", 0)
}

func PublishMultiRebalance(lobbyID string, tableID string) {
	publishMultiRebalance(lobbyID, tableID, 0)
}

func PublishMultiRebalanceAt(lobbyID string, at int64) {
	publishMultiRebalance(lobbyID, "", at)
}

func publishMultiRebalance(lobbyID string, tableID string, at int64) {
	event := mq.MultiPlayersMovedEvent{LobbyID: lobbyID, TableID: tableID, RebalanceAt: at}
	if conf.IsE2E() {
		mq.TestMultiPlayersMovedQueue <- &event
		return
	}
	err := memmq.Publish(mq.MultiPlayersMovedTopicID, &event)
	if err != nil {
		log.Errorf("Failed to send multi rebalance event: %s", err)
	}
}
