package fcm

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	goauth "github.com/glossd/pokergloss/auth"
	"log"
)

var FcmClient *messaging.Client

func Init() {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, goauth.GoogleClientOptions()...)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	FcmClient, err = app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error initializing fcm client: %v\n", err)
	}
}
