package database

import (
	"database/sql"
	"sync"
)

// ConnectionManager 简单的数据库连接管理器
type ConnectionManager struct {
	connections map[string]*sql.DB
	mutex       sync.RWMutex
}

var (
	globalManager *ConnectionManager
	once          sync.Once
)

// GetConnectionManager 获取全局连接管理器实例
func GetConnectionManager() *ConnectionManager {
	once.Do(func() {
		globalManager = &ConnectionManager{
			connections: make(map[string]*sql.DB),
		}
	})
	return globalManager
}

// AddConnection 添加数据库连接
func (cm *ConnectionManager) AddConnection(name string, db *sql.DB) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.connections[name] = db
}

// GetConnection 获取指定名称的数据库连接
func (cm *ConnectionManager) GetConnection(name string) (*sql.DB, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	conn, exists := cm.connections[name]
	return conn, exists
}

// GetDefaultConnection 获取默认数据库连接
func (cm *ConnectionManager) GetDefaultConnection() (*sql.DB, bool) {
	return cm.GetConnection("default")
}

// RemoveConnection 移除数据库连接
func (cm *ConnectionManager) RemoveConnection(name string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.connections, name)
}

// ListConnections 列出所有连接名称
func (cm *ConnectionManager) ListConnections() []string {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	names := make([]string, 0, len(cm.connections))
	for name := range cm.connections {
		names = append(names, name)
	}
	return names
}
