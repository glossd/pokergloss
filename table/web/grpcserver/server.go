package grpcserver

import (
	"context"
	"encoding/json"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/gogrpc/grpctable"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/services/player/timeout"
	"github.com/glossd/pokergloss/table/services/tables"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// This grpc server seems to be used only for survival

func GetTable(ctx context.Context, r *grpctable.GetTableRequest) (*grpctable.GetTableResponse, error) {
	t, err := tables.Find(ctx, r.TableId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	modelTable := model.ToModelTable(t, model.ToPlayerOpenCards)
	bytes, err := json.Marshal(modelTable)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &grpctable.GetTableResponse{TableJson: bytes}, nil
}

func CreateSurvivalTable(ctx context.Context, r *grpctable.CreateSurvivalTableRequest) (*grpctable.CreateSurvivalTableResponse, error) {
	u := domain.User{
		Iden: authid.Identity{
			UserId:   r.User.UserId,
			Username: r.User.Username,
			Picture:  r.User.Picture,
		},
		Stack: r.User.Stack,
	}
	bots := make([]domain.Bot, 0, len(r.Bots))
	for _, bot := range r.Bots {
		bots = append(bots, domain.Bot{
			Name:    bot.Name,
			Picture: bot.Picture,
			Stack:   bot.Stack,
		})
	}
	params := domain.NewSurvivalTableParams{
		Name:              r.Name,
		BigBlind:          r.BigBlind,
		DecisionTimeSec:   r.DecisionTimeSec,
		ThemeID:           domain.ThemeID(r.ThemeId),
		User:              u,
		Bots:              bots,
		LevelIncreaseTime: time.Duration(r.LevelIncreaseTimeSec) * time.Second,
		SurvivalLevel:     r.SurvivalLevel,
	}
	table, err := domain.NewTableSurvival(params)
	if err != nil {
		return nil, err
	}
	err = db.InsertTable(ctx, table)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	mqpub.PublishTimeoutEvent(&timeout.Event{
		Type: timeout.StartGame,
		At:   timeutil.NowAdd(2500 * time.Millisecond),
		Key:  timeout.Key{TableID: table.ID, Position: -1, Version: 0},
	})

	modelTable := model.ToModelTable(table, model.ToPlayerNoCards)
	tableBytes, err := json.Marshal(modelTable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &grpctable.CreateSurvivalTableResponse{TableJson: tableBytes}, nil
}

func MakeAction(ctx context.Context, r *grpctable.MakeActionRequest) (*empty.Empty, error) {
	t, err := tables.Find(ctx, r.TableId)
	if err != nil {
		return nil, err
	}

	err = t.MakeBotAction(int(r.Position), domain.Action{Type: domain.ActionType(r.ActionType), Chips: r.Chips})
	if err != nil {
		return nil, err
	}
	err = actionhandler.Handle(ctx, t)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func DeleteSurvivalTable(ctx context.Context, r *grpctable.DeleteSurvivalTableRequest) (*empty.Empty, error) {
	err := tables.Delete(ctx, r.TableId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func SitBack(ctx context.Context, r *grpctable.SitBackRequest) (*empty.Empty, error) {
	t, err := tables.Find(ctx, r.TableId)
	if err != nil {
		return nil, err
	}

	err = t.SitBotBack(int(r.Position))
	if err != nil {
		return nil, err
	}
	err = db.SetTableContext(ctx, t.ID, db.PlayerUpdateStatus(t.GetSeatUnsafe(int(r.Position))))
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
