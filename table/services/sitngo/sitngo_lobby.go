package sitngo

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services"
	"github.com/glossd/pokergloss/table/services/broadcast"
	"github.com/glossd/pokergloss/table/services/enrich"
	"github.com/glossd/pokergloss/table/services/events"
	"github.com/glossd/pokergloss/table/services/model"
	"github.com/glossd/pokergloss/table/services/paging"
	"github.com/glossd/pokergloss/table/services/player/actionhandler"
	"github.com/glossd/pokergloss/table/web/client/bankclient"
	"github.com/glossd/pokergloss/table/web/client/mqpub"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindAll(ctx context.Context, params paging.Params) ([]*domain.LobbySitAndGo, error) {
	var filter []bson.E
	filter = append(filter, bson.E{Key: "isprivate", Value: false})
	filter = append(filter, db.SkipEmptyFullFilter(params)...)
	return db.FilterSitngoLobbies(ctx, filter, db.PagingOptions(params))
}

func Create(ctx context.Context, params domain.NewLobbySitAndGoParams) (*model.LobbySitAndGo, error) {
	lobby, err := domain.NewLobbySitAndGo(params)
	if err != nil {
		return nil, err
	}

	err = db.InsertSitAndGoLobby(ctx, lobby)
	if err != nil {
		return nil, err
	}

	return model.ToSitAndGoLobby(lobby), nil
}

func Register(ctx context.Context, lobbyID string, position int, iden authid.Identity) error {
	oid, err := primitive.ObjectIDFromHex(lobbyID)
	if err != nil {
		return services.ErrInvalidIdFormat
	}

	lobby, err := db.FindSitAndGoLobby(ctx, oid)
	if err != nil {
		return err
	}

	err = lobby.Register(iden, position)
	if err != nil {
		return err
	}

	err = bankclient.Withdraw(ctx, lobby.BuyIn, iden.UserId, "Registered for Sit&Go tournament")
	if err != nil {
		return err
	}

	err = db.UpdateSitAndGoLobby(ctx, lobby)
	if err != nil {
		err := bankclient.Deposit(lobby.BuyIn, iden.UserId, "Failed to register for Sit&Go")
		if err != nil {
			log.Errorf("IMPORTANT! Failed to deposit back: user %s lost %d", iden, lobby.BuyIn)
		}
		if err == db.ErrVersionNotMatch {
			return services.ErrFormat("please try to register again")
		}
		return err
	}

	broadcast.SendTableEvent(lobbyID, &events.TableEvent{
		Type:    events.SitngoRegisterType,
		Payload: events.M{"position": position, "player": iden},
	})

	handleTable(ctx, lobby)

	return nil
}

func handleTable(ctx context.Context, l *domain.LobbySitAndGo) {
	if l.GetTable() != nil {
		err := db.InsertTable(ctx, l.GetTable())
		if err != nil {
			log.Errorf("Couldn't create sitngo table after last registered player ")
			return
		}

		// todo case of saving table error?

		// little hack, decision timeout sends key with incremented GameFlowVersion
		// but InsertTable does not increment GameFlowVersion
		l.GetTable().GameFlowVersion--
		actionhandler.LaunchDecisionTimeout(l.GetTable())

		enrich.Players(l.GetTable(), l.GetTable().AllPlayers())

		userIDs := make([]string, 0, len(l.Entries))
		for _, entry := range l.Entries {
			userIDs = append(userIDs, entry.UserId)
		}
		mqpub.SendNewsToUsers(userIDs, &events.TableEvent{
			Type: events.SitngoGameStartType,
			Payload: events.M{
				"table":   events.M{"id": l.GetTable().ID},
				"tableId": l.GetTable().ID.Hex(),
				"name":    l.Name,
			},
		})
	}
}

func Unregister(ctx context.Context, lobbyID string, position int, iden authid.Identity) error {
	oid, err := primitive.ObjectIDFromHex(lobbyID)
	if err != nil {
		return services.ErrInvalidIdFormat
	}

	lobby, err := db.FindSitAndGoLobby(ctx, oid)
	if err != nil {
		return err
	}

	err = lobby.Unregister(iden, position)
	if err != nil {
		return err
	}

	err = db.UpdateSitAndGoLobby(ctx, lobby)
	if err != nil {
		if err == db.ErrVersionNotMatch {
			return services.ErrFormat("please try to register again")
		}
		return err
	}

	err = bankclient.Deposit(lobby.BuyIn, iden.UserId, fmt.Sprintf("Left Sit & Go lobby %s", lobby.Name))
	if err != nil {
		log.Errorf("IMPORTANT! failed to send back %d to user %s", lobby.BuyIn, iden)
	}

	var wsEvents []*events.TableEvent
	wsEvents = append(wsEvents, &events.TableEvent{
		Type:    events.SitngoUnregisterType,
		Payload: events.M{"position": position, "player": iden},
	})

	broadcast.SendTableEvents(lobbyID, wsEvents)

	return nil
}
