package gogrpc

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

var ErrConNil = errors.New("grpc connection not initialized yet")

type SetCon func(con *grpc.ClientConn)

func InitConn(hostname string, port int, setCon SetCon) {
	if port == -1 {
		port = 9090
	}
	err := initConn(hostname, port, setCon)
	if err != nil {
		retryInitClient(hostname, port, setCon)
	}
}

func initConn(hostname string, port int, setCon SetCon) error {
	host := fmt.Sprintf("%s:%d", hostname, port)
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	con, err := grpc.DialContext(ctx, host, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Warnf("WARN Failed try connection to the server on %s: %s\n", host, err)
		return err
	}
	log.Infof("Successfully connected to the grpc server on %s:%d\n", host, port)
	setCon(con)
	return nil
}

func retryInitClient(hostname string, port int, setCon SetCon) {
	ticker := time.NewTicker(5 * time.Second)
	go func () {
		for {
			select {
			case <-ticker.C:
				err := initConn(hostname, port, setCon)
				if err == nil {
					ticker.Stop()
					return
				}
			}
		}
	}()
}