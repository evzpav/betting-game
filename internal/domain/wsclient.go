package domain

import (
	"bytes"

	"gitlab.com/evzpav/betting-game/pkg/log"

	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type MessageType string

const (
	RoundType          MessageType = "round"
	StartType          MessageType = "start"
	RestartType        MessageType = "restart"
	EndType            MessageType = "end"
	OverallRankingType MessageType = "overallranking"
)

type Message struct {
	MessageType MessageType `json:"type"`
	Data        interface{} `json:"data"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	Send chan []byte
	log  log.Logger
}

func NewClient(hub *Hub, conn *websocket.Conn, log log.Logger) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		Send: make(chan []byte, 256),
		log:  log,
	}
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		c.hub.Unregister <- c
		if err := c.conn.Close(); err != nil {
			c.log.Error().Sendf("failed to close connection: %v", err)
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.log.Error().Sendf("failed set read deadline: %v", err)
	}

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.log.Error().Sendf("unexpected close error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		c.hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			c.log.Error().Sendf("failed to close connection: %v", err)
		}
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.log.Error().Err(err).Sendf("%v", err)
			}

			if !ok {
				// The hub closed the channel.
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.log.Error().Err(err).Sendf("%v", err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, err = w.Write(message)
			if err != nil {
				c.log.Error().Sendf("failed to write message: %v", err)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.log.Error().Err(err).Sendf("%v", err)
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
