package conn

import "sync"

type manager struct {
	// 全局连接注册表，管理所有活跃链接
	conns map[string]Connection
	// 用户连接注册表，管理每个用户的多端链接，内层的key与conn一致
	userConns map[uint64]map[string]Connection
	mu        sync.RWMutex
}

func NewManager() ConnectionManager {
	return &manager{
		conns:     make(map[string]Connection),
		userConns: make(map[uint64]map[string]Connection),
	}
}

func (m *manager) Register(conn Connection) error {
	connID := conn.ID()
	userID := conn.UserID()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.conns[connID] = conn
	if _, exists := m.userConns[userID]; !exists {
		m.userConns[userID] = make(map[string]Connection)
	}
	m.userConns[userID][connID] = conn

	return nil
}

func (m *manager) UnRegister(conn Connection) error {
	connID := conn.ID()
	userID := conn.UserID()

	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.conns[connID]
	if !ok {
		return ErrConnectionNotFound
	}

	delete(m.conns, connID)
	if userMap, exists := m.userConns[userID]; exists {
		delete(userMap, connID)
		if len(userMap) == 0 {
			delete(m.userConns, userID)
		}
	}

	return nil
}

func (m *manager) GetConnection(connID string) (Connection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, ok := m.conns[connID]
	if !ok {
		return nil, ErrConnectionNotFound
	}

	return conn, nil
}

func (m *manager) GetConnectionsByUserID(userID uint64) []Connection {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var conns []Connection
	for _, conn := range m.userConns[userID] {
		conns = append(conns, conn)
	}

	return conns
}
