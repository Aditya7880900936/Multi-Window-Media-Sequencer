package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()

			if _, exists := h.Clients[client]; exists {
				delete(h.Clients, client)
				close(client.Send)
			}

			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()

			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					// Client is too slow; skip this broadcast.
				}
			}

			h.mu.RUnlock()
		}
	}
}
