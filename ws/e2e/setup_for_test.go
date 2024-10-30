package e2e

import (
	"github.com/glossd/pokergloss/auth"
	conf "github.com/glossd/pokergloss/goconf"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	conf.IsE2EVar = true
	os.Setenv("PB_JWT_VERIFICATION_DISABLE", "true")
	auth.Init()
	code := m.Run()
	os.Exit(code)
}
