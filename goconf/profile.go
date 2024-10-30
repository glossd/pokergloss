package goconf

import "os"

var Profile = os.Getenv("PG_PROFILE")
var IsE2EVar = Profile == "e2e"
var IsLocalVar = Profile == "local"
var IsDevVar = Profile == "dev"

func IsProd() bool {
	return !IsLocal()
}

func IsDev() bool {
	return IsDevVar
}

func IsLocal() bool {
	return IsE2EVar || IsLocalVar
}

func IsLocalOnly() bool {
	return IsLocalVar
}

func IsE2E() bool {
	return IsE2EVar
}
