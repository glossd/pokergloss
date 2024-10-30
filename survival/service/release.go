package service

import (
	"context"
	"github.com/glossd/pokergloss/auth/authid"
	"github.com/glossd/pokergloss/survival/db"
	"github.com/glossd/pokergloss/survival/web/mq/mqpub"
)

func Release(ctx context.Context, iden authid.Identity) error {
	s, err := db.Find(ctx, iden.UserId)
	if err != nil {
		return err
	}

	err = db.Delete(ctx, iden.UserId)
	if err != nil {
		return err
	}

	if s.IsIdle || s.IsAnonymous {
		return nil
	}

	if s.GetWheelOfFortune().WonSlot().Item != nil {
		_ = mqpub.GiftItem(s)
	} else {
		_ = mqpub.Deposit(s)
	}

	return nil
}
