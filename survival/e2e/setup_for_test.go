package e2e

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/web/router"
	"github.com/pokerblow/go-httptestutil"
	"github.com/pokerblow/mongotest"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

var defaultToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjhmNDMyMDRhMTc5MTVlOGJlN2NjZDdjYjI2NGRmNmVhMzgzYzQ5YWIiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoicG9rZXIiLCJwaWN0dXJlIjoiaHR0cHM6Ly9zdG9yYWdlLmdvb2dsZWFwaXMuY29tL3Bva2VyYmxvdy1hdmF0YXJzL3NZbmRSa01GMlZNSFNHR2JIRXM1UWozczZqazItWGRweiIsInVzZXJuYW1lIjoicG9rZXIiLCJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vcG9rZXJibG93IiwiYXVkIjoicG9rZXJibG93IiwiYXV0aF90aW1lIjoxNjI1NTAzMjkxLCJ1c2VyX2lkIjoic1luZFJrTUYyVk1IU0dHYkhFczVRajNzNmprMiIsInN1YiI6InNZbmRSa01GMlZNSFNHR2JIRXM1UWozczZqazIiLCJpYXQiOjE2MjYwMDc2MjYsImV4cCI6MTYyNjAxMTIyNiwiZW1haWwiOiJkZW5pc2dsb3Rvdi4xOTExQG1haWwucnUiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJlbWFpbCI6WyJkZW5pc2dsb3Rvdi4xOTExQG1haWwucnUiXX0sInNpZ25faW5fcHJvdmlkZXIiOiJwYXNzd29yZCJ9fQ.Vy9wUEIuO52mneQwZymI9B9pQ5xRBkUu1NRuREyKfzPoUks2DSjcnxtR7ysb9buVJz6h93rR5Ns88IN2oOPUppj6AN7gXHqnyH3q-SLPjlzS8Ri7dt9csUh6CR56k0r493U-K9Okp7gQJXufXGhWHTfcZW0aCc0o-QiuP8whNCoKx9n-Gt54y8pIfjgQcclytWE18tX_sduWBFUS_A_dkyFGsy2fKbBvwuT311QHvijo-KN3BKX0ZvRsPVnwXMwRv0dQPJHKZtori8WH8_55c87mTDLXAAA-_P_IcNKoDxNl68qdWa42u3qNZ0qVF4qJHmbGEC7qAbO1dO4KU8RWkQ"
var defaultIdentity = authid.Identity{UserId: "sYndRkMF2VMHSGGbHEs5Qj3s6jk2", Username: "poker", Picture: "https://storage.googleapis.com/pokerblow-avatars/sYndRkMF2VMHSGGbHEs5Qj3s6jk2-Xdpz"}

var testRouter = httptestutil.NewRouter(router.New(gin.New())).BasePath(router.BasePath).Headers(authHeaders(defaultToken))

func authHeaders(token string) map[string]string {
	return map[string]string{"Authorization": "Bearer " + token}
}

func TestMain(m *testing.M) {
	conf.IsE2EVar = true
	os.Setenv("PB_JWT_VERIFICATION_DISABLE", "true")
	auth.Init()

	cc := mongotest.StartMongoContainer("4.2")

	_, _, err := db.InitWithURI(cc.GetMongoURI(db.DbName))
	if err != nil {
		cc.KillMongoContainer()
		log.Fatal(err)
	}

	code := m.Run()

	cc.KillMongoContainer()

	os.Exit(code)
}

func cleanUpDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := db.SurvivalCol().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = db.CardCol().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = db.ScoreCol().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db.CreateCardsCol(ctx)
}
