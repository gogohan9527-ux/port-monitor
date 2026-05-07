package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]struct{}
	last    []byte
}

func NewHub() *Hub {
	return &Hub{clients: make(map[*Client]struct{})}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	last := h.last
	h.mu.Unlock()
	if last != nil {
		select {
		case c.send <- last:
		default:
		}
	}
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.send)
	}
	h.mu.Unlock()
}

func (h *Hub) Broadcast(payload []byte) {
	h.mu.Lock()
	h.last = payload
	for c := range h.clients {
		select {
		case c.send <- payload:
		default:
			// slow client; drop frame
		}
	}
	h.mu.Unlock()
}

func (h *Hub) Snapshot() []byte {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.last
}
