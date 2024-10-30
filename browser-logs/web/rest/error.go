package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/glossd/pokergloss/auth"
	"io/ioutil"
	"time"
)

type Example struct {
	ID string
}

// @ID postError
// @Accept  json
// @Param input body string true "input json data"
// @Success 204 "Posted error"
// @Failure 400 {object} ErrorRes
// @Failure 500 {object} ErrorRes
// @Router /errors [post]
func PostError(c *gin.Context) {
	var jsonData map[string]interface{}
	data, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(data, &jsonData); err != nil {
		jsonData["msg"] = data
	}

	iden, ok := auth.IdSafe(c)
	if ok {
		jsonData["userId"] = iden.UserId
		jsonData["username"] = iden.Username
	} else {
		delete(jsonData, "userId")
		delete(jsonData, "username")
	}

	jsonData["severity"] = "WARN"
	jsonData["time"] = time.Now().Format(time.RFC3339Nano)

	//logJson, err := json.Marshal(jsonData)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, E(err))
	//	return
	//}
	//
	//log.Warn(string(logJson))

	c.Writer.WriteHeader(204)
}
