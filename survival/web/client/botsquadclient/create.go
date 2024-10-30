package botsquadclient

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	botsquad "github.com/glossd/pokergloss/survival/bot-squad"
	botsquadconf "github.com/glossd/pokergloss/survival/bot-squad/conf"
	"github.com/glossd/pokergloss/survival/domain"
	log "github.com/sirupsen/logrus"
	"strings"
)

func Create(ctx context.Context, forUser authid.Identity, tableID string, params domain.TableParams) error {
	var positions string
	var weakPositions string
	var looseness string
	var aggression string
	for i, bot := range params.Bots {
		positions = addIntStr(positions, i+1)
		weakPositions = addIntStr(weakPositions, i+1)
		looseness = addFloatStr(looseness, bot.Looseness)
		aggression = addFloatStr(aggression, bot.Aggression)
	}

	err := botsquad.Create(botsquadconf.Config{
		UserID:  forUser.UserId,
		TableID: tableID,
		Squad: botsquadconf.Squad{
			Positions:  strings.Split(positions, " "),
			Looseness:  strings.Split(looseness, " "),
			Aggression: strings.Split(aggression, " "),
		},
	})
	if err != nil {
		log.Errorf("Failed to create bot-squad: %v", err)
		return err
	}
	return nil
}

func addIntStr(acc string, v int) string {
	if acc == "" {
		return fmt.Sprintf("%d", v)
	} else {
		return fmt.Sprintf("%s %d", acc, v)
	}
}

func addFloatStr(acc string, v float64) string {
	if acc == "" {
		return fmt.Sprintf("%.2f", v)
	} else {
		return fmt.Sprintf("%s %.2f", acc, v)
	}
}

func Delete(ctx context.Context, tableID string, userID string) error {
	if conf.IsE2E() {
		return nil
	}
	err := botsquad.Delete(userID)
	if err != nil {
		log.Errorf("Failed to delete bot squad, userID=%s tableID=%s: %v", userID, tableID, err)
		return err
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
