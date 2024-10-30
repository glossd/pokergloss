package conf

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	goauth "github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/goconf"
	log "github.com/sirupsen/logrus"
	"time"
)

const timeout = time.Second * 7

var AuthClient *auth.Client
var StorageClient *storage.Client

func InitAuthClient() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var err error
	app, err := firebase.NewApp(ctx, nil, goauth.GoogleClientOptions()...)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	firebaseCtx, authCancel := context.WithTimeout(context.Background(), timeout)
	defer authCancel()

	AuthClient, err = app.Auth(firebaseCtx)
	if err != nil {
		log.Fatalf("error initializing firebase auth: %v\n", err)
	}

	if goconf.IsProd() {
		StorageClient, err = app.Storage(firebaseCtx)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
