package e2e

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/messenger/db"
	"github.com/glossd/pokergloss/messenger/web/router"
	"github.com/pokerblow/go-httptestutil"
	"github.com/pokerblow/mongotest"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

var testRouter = httptestutil.NewRouter(router.New(gin.New())).BasePath(router.BasePath).Headers(authHeaders(firstToken))

func authHeaders(token string) map[string]string {
	if token == "" {
		return map[string]string{"Authorization": ""}
	}
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

func cleanUp() {
	ctx := context.Background()
	err := db.MessageCol().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = db.UserChatListCol().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = db.ChatCol().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

var firstIdentity = authid.Identity{UserId: "2fR5MYyjqMSMgzLtdkdrKZkc18t1", Username: "denis"}
var secondIdentity = authid.Identity{UserId: "gY5mfr38daQXLuOxC6mPuz0LS6a2", Username: "denis2"}

var firstToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjUxMDM2YWYyZDgzOWE4NDJhZjQzY2VjZmJiZDU4YWYxYTc1OGVlYTIiLCJ0eXAiOiJKV1QifQ.eyJ1c2VybmFtZSI6ImRlbmlzIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL3Bva2VyYmxvdyIsImF1ZCI6InBva2VyYmxvdyIsImF1dGhfdGltZSI6MTU5ODY0Njg3OCwidXNlcl9pZCI6IjJmUjVNWXlqcU1TTWd6THRka2RyS1prYzE4dDEiLCJzdWIiOiIyZlI1TVl5anFNU01nekx0ZGtkcktaa2MxOHQxIiwiaWF0IjoxNTk4NjQ2ODc4LCJleHAiOjE1OTg2NTA0NzgsImVtYWlsIjoiZGVuaXNnbG90b3Y5OEBtYWlsLnJ1IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7ImVtYWlsIjpbImRlbmlzZ2xvdG92OThAbWFpbC5ydSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.Kcn9QroIR-62xGlzGTHdvx2uNRfiqAtUJyYBfzg74Mt_v4XozZTW-6O_teFmFoRasJsOr49uW4i9ntkkgoc6FgDoo1jTi_1yMMx3_gNS9qSAMMcmMscqelOHQdgxsi9mJMwltHqHOf-AsoYl7qbc_HCf5ShYBtljZlkUXY_pMGvy0ePupNMFiWxoTYmNiIelaz0d-O9oVzns8XOm6O6A5qsFLx6hnNNsS7cBMMbc9zqhsySAZhMzYHdkd-LvL8QtUUjAQQfsbit9hPFa4irFEf7gfYOX61kUeKlMBP6-f6O8Q8GmgieBJ30Ly8wQwJcgHcYkTAFlZhD8FRO08QesgQ"
var secondPlayerToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjUxMDM2YWYyZDgzOWE4NDJhZjQzY2VjZmJiZDU4YWYxYTc1OGVlYTIiLCJ0eXAiOiJKV1QifQ.eyJ1c2VybmFtZSI6ImRlbmlzMiIsImlzcyI6Imh0dHBzOi8vc2VjdXJldG9rZW4uZ29vZ2xlLmNvbS9wb2tlcmJsb3ciLCJhdWQiOiJwb2tlcmJsb3ciLCJhdXRoX3RpbWUiOjE1OTkwMzg1MzMsInVzZXJfaWQiOiJnWTVtZnIzOGRhUVhMdU94QzZtUHV6MExTNmEyIiwic3ViIjoiZ1k1bWZyMzhkYVFYTHVPeEM2bVB1ejBMUzZhMiIsImlhdCI6MTU5OTA1ODI3MCwiZXhwIjoxNTk5MDYxODcwLCJlbWFpbCI6ImRlbmlzZ2xvdG92LjE5MTFAbWFpbC5ydSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJlbWFpbCI6WyJkZW5pc2dsb3Rvdi4xOTExQG1haWwucnUiXX0sInNpZ25faW5fcHJvdmlkZXIiOiJwYXNzd29yZCJ9fQ.VbWRAmmcF_WIlbo2AeQ3a0yuRvRpwDm502Opt9rdn-sq8DbmK0XcUOOBcm5n0r8iKQUxXtEbHXDyJPdaNuARt6zsU29Ij-_EoQAvDLGPCpaAmch1GkqjMwmYeJILQtz_cG6R3SJCOmzzHouqIgx6kfI9PJNwqhHOHCU3Vp3GsMIZT11-HAKmvYwU0g8rEaSnfOxMsT6LypomgcjiwXU-vf6yoB4qnbtjxpp8lTHRU24kYCR5rFwP_NjVXRwFJ6-vDxgsfESS4a_Vk6u29C19KMKhsao7ilerbpG2XTU44wEi3AZfZhDkKvKxRUPWlICV7OK3zSywI9a7DL-vj76g_w"
