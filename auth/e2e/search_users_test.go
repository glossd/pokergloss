package e2e

import (
	"context"
	"github.com/stretchr/testify/assert"
"github.com/glossd/pokergloss/auth/authsearch"
"github.com/glossd/pokergloss/auth/authunsafe"
"testing"
"time"
)

const denisUserID = "2fR5MYyjqMSMgzLtdkdrKZkc18t1"
const denis2UserID = "gY5mfr38daQXLuOxC6mPuz0LS6a2"

func TestGetIdentities(t *testing.T) {
	authunsafe.FirebaseClient = initFirebaseClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	idens, err := authsearch.GetIdentities(ctx, []string{denisUserID, denis2UserID})
	assert.Nil(t, err)
	assert.Len(t, idens, 2)
	assert.EqualValues(t, "denis", idens[0].Username)
	assert.EqualValues(t, "denis2", idens[1].Username)
}
