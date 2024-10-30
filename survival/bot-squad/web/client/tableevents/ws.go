package tableevents

import (
	"fmt"
	"github.com/glossd/pokergloss/survival/bot-squad/conf"
	"github.com/glossd/pokergloss/survival/bot-squad/domain"
	"github.com/glossd/pokergloss/survival/bot-squad/service"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strings"
)

func StreamEventsWS(config conf.Config) {
	url := fmt.Sprintf("wss://%s/api/table-stream/tables/%s?token=%s", config.TableEvents.Host, config.TableID, config.Token)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("ws dial failed:", err)
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Errorf("ws read message: %s", err)
			return
		}
		for _, eventsJSON := range strings.Split(string(message), "\n") {
			result := gjson.Get(eventsJSON, "@this")
			var events []*domain.Event
			result.ForEach(func(key, value gjson.Result) bool {
				events = append(events, domain.NewEventBytes([]byte(value.String())))
				return true
			})
			service.HandleEvents(config, events)
		}
	}
}
