package websocket

import (
	"lock-stock-v2/external/websocket"
	"log"
	"sync"
)

type Manager struct {
	Clients        map[*websocket.Client]bool
	BroadcastChan  chan []byte
	RegisterChan   chan *websocket.Client
	UnregisterChan chan *websocket.Client
	mu             sync.Mutex
}

func NewWebSocketManager() *Manager {
	return &Manager{
		Clients:        make(map[*websocket.Client]bool),
		BroadcastChan:  make(chan []byte),
		RegisterChan:   make(chan *websocket.Client),
		UnregisterChan: make(chan *websocket.Client),
	}
}

func (manager *Manager) Run() {
	for {
		select {
		case client := <-manager.RegisterChan:
			log.Printf("Client received from register channel: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
			manager.mu.Lock()
			manager.Clients[client] = true
			manager.mu.Unlock()
			log.Printf("Client successfully registered: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)

		case client := <-manager.UnregisterChan:
			log.Printf("Client received from unregister channel: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
			manager.mu.Lock()
			if _, ok := manager.Clients[client]; ok {
				delete(manager.Clients, client)
				close(client.Send)
				log.Printf("Client disconnected: %+v, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
			}
			manager.mu.Unlock()

		case message := <-manager.BroadcastChan:
			manager.mu.Lock()
			log.Printf("Broadcasting message to all clients: %s\n", string(message))
			for client := range manager.Clients {
				select {
				case client.Send <- message:
				default:
					log.Printf("Failed to send message to client: %+v\n", client.Conn.RemoteAddr())
					close(client.Send)
					delete(manager.Clients, client)
				}
			}
			manager.mu.Unlock()
		}
	}
}

func (manager *Manager) PublishToRoom(roomID string, message []byte) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	log.Printf("Current clients: %+v\n", manager.Clients)
	log.Printf("Current clients before publishing: %+v\n", manager.Clients)

	log.Printf("Publishing message to room: %s, Message: %s\n", roomID, string(message))

	for client := range manager.Clients {
		log.Printf("Checking client: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
		if client.RoomID == roomID {
			log.Printf("Client matched for RoomID: %s\n", roomID)
			select {
			case client.Send <- message:
				log.Printf("Message sent to client: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
			default:
				log.Printf("Failed to send message to client: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
				close(client.Send)
				delete(manager.Clients, client)
			}
		}
	}
}

func (manager *Manager) Register(client *websocket.Client) {
	manager.RegisterChan <- client
}

func (manager *Manager) Unregister(client *websocket.Client) {
	manager.UnregisterChan <- client
}

func (manager *Manager) Broadcast(message []byte) {
	manager.BroadcastChan <- message
}
