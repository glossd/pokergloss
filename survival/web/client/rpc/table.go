package rpc

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gogrpc/grpctable"
	"github.com/glossd/pokergloss/survival/domain"
	grpc "github.com/glossd/pokergloss/table/web/grpcserver"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func CreateTable(ctx context.Context, params domain.TableParams) (string, error) {

	var bots []*grpctable.Bot
	for _, bot := range params.Bots {
		bots = append(bots, &grpctable.Bot{
			Name:    bot.Name,
			Picture: bot.Picture,
			Stack:   bot.Stack,
		})
	}
	resp, err := grpc.CreateSurvivalTable(ctx, &grpctable.CreateSurvivalTableRequest{
		Name:                 params.Name,
		BigBlind:             params.BigBlind,
		DecisionTimeSec:      params.DecisionTimeoutSec,
		ThemeId:              params.ThemeID,
		LevelIncreaseTimeSec: params.LevelIncreaseTimeSec,
		User: &grpctable.User{
			UserId:   params.User.UserId,
			Username: params.User.Username,
			Picture:  params.User.Picture,
			Stack:    params.UserStack,
		},
		SurvivalLevel: params.SurvivalLevel,
		Bots:          bots,
	})
	if err != nil {
		log.Errorf("Failed to create table: %s", err)
		return "", err
	}
	tableID := gjson.GetBytes(resp.TableJson, "id").String()

	return tableID, nil
}

func DeleteTable(ctx context.Context, tableID string) error {
	if conf.IsE2E() {
		return nil
	}
	_, err := grpc.DeleteSurvivalTable(ctx, &grpctable.DeleteSurvivalTableRequest{TableId: tableID})
	if err != nil {
		log.Errorf("rpc.DeleteTable failed: %s", err)
	}
	return nil
}
