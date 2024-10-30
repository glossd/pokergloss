package cleaning

import (
	"context"
	"github.com/glossd/pokergloss/table/db"
	"github.com/glossd/pokergloss/table/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CleanFinishedLobbies() {
	_, _ = CleanFinishedLobbiesOf(db.ColSitAndGoLobby())
	_, _ = CleanFinishedLobbiesOf(db.ColLobbyMulti())
}

func CleanFinishedLobbiesOf(collection *mongo.Collection) (int64, error) {
	many, err := collection.DeleteMany(context.Background(), bson.M{"status": domain.LobbyFinished})
	if err != nil {
		log.Errorf("Failed to clean finished %s lobbies: %s", collection.Name(), err)
		return 0, err
	}
	if many.DeletedCount > 0 {
		log.Infof("Finished cleaning finished %s lobbies, count=%d", collection.Name(), many.DeletedCount)
	}
	return many.DeletedCount, nil
}
