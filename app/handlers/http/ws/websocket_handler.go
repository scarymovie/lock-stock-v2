package ws

import (
	gorillaWs "github.com/gorilla/websocket"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"lock-stock-v2/internal/websocket"
	"log"
	"net/http"
)

type WebSocketHandler struct {
	Manager websocket.Manager
}

func NewWebSocketHandler(manager websocket.Manager) *WebSocketHandler {
	return &WebSocketHandler{Manager: manager}
}

func (h *WebSocketHandler) ConnectWebSocket(w http.ResponseWriter, r *http.Request, roomId openapi_types.UUID) {
	upgrader := gorillaWs.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	roomID := roomId.String()
	client := &websocket.Client{
		Conn:   conn,
		Send:   make(chan []byte),
		RoomID: roomID,
	}

	log.Printf("Registering client: %s, RoomID: %s\n", conn.RemoteAddr(), roomID)
	h.Manager.Register(client)
	log.Printf("Client sent to register channel: %s, RoomID: %s\n", conn.RemoteAddr(), roomID)

	go func() {
		log.Printf("Starting handleMessages for client: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
		h.handleMessages(client)
	}()
	go func() {
		log.Printf("Starting writeMessages for client: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
		h.writeMessages(client)
	}()
}

func (h *WebSocketHandler) writeMessages(client *websocket.Client) {
	defer func() {
		log.Printf("Closing connection for client: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
		client.Conn.Close()
	}()

	for message := range client.Send {
		log.Printf("Sending message to client: %s, Message: %s\n", client.Conn.RemoteAddr(), string(message))
		err := client.Conn.WriteMessage(gorillaWs.TextMessage, message)
		if err != nil {
			log.Printf("Error sending message to client: %s, Error: %v\n", client.Conn.RemoteAddr(), err)
			break
		}
	}
}

func (h *WebSocketHandler) handleMessages(client *websocket.Client) {
	defer func() {
		log.Printf("Unregistering client: %s, RoomID: %s\n", client.Conn.RemoteAddr(), client.RoomID)
		h.Manager.Unregister(client)
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from client: %s, Error: %v\n", client.Conn.RemoteAddr(), err)
			break
		}
		log.Printf("Message received from client: %s, Message: %s\n", client.Conn.RemoteAddr(), string(message))
		h.Manager.Broadcast(message)
	}
}
