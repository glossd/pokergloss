package e2e

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/assignment/db"
	"github.com/glossd/pokergloss/assignment/web/rest/router"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/pokerblow/go-httptestutil"
	"github.com/pokerblow/mongotest"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

var defaultToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjQ5YWQ5YmM1ZThlNDQ3OTNhMjEwOWI1NmUzNjFhMjNiNDE4ODA4NzUiLCJ0eXAiOiJKV1QifQ.eyJwaWN0dXJlIjoiaHR0cHM6Ly9zdG9yYWdlLmdvb2dsZWFwaXMuY29tL3Bva2VyYmxvdy1hdmF0YXJzLzJmUjVNWXlqcU1TTWd6THRka2RyS1prYzE4dDEiLCJ1c2VybmFtZSI6ImRlbmlzIiwiaXNzIjoiaHR0cHM6Ly9zZWN1cmV0b2tlbi5nb29nbGUuY29tL3Bva2VyYmxvdyIsImF1ZCI6InBva2VyYmxvdyIsImF1dGhfdGltZSI6MTU5OTE0NDc1OCwidXNlcl9pZCI6IjJmUjVNWXlqcU1TTWd6THRka2RyS1prYzE4dDEiLCJzdWIiOiIyZlI1TVl5anFNU01nekx0ZGtkcktaa2MxOHQxIiwiaWF0IjoxNTk5ODI0OTY2LCJleHAiOjE1OTk4Mjg1NjYsImVtYWlsIjoiZGVuaXNnbG90b3Y5OEBtYWlsLnJ1IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJmaXJlYmFzZSI6eyJpZGVudGl0aWVzIjp7ImVtYWlsIjpbImRlbmlzZ2xvdG92OThAbWFpbC5ydSJdfSwic2lnbl9pbl9wcm92aWRlciI6InBhc3N3b3JkIn19.uUFnfxWFGXAlPFjS9Lf5QjB_YLx824dqDjnOQi-xYquaCRjE9yMHZ-Tk8MKLhbjjWWWK0eleghfqIbZQmTowXhOkq9h5fCb5RgWUrRe2otRm5APhWrX-xXkkDpl5uEcYj9xly-1MHrCUBQdUCFaQdkwD960LJf-8saFmv_FWND-QveRyJxja10Dt8L6I6YI_NjoZTgYe24JamNNgxkFrfqFNN6fPN_X_SQYJXHetkX0zm_SL6UODwgId6ioTSOM0QrLt1ZcdGiq1eNuX8s9aQqHs8XYPRm4Bt0zARX3J4RReX9OZa64O9d26WDnLVz_QdPzzIz8itxpiZnYIjnGVRA"
var defaultIdentity = authid.Identity{UserId: "2fR5MYyjqMSMgzLtdkdrKZkc18t1", Username: "denis", Picture: "https://storage.googleapis.com/pokerblow-avatars/2fR5MYyjqMSMgzLtdkdrKZkc18t1"}

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

func cleanUp() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := db.ColUserDaily().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = db.ColDaily().Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}
	db.CleanCache()
}
