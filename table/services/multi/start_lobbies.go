package multi

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/enrich"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	log "github.com/sirupsen/logrus"
	"strings"
)

func LaunchMultiLobbies(startAt int64) error {
	if conf.IsE2E() {
		ctx := context.Background()
		err := StartMultiLobbies(ctx, startAt)
		if err != nil {
			return err
		}
		return nil
	}
	mqpub.PublishStartMulti(startAt)
	return nil
}

func StartMultiLobbies(ctx context.Context, startAt int64) error {
	lobbies, err := db.FindMultiLobbiesByStartAt(ctx, startAt)
	if err != nil {
		log.Errorf("Failed to start freerolls, find failure, startAt=%d: %s", startAt, err)
		return err
	}
	for _, lobby := range lobbies {
		err := startLobby(ctx, lobby)
		if err != nil {
			return err
		}
	}
	return nil
}

func startLobby(ctx context.Context, lobby *domain.LobbyMulti) error {
	if lobby.Status == domain.LobbyStarted {
		return nil
	}

	lobby.Start()

	if lobby.Status == domain.LobbyFinished {
		_ = db.DeleteLobbyMulti(ctx, lobby.ID)
		return nil
	}

	// using upsert instead of insert, because daily tournament has the same lobbyID
	err := db.UpsertRebalanceConfig(ctx, &db.RebalancerConfig{LobbyID: lobby.ID})
	if err != nil {
		return err
	}

	err = db.InsertManyTables(ctx, lobby.GetTables())
	if err != nil {
		// case where update lobby failed
		if !strings.Contains(err.Error(), "duplicate key") {
			log.Errorf("Failed to start multi, insert freeroll tables: %s", err)
			return err
		}
	}

	err = db.UpdateLobbyMulti(ctx, lobby)
	if err != nil {
		log.Errorf("Failed to start multi, update failure, startAt=%d: %s", lobby.StartAt, err)
		return err
	}

	mqpub.PublishMultiRebalanceStart(lobby.ID.Hex())
	sendGameStartEventToPlayersDirectly(lobby)
	sendMultiLobby(lobby)
	launchDecisionTimeouts(lobby)

	for _, table := range lobby.GetTables() {
		enrich.Players(table, table.AllPlayers())
	}

	log.Infof("Successfully started multi lobby: %s", lobby.ID.Hex())
	return nil
}

func launchDecisionTimeouts(lobby *domain.LobbyMulti) {
	for _, table := range lobby.GetTables() {
		// todo fix. little hack, LaunchDecisionTimeout increments game flow version
		table.GameFlowVersion--
		actionhandler.LaunchDecisionTimeout(table)
	}
}

func sendGameStartEventToPlayersDirectly(lobby *domain.LobbyMulti) {
	for _, t := range lobby.GetTables() {
		players := t.AllPlayers()
		userIDs := make([]string, 0, len(players))
		for _, p := range players {
			userIDs = append(userIDs, p.UserId)
		}
		mqpub.SendNewsToUsers(userIDs, &events.TableEvent{Type: events.MultiGameStartType, Payload: events.M{
			"tableId": t.ID.Hex(),
			"name":    lobby.Name,
		}})
	}
}

func sendMultiLobby(lobby *domain.LobbyMulti) {
	mqpub.SendTableMessage(lobby.ID.Hex(), []*events.TableEvent{{
		Type: events.MultiLobby,
		Payload: events.M{
			"lobby": model.ToLobbyMulti(lobby),
		}}})
}
