package websocket

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// Match rooms: matchID -> set of clients
	rooms map[uuid.UUID]map[*Client]bool

	// Mutex for rooms
	roomsMu sync.RWMutex
}

// NewHub creates a new hub
func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[uuid.UUID]map[*Client]bool),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// Remove from any room
				h.roomsMu.Lock()
				for matchID, room := range h.rooms {
					if _, ok := room[client]; ok {
						delete(room, client)
						if len(room) == 0 {
							delete(h.rooms, matchID)
						}
					}
				}
				h.roomsMu.Unlock()
			}
		case message := <-h.Broadcast:
			// Broadcast to all clients
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

// JoinRoom adds a client to a match room
func (h *Hub) JoinRoom(client *Client, matchID uuid.UUID) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	if _, ok := h.rooms[matchID]; !ok {
		h.rooms[matchID] = make(map[*Client]bool)
	}
	h.rooms[matchID][client] = true
	client.SetMatchID(&matchID)
}

// LeaveRoom removes a client from a match room
func (h *Hub) LeaveRoom(client *Client, matchID uuid.UUID) {
	h.roomsMu.Lock()
	defer h.roomsMu.Unlock()

	if room, ok := h.rooms[matchID]; ok {
		delete(room, client)
		if len(room) == 0 {
			delete(h.rooms, matchID)
		}
	}
	client.SetMatchID(nil)
}

// BroadcastToRoom sends a message to all clients in a room
func (h *Hub) BroadcastToRoom(matchID uuid.UUID, message []byte) {
	h.roomsMu.RLock()
	defer h.roomsMu.RUnlock()

	if room, ok := h.rooms[matchID]; ok {
		for client := range room {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(room, client)
			}
		}
	}
}

// GetRoomClients returns clients in a room
func (h *Hub) GetRoomClients(matchID uuid.UUID) []*Client {
	h.roomsMu.RLock()
	defer h.roomsMu.RUnlock()

	var clients []*Client
	if room, ok := h.rooms[matchID]; ok {
		for client := range room {
			clients = append(clients, client)
		}
	}
	return clients
}

// LogStats logs hub statistics (for debugging)
func (h *Hub) LogStats() {
	h.roomsMu.RLock()
	defer h.roomsMu.RUnlock()

	log.Printf("Hub stats: %d clients, %d rooms", len(h.clients), len(h.rooms))
	for matchID, room := range h.rooms {
		log.Printf("  Room %s: %d clients", matchID, len(room))
	}
}