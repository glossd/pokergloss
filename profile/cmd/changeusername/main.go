package main

import (
	"context"
	"flag"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/profile/conf"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/service"
	"log"
)

var username string
var newUsername string

func main() {
	flag.StringVar(&username, "old", "", "username of the user")
	flag.StringVar(&newUsername, "new", "", "new username of the user")
	flag.Parse()
	if username == "" {
		log.Fatalf("-old is not specified")
		return
	}
	if newUsername == "" {
		log.Fatalf("-new is not specified")
		return
	}

	conf.InitAuthClient()
	_, _, err := db.Init()
	if err != nil {
		log.Fatalf("db.Init: %s", err)
	}

	doc, err := db.FindProfileNoCtx(username)
	if err != nil {
		log.Fatalf("Find username: %s", err)
	}

	err = service.ChangeUsername(context.Background(), newUsername, authid.Identity{
		UserId:   doc.UserID,
		Username: doc.Username,
	})
	if err != nil {
		log.Fatalf("Failed to change username: %s", err)
	}
	log.Printf("Changed from %s to %s\n", username, newUsername)
}
