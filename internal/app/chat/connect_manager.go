package domain

import (
	"net"
	"sync"
)

type ConnectionManager struct {
	connections map[uint32]*net.Conn
	mutex       sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]*net.Conn),
	}
}

func (cm *ConnectionManager) AddConnection(userID uint32, conn *net.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connections[userID] = conn
}

func (cm *ConnectionManager) GetConnection(userID uint32) (*net.Conn, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	conn, ok := cm.connections[userID]
	return conn, ok
}

func (cm *ConnectionManager) RemoveConnection(userID uint32) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.connections, userID)
}
