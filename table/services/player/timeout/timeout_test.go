package timeout

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"testing"
)

func TestKey(t *testing.T) {
	timeoutMap := sync.Map{}
	oid := primitive.NewObjectID()
	timeoutMap.Store(Key{TableID: oid, Position: 0}, "my value")
	got, ok := timeoutMap.Load(Key{TableID: oid, Position: 0})
	assert.True(t, ok)
	assert.EqualValues(t, got, "my value")
}
