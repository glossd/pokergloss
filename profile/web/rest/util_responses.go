package rest

// Defines http error response
type ErrorRes struct {
	Message string `json:"message"`
}

func E(err error) *ErrorRes {
	return &ErrorRes{Message: err.Error()}
}

func Efmt(description string) *ErrorRes {
	return &ErrorRes{Message: description}
}


type OkRes struct {
	Message string `json:"message"`
}

func Ok(msg string) *OkRes {
	return &OkRes{Message: msg}
}
