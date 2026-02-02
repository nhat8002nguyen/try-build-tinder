package websocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan Message
	UserID uuid.UUID
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("WebSocket read error")
			}
			break
		}

		var incoming IncomingMessage
		if err := json.Unmarshal(messageBytes, &incoming); err != nil {
			log.Error().Err(err).Msg("Failed to parse WebSocket message")
			continue
		}

		c.handleMessage(incoming)
	}
}

func (c *Client) handleMessage(msg IncomingMessage) {
	switch msg.Type {
	case "ping":
		c.send <- Message{Type: "pong", Payload: nil}
	case "typing":
		var payload struct {
			MatchID string `json:"match_id"`
		}
		if err := json.Unmarshal(msg.Payload, &payload); err == nil {
			c.hub.Broadcast(Message{
				Type: "typing",
				Payload: map[string]interface{}{
					"user_id":  c.UserID.String(),
					"match_id": payload.MatchID,
				},
			})
		}
	default:
		log.Debug().Str("type", msg.Type).Msg("Unknown message type")
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal message")
				continue
			}

			w.Write(data)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("WebSocket upgrade failed")
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan Message, 256),
		UserID: userID,
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}
