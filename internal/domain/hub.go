package domain

import "github.com/gorilla/websocket"

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

type MessageType string

const (
	RoundType          MessageType = "round"
	StartType          MessageType = "start"
	EndType            MessageType = "end"
	OverallRankingType MessageType = "overallranking"
	IntervalTickerType MessageType = "intervalTicker"
)

type Message struct {
	MessageType MessageType `json:"type"`
	Data        interface{} `json:"data"`
}

type HubService interface {
	RunHub()
	AddNewWebsocketClient(conn *websocket.Conn)
	Broadcast(messageType MessageType, data interface{}) error
}
