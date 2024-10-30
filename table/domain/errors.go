package domain

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PokerErr struct {
	description string
}

func (e *PokerErr) Error() string {
	return e.description
}

func (e *PokerErr) GRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.description)
}

func E(format string, a ...interface{}) *PokerErr {
	return &PokerErr{description: fmt.Sprintf(format, a...)}
}

