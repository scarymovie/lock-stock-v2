package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte
	RoomID string
}

type WebSocketManager struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (manager *WebSocketManager) Run() {
	for {
		select {
		case client := <-manager.Register:
			log.Printf("Client received from register channel: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
			manager.mu.Lock()
			manager.Clients[client] = true
			manager.mu.Unlock()
			log.Printf("Client successfully registered: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)

		case client := <-manager.Unregister:
			log.Printf("Client received from unregister channel: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
			manager.mu.Lock()
			if _, ok := manager.Clients[client]; ok {
				delete(manager.Clients, client)
				close(client.Send)
				log.Printf("Client disconnected: %+v, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
			}
			manager.mu.Unlock()

		case message := <-manager.Broadcast:
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

func (manager *WebSocketManager) PublishToRoom(roomID string, message []byte) {
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
