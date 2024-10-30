package domain

import (
	conf "github.com/glossd/pokergloss/goconf"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// cicd has env var RakePercent=0.01 and FeePercent=0.02
	conf.Props.Table.RakePercent = 0.0
	conf.Props.Tournament.FeePercent = 0.0
	code := m.Run()
	os.Exit(code)
}
