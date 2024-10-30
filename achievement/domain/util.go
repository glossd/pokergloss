package domain

import "fmt"

type Error struct {
	description string
}

func newEf(descr string, a ...interface{}) *Error {
	return &Error{description: fmt.Sprintf(descr, a...)}
}

func (e *Error) Error() string {
	return e.description
}
