package services

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidIdFormat = ErrFormat("invalid id format")

type ServiceErr struct {
	description string
}

func (e *ServiceErr) Error() string {
	return e.description
}

func (e *ServiceErr) GRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.description)
}

func ErrFormat(format string, a ...interface{}) *ServiceErr {
	return &ServiceErr{description: fmt.Sprintf(format, a...)}
}

func E(err error) *ServiceErr {
	if err == nil {
		return nil
	}
	return &ServiceErr{description: err.Error()}
}
