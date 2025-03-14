package websocket

import "github.com/gorilla/websocket"

type Manager interface {
	Run()
	PublishToRoom(roomID string, message []byte)
	Register(client *Client)
	Unregister(client *Client)
	Broadcast(message []byte)
}

type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte
	RoomID string
}
