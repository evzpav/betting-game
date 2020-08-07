package domain

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
)

type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

type Game struct {
	GameRunning  bool
	GameSnapshot []string
	RoundCounter int
	StopGame     chan bool
	Cron         *cron.Cron
}

type GameService interface {
	ServeWs(w http.ResponseWriter, r *http.Request)
	RunHub()
	Register(c *Client)
	Unregister(c *Client)
	Broadcast([]byte)
	RegisterNewClient(*websocket.Conn)
}
