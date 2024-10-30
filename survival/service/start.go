package service

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/domain"
	"github.com/glossd/pokergloss/survival/web/client/botsquadclient"
	"github.com/glossd/pokergloss/survival/web/client/rpc"
)

var ErrNotEnoughTickets = domain.E("Not enough tickets")

type StartResult struct {
	TableID        string
	AlreadyStarted bool
}

func Start(ctx context.Context, iden authid.Identity, params domain.Params) (*StartResult, error) {
	if s, err := db.Find(ctx, iden.UserId); err == nil {
		if s.WoF == nil {
			return &StartResult{TableID: s.TableID, AlreadyStarted: true}, nil
		} else {
			// no prize for user
			err = db.Delete(ctx, iden.UserId)
			if err != nil {
				return nil, err
			}
		}
	}

	s := domain.New(iden, params)
	if params.Anonymous {
		count, err := db.FindAnonymousCounter(ctx)
		if err != nil {
			return nil, err
		}
		if count > conf.Props.MaxAnonymousCounter {
			return nil, ErrMaxAnonymousReached
		}
		err = db.IncAnonymousCounter(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		if !params.Idle {
			deced, err := DecCardTickets(ctx, iden.UserId)
			if err != nil {
				return nil, err
			}
			if !deced {
				return nil, ErrNotEnoughTickets
			}
		}
	}

	restore := func() {
		if params.Anonymous {
			_ = db.DecAnonymousCounter(ctx)
		} else if !params.Idle {
			_ = db.CardIncTicket(ctx, iden.UserId)
		}
	}

	err := createTableAndSquad(ctx, s)
	if err != nil {
		restore()
		return nil, err
	}

	err = db.Insert(ctx, s)
	if err != nil {
		restore()
		_ = botsquadclient.Delete(ctx, s.TableID, iden.UserId)
		_ = rpc.DeleteTable(ctx, s.TableID)
		return nil, err
	}
	return &StartResult{TableID: s.TableID}, nil
}

func createTableAndSquad(ctx context.Context, s *domain.Survival) error {
	if conf.IsE2E() {
		s.TableID = "e2eid"
		return nil
	}
	tableID, err := rpc.CreateTable(ctx, s.GetTableParams())
	if err != nil {
		return err
	}

	err = botsquadclient.Create(ctx, s.Iden, tableID, s.GetTableParams())
	if err != nil {
		_ = rpc.DeleteTable(ctx, tableID)
		return err
	}
	s.TableID = tableID
	return err
}
