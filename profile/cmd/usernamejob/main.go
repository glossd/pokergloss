package main

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/profile/conf"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/service"
	"log"
	"strings"
)

func main() {
	conf.InitAuthClient()
	_, _, err := db.Init()
	if err != nil {
		log.Fatalf("db.Init: %s", err)
	}

	var nextPageToken string
	usersIter := conf.AuthClient.Users(context.Background(), nextPageToken)
	for {
		user, err := usersIter.Next()
		if err != nil {
			log.Fatalf("Failed fetch users: %s", err)
		}

		username := user.CustomClaims["username"].(string)
		newUsername := strings.ReplaceAll(username, "-", "_")
		if username != newUsername {
			err = service.ChangeUsername(context.Background(), newUsername, authid.Identity{
				UserId:   user.UID,
				Username: username,
			})
			if err != nil {
				log.Fatalf("Failed to change username: %s", err)
			}
			log.Printf("Changed from %s to %s\n", username, newUsername)
		} else {
			log.Printf("Skipping %s\n", username)
		}
	}
}
