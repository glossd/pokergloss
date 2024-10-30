package e2e

import (
	"github.com/glossd/pokergloss/table-history/db"
	"github.com/glossd/pokergloss/table-history/domain"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"testing"
	"time"
)

func TestDeleteAllBefore(t *testing.T) {
	err := db.InsertManyEvents([]*domain.Event{
		{Type: "type", TableIDs: []string{"tableid"}, CreatedAt: bsonx.DateTime(time.Now().UnixNano() / 1e6), Payload: "payload"},
		{Type: "type", TableIDs: []string{"tableid"}, CreatedAt: bsonx.DateTime(time.Now().AddDate(0, 0, 2).UnixNano() / 1e6), Payload: "payload"},
	})
	assert.Nil(t, err)
	db.DeleteAllBefore(time.Now().AddDate(0, 0, 1))
	all, err := db.FindAll()
	assert.Nil(t, err)
	assert.EqualValues(t, 1, len(all))
}
