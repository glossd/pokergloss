package domain

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"testing"
)

func TestMultiLobby_31thPlayerLeftWithoutTable(t *testing.T) {
	data, err := ioutil.ReadFile("./snapshots/multi-lobby-13-03-2021.json")
	assert.Nil(t, err)
	var l LobbyMulti
	assert.Nil(t, bson.UnmarshalExtJSON(data, true, &l))
	oid, err := primitive.ObjectIDFromHex("5fea7180df0ccde21ef5e2c0")
	assert.Nil(t, err)
	l.ID = oid

	l.Start()
	assert.EqualValues(t, 6, len(l.GetTables()))
	assert.EqualValues(t, LobbyStarted, l.Status)
}
