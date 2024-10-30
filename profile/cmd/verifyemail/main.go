package main

import (
	"context"
	"firebase.google.com/go/auth"
	"flag"
	"github.com/glossd/pokergloss/profile/conf"
	"log"
	"time"
)

var email string

func main() {
	flag.StringVar(&email, "e", "", "email of the user")
	flag.Parse()
	if email == "" {
		log.Fatalf("email is not specified")
		return
	}

	conf.InitAuthClient()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := conf.AuthClient.GetUserByEmail(ctx, email)
	if err != nil {
		log.Fatalf("Failed fetch user: %s", err)
		return
	}

	params := (&auth.UserToUpdate{}).
		EmailVerified(true)

	_, err = conf.AuthClient.UpdateUser(ctx, user.UID, params)
	if err != nil {
		log.Fatalf("Failed verify email: %s", err)
		return
	}

	log.Println("Successfully verified email")
}
