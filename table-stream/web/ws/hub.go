package ws

import (
	"github.com/glossd/pokergloss/gomq/mqws"
	"github.com/glossd/pokergloss/table-stream/web/model"
	"github.com/glossd/pokergloss/table-stream/web/mq/mqpub"
)

type Hub struct {
	tableID string
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	direct chan *mqws.TableUserEvents

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func (h *Hub) initHub(tableID string) {
	h.tableID = tableID
	h.broadcast = make(chan []byte)
	h.direct = make(chan *mqws.TableUserEvents)
	h.register = make(chan *Client)
	h.unregister = make(chan *Client)
	h.clients = make(map[*Client]bool)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			mqpub.Publish(&mqws.TableMessage{
				ToEntityIds: []string{h.tableID},
				Events:      []*mqws.Event{model.BuildNewConnection(client.iden)},
			})
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				mqpub.Publish(&mqws.TableMessage{
					ToEntityIds: []string{h.tableID},
					Events:      []*mqws.Event{model.BuildLeftConnection(client.iden)},
				})
			}
		case tableEvents := <-h.direct:
			for client := range h.clients {
				allEvents := make([]*mqws.Event, 0,
					len(tableEvents.BeforeEvents)+len(tableEvents.AfterEvents)+len(tableEvents.NotFoundUsersEvents))
				allEvents = append(allEvents, tableEvents.BeforeEvents...)
				if client.iden != nil {
					if eventsWrapper, ok := tableEvents.UserEvents[client.iden.UserId]; ok {
						allEvents = append(allEvents, eventsWrapper.Events...)
					} else {
						allEvents = append(allEvents, tableEvents.NotFoundUsersEvents...)
					}
				} else {
					allEvents = append(allEvents, tableEvents.NotFoundUsersEvents...)
				}

				allEvents = append(allEvents, tableEvents.AfterEvents...)

				select {
				case client.send <- mqws.EventsToJson(allEvents):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
