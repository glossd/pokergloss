package service

import (
	"fmt"
)

var ErrMaxAnonymousReached = fmt.Errorf("anomynous survivals reached maximum")
