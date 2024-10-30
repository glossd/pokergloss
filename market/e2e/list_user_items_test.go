package e2e

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"testing"
)

func TestListUserItems(t *testing.T) {
	t.Cleanup(cleanUp)
	buyGlassOfWine(t)
	body := listUserItems(t)
	assert.EqualValues(t, 2, gjson.Get(body, "#").Int())
}

func listUserItems(t *testing.T) string {
	rr := testRouter.Request(t, http.MethodGet, fmt.Sprintf("/users/%s/items", defaultIdentity.UserId), nil, nil)
	assert.EqualValues(t, http.StatusOK, rr.Code, rr.Body.String())
	body := rr.Body.String()
	return body
}
