package gomq

type ackableError struct {
	message string
}

func NewAckableError(msg string) error {
	return &ackableError{
		message: msg,
	}
}

func WrapInAckableError(err error) error {
	if err == nil {
		return &ackableError{message: "nil"}
	}
	return &ackableError{message: err.Error()}
}

func IsAckableError(err error) bool {
	if err == nil {
		return true
	}
	_, ok := err.(*ackableError)
	return ok
}

func (e *ackableError) Error() string {
	return e.message
}
