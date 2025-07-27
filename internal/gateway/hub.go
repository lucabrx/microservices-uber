package gateway

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	pb_driver "github.com/lukabrx/uber-clone/api/proto/driver/v1"
)

type Hub struct {
	clients      map[*websocket.Conn]bool
	mu           sync.Mutex
	driverClient pb_driver.DriverServiceClient
}

func NewHub(driverClient pb_driver.DriverServiceClient) *Hub {
	return &Hub{
		clients:      make(map[*websocket.Conn]bool),
		driverClient: driverClient,
	}
}

func (h *Hub) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()
	log.Printf("Client added. Total clients: %d", len(h.clients))
}

func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()
	log.Printf("Client removed. Total clients: %d", len(h.clients))
}

func (h *Hub) Broadcast(drivers []*pb_driver.Driver) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		if err := client.WriteJSON(drivers); err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			// On error, assume the client has disconnected and remove them.
			client.Close()
			delete(h.clients, client)
		}
	}
}
