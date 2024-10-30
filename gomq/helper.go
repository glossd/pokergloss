package gomq

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

type M map[string]interface{}

func (m M) JSON() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Errorf("M.JSON couldn't marshal: %v", err)
		return "{}"
	}
	return string(bytes)
}
