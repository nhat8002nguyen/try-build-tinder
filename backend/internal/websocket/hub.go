package websocket

import (
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	clients    map[*Client]bool
	userConns  map[uuid.UUID]map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		userConns:  make(map[uuid.UUID]map[*Client]bool),
		broadcast:  make(chan Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if h.userConns[client.UserID] == nil {
				h.userConns[client.UserID] = make(map[*Client]bool)
			}
			h.userConns[client.UserID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if h.userConns[client.UserID] != nil {
					delete(h.userConns[client.UserID], client)
					if len(h.userConns[client.UserID]) == 0 {
						delete(h.userConns, client.UserID)
					}
				}
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.userConns[userID]; ok {
		for client := range clients {
			select {
			case client.send <- message:
			default:
			}
		}
	}
}

func (h *Hub) Broadcast(message Message) {
	h.broadcast <- message
}

func (h *Hub) IsUserOnline(userID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.userConns[userID]
	return ok
}

func (h *Hub) GetOnlineUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.userConns)
}
