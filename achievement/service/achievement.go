package service

import (
	"context"
	"github.com/glossd/pokergloss/achievement/db"
	"github.com/glossd/pokergloss/achievement/domain"
	"github.com/glossd/pokergloss/achievement/model"
	"github.com/glossd/pokergloss/achievement/web/mq/bank"
	"github.com/glossd/pokergloss/achievement/web/mq/ws"
	"github.com/glossd/pokergloss/gomq/mqtable"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"sort"
	"time"
)

func FindSortedAchievements(ctx context.Context, userId string) ([]*model.Achievement, error) {
	as, err := FindOrBuildAS(ctx, userId)
	if err != nil {
		return nil, err
	}
	achievements := model.ToAchievements(as)
	sort.Slice(achievements, func(i, j int) bool {
		return achievements[i].Level > achievements[j].Level
	})
	return achievements, nil
}

func FindOrBuildAS(ctx context.Context, userID string) (*domain.AchievementStore, error) {
	as, err := db.FindAchievementStore(ctx, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			as = domain.NewAchievementStore(userID)
		} else {
			return nil, err
		}
	}
	return as, nil
}

func FindOrBuildASNoCtx(userID string) (*domain.AchievementStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return FindOrBuildAS(ctx, userID)
}

// For tests.
func UpdateAchievementStoreNoCtx(msg *mqtable.GameEnd) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	UpdateAchievementStore(ctx, msg)
}

func UpdateAchievementStore(ctx context.Context, msg *mqtable.GameEnd) {
	// it's either msg.Winners or msg.TournamentWinners

	for _, winner := range msg.Winners {
		as, err := FindOrBuildAS(ctx, winner.UserId)
		if err != nil {
			log.Errorf("UpdateHandCounter, find failed: %s", err)
			continue
		}

		ge := domain.NewGameEnd(msg)
		as.Update(winner, ge)

		err = db.UpsertAchievementStore(ctx, as)
		if err != nil {
			log.Errorf("Failed to save update %+v , for user %s: %s", winner, as.UserID, err)
			continue
		}

		timeToWait := time.Duration(ge.GameStartAt-time.Now().Unix()-1) * time.Second
		if as.HandsCounter.GetPrize().Chips > 0 {
			bank.DepositHandAchievement(as)
			time.AfterFunc(timeToWait, func() { ws.PublishHandEvent(as) })
		}

		for _, c := range []domain.Counter{as.WinCounter, as.BustCounter, as.DefeatCounter} {
			bank.DepositCounterAchievement(as.UserID, c)
			time.AfterFunc(timeToWait, func() { ws.PublishNewCounterAchievement(as.UserID, c) })
		}
	}
}

func UpdateAchievementStoreTournament(ctx context.Context, te *mqtable.TournamentEnd) {
	for _, winner := range te.TournamentWinners {
		as, err := FindOrBuildAS(ctx, winner.UserId)
		if err != nil {
			log.Errorf("UpdateHandCounter, find failed: %s", err)
			continue
		}

		switch te.Type {
		case mqtable.TournamentEnd_SITNGO:
			as.SitngoWinCounter.Inc()
		case mqtable.TournamentEnd_MULTI:
			as.MultiWinCounter.Inc()
		}

		err = db.UpsertAchievementStore(ctx, as)
		if err != nil {
			log.Errorf("Failed to save tournament winner %+v , for user %s: %s", winner, as.UserID, err)
			continue
		}

		if as.SitngoWinCounter.GetPrize().Chips > 0 {
			bank.DepositCounterAchievement(as.UserID, as.SitngoWinCounter)
			ws.PublishNewCounterAchievement(as.UserID, as.SitngoWinCounter)
		}
		if as.MultiWinCounter.GetPrize().Chips > 0 {
			bank.DepositCounterAchievement(as.UserID, as.MultiWinCounter)
			ws.PublishNewCounterAchievement(as.UserID, as.MultiWinCounter)
		}
	}
}
