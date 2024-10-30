package e2e

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestListItems(t *testing.T) {
	rr := testRouter.Request(t, http.MethodGet, "/items", nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
}
