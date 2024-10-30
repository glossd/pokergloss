package service

import "errors"

var blackList = map[string]struct{}{
	"AVD7AkuHstalWk6dAkRrx8dq2963":{}, // EdvarD
	"ti2RZW8ej4OulAluHnFKTFZsCbG3":{}, // Artur
}

var ErrUserBlackListed = errors.New("you are blacklisted")

func IsUserInBlackList(userID string) bool {
	_, ok := blackList[userID]
	return ok
}
