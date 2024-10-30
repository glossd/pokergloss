package idp

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	goauth "github.com/glossd/pokergloss/auth"
	log "github.com/sirupsen/logrus"
)

var client *auth.Client

func Init() {
	app, err := firebase.NewApp(context.Background(), nil, goauth.GoogleClientOptions()...)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
}
