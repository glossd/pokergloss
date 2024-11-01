package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/achievement"
	"github.com/glossd/pokergloss/assignment"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/bank"
	"github.com/glossd/pokergloss/bonus"
	browserlogs "github.com/glossd/pokergloss/browser-logs"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/mail"
	"github.com/glossd/pokergloss/market"
	"github.com/glossd/pokergloss/messenger"
	"github.com/glossd/pokergloss/profile"
	"github.com/glossd/pokergloss/survival"
	"github.com/glossd/pokergloss/table"
	tablechat "github.com/glossd/pokergloss/table-chat"
	tablehistory "github.com/glossd/pokergloss/table-history"
	tablestream "github.com/glossd/pokergloss/table-stream"
	"github.com/glossd/pokergloss/ws"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	auth.Init()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	shutdowns := []func(ctx context.Context){
		achievement.Run(r),
		assignment.Run(r),
		bank.Run(r),
		bonus.Run(r),
		browserlogs.Run(r),
		table.Run(r),
		mail.Run(r),
		market.Run(r),
		messenger.Run(r),
		profile.Run(r),
		survival.Run(r),
		tablechat.Run(r),
		tablehistory.Run(r),
		tablestream.Run(r),
		ws.Run(r),
	}
	runWithGracefulShutDown(r, shutdowns)
}

// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
func runWithGracefulShutDown(r *gin.Engine, serviceShutdowns []func(context.Context)) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Props.Port),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Infof("Starting server on %d", conf.Props.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infof("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	for _, shutdown := range serviceShutdowns {
		shutdown(ctx)
	}

	log.Infof("Server exiting")
}
