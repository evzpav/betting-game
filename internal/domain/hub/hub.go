package hub

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"gitlab.com/evzpav/betting-game/internal/domain"
	"gitlab.com/evzpav/betting-game/pkg/log"
)

type service struct {
	hub *domain.Hub
	log log.Logger
}

func NewService(log log.Logger) *service {
	return &service{
		hub: newHub(),
		log: log,
	}
}

func newHub() *domain.Hub {
	return &domain.Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *domain.Client),
		Unregister: make(chan *domain.Client),
		Clients:    make(map[*domain.Client]bool),
	}
}

func (s *service) RunHub() {
	for {
		select {
		case client := <-s.hub.Register:
			s.hub.Clients[client] = true
		case client := <-s.hub.Unregister:
			if _, ok := s.hub.Clients[client]; ok {
				delete(s.hub.Clients, client)
				close(client.Send)
			}
		case message := <-s.hub.Broadcast:
			for client := range s.hub.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(s.hub.Clients, client)
				}
			}
		}
	}
}

func (s *service) AddNewWebsocketClient(conn *websocket.Conn) {
	cli := domain.NewClient(conn, s.log)
	s.Register(cli)

	go cli.WritePump()
	go cli.ReadPump(s.hub)
}

func (s *service) Register(c *domain.Client) {
	s.hub.Register <- c
}

func (s *service) Unregister(c *domain.Client) {
	s.hub.Unregister <- c
}

func (s *service) Broadcast(messageType domain.MessageType, data interface{}) error {
	var wsMessage = domain.Message{
		MessageType: messageType,
		Data:        data,
	}

	bs, err := json.Marshal(wsMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	s.hub.Broadcast <- bs
	return nil
}
