package e2e

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/profile/db"
	"github.com/glossd/pokergloss/profile/domain"
	"github.com/glossd/pokergloss/profile/web/router"
	"github.com/pokerblow/go-httptestutil"
	"github.com/pokerblow/mongotest"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var testRouter = httptestutil.NewRouter(router.New()).BasePath(router.BasePath)

var defaultToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjUxMDM2YWYyZDgzOWE4NDJhZjQzY2VjZmJiZDU4YWYxYTc1OGVlYTIiLCJ0eXAiOiJKV1QifQ.eyJ1c2VybmFtZSI6ImRlbmlzIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL3Bva2VyYmxvdyIsImF1ZCI6InBva2VyYmxvdyIsImF1dGhfdGltZSI6MTU5ODY0Njg3OCwidXNlcl9pZCI6IjJmUjVNWXlqcU1TTWd6THRka2RyS1prYzE4dDEiLCJzdWIiOiIyZlI1TVl5anFNU01nekx0ZGtkcktaa2MxOHQxIiwiaWF0IjoxNTk4NjQ2ODc4LCJleHAiOjE1OTg2NTA0NzgsImVtYWlsIjoiZGVuaXNnbG90b3Y5OEBtYWlsLnJ1IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7ImVtYWlsIjpbImRlbmlzZ2xvdG92OThAbWFpbC5ydSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.Kcn9QroIR-62xGlzGTHdvx2uNRfiqAtUJyYBfzg74Mt_v4XozZTW-6O_teFmFoRasJsOr49uW4i9ntkkgoc6FgDoo1jTi_1yMMx3_gNS9qSAMMcmMscqelOHQdgxsi9mJMwltHqHOf-AsoYl7qbc_HCf5ShYBtljZlkUXY_pMGvy0ePupNMFiWxoTYmNiIelaz0d-O9oVzns8XOm6O6A5qsFLx6hnNNsS7cBMMbc9zqhsySAZhMzYHdkd-LvL8QtUUjAQQfsbit9hPFa4irFEf7gfYOX61kUeKlMBP6-f6O8Q8GmgieBJ30Ly8wQwJcgHcYkTAFlZhD8FRO08QesgQ"
var secondToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjRlMDBlOGZlNWYyYzg4Y2YwYzcwNDRmMzA3ZjdlNzM5Nzg4ZTRmMWUiLCJ0eXAiOiJKV1QifQ.eyJuYW1lIjoiU2VhX01hbiIsInBpY3R1cmUiOiJodHRwczovL3N0b3JhZ2UuZ29vZ2xlYXBpcy5jb20vcG9rZXJibG93LWF2YXRhcnMvY0xoaHJDQzBzTVVkcElIR2xuZ29kbjFsS0xsMiIsInVzZXJuYW1lIjoiU2VhX01hbiIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9wb2tlcmJsb3ciLCJhdWQiOiJwb2tlcmJsb3ciLCJhdXRoX3RpbWUiOjE2MTY0MDM4NzQsInVzZXJfaWQiOiJjTGhockNDMHNNVWRwSUhHbG5nb2RuMWxLTGwyIiwic3ViIjoiY0xoaHJDQzBzTVVkcElIR2xuZ29kbjFsS0xsMiIsImlhdCI6MTYxNjQwMzg3NCwiZXhwIjoxNjE2NDA3NDc0LCJlbWFpbCI6InNlYW1hbkBwb2tlcmJsb3cuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZW1haWwiOlsic2VhbWFuQHBva2VyYmxvdy5jb20iXX0sInNpZ25faW5fcHJvdmlkZXIiOiJwYXNzd29yZCJ9fQ.GS18eIPAWd15FnGicQDiLcD5iAohFC19uhS3_heZ8WlB2_Z1fqBtjefxiDe4OsLzvbg9xHaOnyzyRim0YS1DEL-qjYu95dlLsc9hhKBUWCS8txhW7YFL88-W7cUpNv4zYOMZWchn7hvssa6zebwxzO7-RtRDI-8KXPOZEaU3qu5BaGBlr59R5Q2NfxX4E5VLAWAgfmekZk-jy53WMHrCPEMLFqmqLcAfVchz63PwczBARJsNwyPjQkswNAU-aXGTYvzi_O4uedgtRQloh3jdQgrlK1r_X2CSwbz8WO-2sze__oXuo8MbPzgK2-QFGvtJ9s7mA9uaI-PfnRXaqtASGg"

var authTestRouter = httptestutil.NewRouter(router.New()).BasePath(router.BasePath).Headers(authHeaders(defaultToken))

var defaultIdentity = authid.Identity{UserId: "2fR5MYyjqMSMgzLtdkdrKZkc18t1", Username: "denis"}
var secondIdentity = authid.Identity{UserId: "cLhhrCC0sMUdpIHGlngodn1lKLl2", Username: "Sea_Man"}

var defaultTechInfo = domain.TechInfo{}

func authHeaders(token string) map[string]string {
	return map[string]string{"Authorization": "Bearer " + token}
}

func cleanUpDB() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := db.ProfileCol().DeleteMany(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
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

func createUser(t *testing.T) {
	body := "username=den&email=den@mail.ru&password=123456"
	rr := testRouter.Request(t, http.MethodPost, "/signup", &body, map[string]string{"Content-type": "application/x-www-form-urlencoded"})
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}

func createUserWithUsername(t *testing.T, username string) {
	body := fmt.Sprintf("username=%s&email=den@mail.ru&password=123456", username)
	rr := testRouter.Request(t, http.MethodPost, "/signup", &body, map[string]string{"Content-type": "application/x-www-form-urlencoded"})
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}

func createUserFull(t *testing.T, username, email, password string) {
	body := fmt.Sprintf("username=%s&email=%s&password=%s", username, email, password)
	rr := testRouter.Request(t, http.MethodPost, "/signup", &body, map[string]string{"Content-type": "application/x-www-form-urlencoded"})
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}
