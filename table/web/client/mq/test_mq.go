package mq

import (
	"github.com/glossd/pokergloss/gomq/mqtable"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table/services/player/timeout"
)

var TestMQ = make(chan *mqws.TableMessage, 64)
var TestNewsMQ = make(chan *mqws.Message, 64)
var TestGameEndMQ = make(chan *mqtable.GameEnd, 16)
var TimeoutTestMQ = make(chan *timeout.Event, 64)
var IsTimeoutTestMQEnabled = false
var TestMultiPlayersMovedQueue = make(chan *MultiPlayersMovedEvent, 64)
var IsMultiPlayersMovedEnabledTest = false

func ResetTestMQ() {
	TestMQ = make(chan *mqws.TableMessage, 64)
	TestNewsMQ = make(chan *mqws.Message, 64)
	TestGameEndMQ = make(chan *mqtable.GameEnd, 16)
	IsTimeoutTestMQEnabled = true
	for len(TimeoutTestMQ) > 0 {
		<-TimeoutTestMQ
	}
	for len(TestMultiPlayersMovedQueue) > 0 {
		<-TestMultiPlayersMovedQueue
	}
}
