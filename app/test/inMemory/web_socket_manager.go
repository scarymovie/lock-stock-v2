package inMemory

import (
	externalWs "lock-stock-v2/internal/websocket"
	"sync"
)

var _ externalWs.Manager = (*InMemoryWebSocketManager)(nil)

type InMemoryWebSocketManager struct {
	mu             sync.Mutex
	clients        map[string][]*externalWs.Client
	broadcasted    [][]byte
	messagesByRoom map[string][][]byte
}

func NewInMemoryWebSocketManager() *InMemoryWebSocketManager {
	return &InMemoryWebSocketManager{
		clients:        make(map[string][]*externalWs.Client),
		messagesByRoom: make(map[string][][]byte),
		broadcasted:    make([][]byte, 0),
	}
}

func (m *InMemoryWebSocketManager) Run() {
}

func (m *InMemoryWebSocketManager) PublishToRoom(roomID string, message []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messagesByRoom[roomID] = append(m.messagesByRoom[roomID], message)
}

func (m *InMemoryWebSocketManager) Register(client *externalWs.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	roomID := client.RoomID
	m.clients[roomID] = append(m.clients[roomID], client)
}

func (m *InMemoryWebSocketManager) Unregister(client *externalWs.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	roomID := client.RoomID
	clients := m.clients[roomID]

	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	m.clients[roomID] = clients
}

func (m *InMemoryWebSocketManager) Broadcast(message []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.broadcasted = append(m.broadcasted, message)
}

// -----
// Дополнительные методы (не из интерфейса), чтобы в тестах было удобнее проверять результат
// -----

// GetMessagesForRoom возвращает список сообщений, отправленных в конкретную комнату.
func (m *InMemoryWebSocketManager) GetMessagesForRoom(roomID string) [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.messagesByRoom[roomID]
}

// GetBroadcastedMessages возвращает все сообщения, отосланные через Broadcast.
func (m *InMemoryWebSocketManager) GetBroadcastedMessages() [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.broadcasted
}

// GetClientsForRoom возвращает список клиентов, зарегистрированных в комнате.
func (m *InMemoryWebSocketManager) GetClientsForRoom(roomID string) []*externalWs.Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clients[roomID]
}
