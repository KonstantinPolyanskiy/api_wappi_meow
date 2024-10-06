package clientPool

import (
	"go.mau.fi/whatsmeow"
	"sync"
)

// ClientPool предоставляет методы для управления множеством клиентов Whatsmeow
type ClientPool struct {
	clients map[string]*whatsmeow.Client
	mu      sync.RWMutex
}

// NewClientPool создает новый экземпляр ClientPool
func NewClientPool() *ClientPool {
	return &ClientPool{
		clients: make(map[string]*whatsmeow.Client),
	}
}

// AddClient добавляет клиента Whatsmeow в пул
func (cp *ClientPool) AddClient(id string, client *whatsmeow.Client) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.clients[id] = client
}

// GetClient получает клиента из пула по его ID
func (cp *ClientPool) GetClient(id string) (*whatsmeow.Client, bool) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	client, exists := cp.clients[id]
	return client, exists
}

// RemoveClient удаляет клиента из пула по его ID
func (cp *ClientPool) RemoveClient(id string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if client, exists := cp.clients[id]; exists {
		client.Disconnect()
		delete(cp.clients, id)
	}
}

// ListClients возвращает срез всех клиентов в пуле
func (cp *ClientPool) ListClients() []*whatsmeow.Client {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	clients := make([]*whatsmeow.Client, 0, len(cp.clients))
	for _, client := range cp.clients {
		clients = append(clients, client)
	}
	return clients
}

func (cp *ClientPool) UpdateClientID(oldID, newID string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if client, exists := cp.clients[oldID]; exists {
		delete(cp.clients, oldID)
		cp.clients[newID] = client
	}
}
