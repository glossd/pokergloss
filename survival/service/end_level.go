package service

import (
	"context"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/gomq"
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/web/client/botsquadclient"
	"github.com/glossd/pokergloss/survival/web/client/rpc"
	"github.com/glossd/pokergloss/survival/web/mq/mqpub"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func EndLevel(ctx context.Context, userID, tableID string, isUserLost bool) error {
	err := botsquadclient.Delete(ctx, tableID, userID)
	if err != nil {
		return err
	}

	_ = rpc.DeleteTable(ctx, tableID)

	s, err := db.Find(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Errorf("No survival found to end level userID=%s", userID)
			return nil
		} else {
			return err
		}
	}

	if isUserLost {
		s.CreateWheelOfFortune()
		if s.IsAnonymous {
			err = db.DecAnonymousCounter(ctx)
			if err != nil {
				return err
			}
		} else if !s.IsIdle {
			err = db.UpsertScore(ctx, userID, s.Level)
			if err != nil {
				return err
			}
			err = db.Update(ctx, s)
			if err != nil {
				return err
			}
		}
		if s.GetWheelOfFortune() == nil {
			_ = mqpub.PublishEmptyWheel(s)
			err := db.Delete(ctx, s.UserID)
			if err != nil {
				return err
			}
		} else {
			_ = mqpub.PublishWheel(s)
		}
		return nil
	} else {
		s.NewLevel()
		err := createTableAndSquad(ctx, s)
		if err != nil {
			return err
		}
		err = db.Update(ctx, s)
		if err != nil {
			_ = botsquadclient.Delete(ctx, s.TableID, userID)
			_ = rpc.DeleteTable(ctx, s.TableID)
			return err
		}

		if !conf.IsE2E() {
			err = mqws.PublishNews(&mqws.Message{ToUserIds: []string{s.UserID}, Events: []*mqws.Event{{Type: "multiPlayerMove", Payload: gomq.M{"tableId": s.TableID}.JSON()}}})
			if err != nil {
				log.Errorf("Failed to publish multiMovePlayer: %s", err)
			}
		}
	}

	return nil
}
