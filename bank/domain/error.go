package domain

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Err struct {
	description string
}

func (e *Err) Error() string {
	return e.description
}

func (e Err) GRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.description)
}

func E(format string, a ...interface{}) *Err {
	return &Err{description: fmt.Sprintf(format, a...)}
}