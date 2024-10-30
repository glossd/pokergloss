package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/glossd/pokergloss/gogrpc/grpctable"
	"github.com/glossd/pokergloss/survival/bot-squad/domain"
	grpc "github.com/glossd/pokergloss/table/web/grpcserver"
)

func GetTable(tableID string) (*domain.Table, error) {
	resp, err := grpc.GetTable(context.Background(), &grpctable.GetTableRequest{TableId: tableID})
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %s", err)
	}

	var t domain.Table
	err = json.Unmarshal(resp.TableJson, &t)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the table: %s", err)
	}
	return &t, nil
}
