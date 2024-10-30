package auth

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/glossd/pokergloss/auth/authconf"
	"github.com/glossd/pokergloss/auth/authunsafe"
	"github.com/glossd/pokergloss/goconf"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func Init() {
	if authconf.JwtVerificationDisabled() {
		log.Println("WARN Jwt signature verification disabled")
	} else {
		ctx := context.Background()

		app, err := firebase.NewApp(ctx, &firebase.Config{
			ProjectID: goconf.Props.GCP.ProjectID,
		}, GoogleClientOptions()...)
		if err != nil {
			log.Fatalf("error initializing firebase: %v\n", err)
		}

		authunsafe.FirebaseClient, err = app.Auth(ctx)
		if err != nil {
			log.Fatalf("Couldn't create client: %s", err)
		}
	}
}

func GoogleClientOptions() []option.ClientOption {
	var opts []option.ClientOption
	if goconf.IsProd() {
		opts = append(opts, option.WithCredentialsJSON([]byte(goconf.Props.GCP.Credentials)))
	}
	return opts
}

func InitCustomSetup(client *auth.Client) {
	authunsafe.FirebaseClient = client
}
