package domain

import "fmt"

type Err struct {
	description string
}

func (e *Err) Error() string {
	return e.description
}

func E(format string, a ...interface{}) *Err {
	return &Err{description: fmt.Sprintf(format, a...)}
}
