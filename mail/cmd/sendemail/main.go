package main

import (
	"github.com/glossd/pokergloss/mail/sender"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := sender.SendEmail("glossde@hotmail.com", "PokerGloss", "hello there!")
	if err != nil {
		log.Fatal(err)
	}
}
