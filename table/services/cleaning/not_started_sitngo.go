package cleaning

import (
	"context"
	"fmt"
	conf "github.com/glossd/pokergloss/goconf"
	"github.com/glossd/pokergloss/goconf/timeutil"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	"github.com/glossd/pokergloss/table/services/sitngo"
	"github.com/glossd/pokergloss/table/web/client/bankclient"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func CleanNotStartedLobbies() {
	deleteTime := timeutil.NowAdd(-conf.Props.Cleaning.SitngoStartTimeout)
	filter := bson.D{
		{"tournamentlobby.status", domain.LobbyRegistering},
		{"createdat", bson.M{"$lte": deleteTime}},
	}
	err := db.ForEachSitngoLobby(filter, func(l *domain.LobbySitAndGo) {
		if _, ok := sitngo.PersistentSitNGoIDs[l.ID]; ok {
			return
		}
		var deposits []bankclient.UserDeposit
		for _, entry := range l.Entries {
			deposits = append(deposits, bankclient.UserDeposit{
				UserID:      entry.UserId,
				Amount:      l.BuyIn,
				Description: fmt.Sprintf("SitNGo expired %s", l.Name),
			})
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := db.DeleteSitngo(ctx, l.ID)
		if err != nil {
			return
		}
		err = bankclient.MultiDeposit(ctx, deposits)
		if err != nil {
			return
		}
	})
	if err != nil {
		log.Errorf("Failed to clean not started sitngo lobbies: %s", err)
	}
}
