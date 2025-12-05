package session

import (
	"sync"
	"time"
)

type SessionManager struct {
	sessions map[SessionID]*BaseSession
	mutex    *sync.RWMutex
}

func NewSessionManager() *SessionManager {
	r := &SessionManager{
		sessions: make(map[SessionID]*BaseSession),
		mutex:    &sync.RWMutex{},
	}

	// 启动一个后台 goroutine 定期清理过期会话
	go r.cleanupExpiredSessions()
	return r
}

func (m *SessionManager) GetSession(connID string) (*BaseSession, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	session, ok := m.sessions[SessionID(connID)]
	if !ok || session.IsExpired() {
		return nil, false
	}
	return session, ok
}

func (m *SessionManager) AddSession(session *BaseSession) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sessions[session.ID] = session
}

func (m *SessionManager) RemoveSession(connID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.sessions, SessionID(connID))
}

func (m *SessionManager) GetSessionsByType(sessionType SessionType) []*BaseSession {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	sessions := make([]*BaseSession, 0)
	for _, session := range m.sessions {
		if session.Type == sessionType && !session.IsExpired() {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

func (m *SessionManager) cleanupExpiredSessions() {
	// 创建一个10分钟的定时器，用于定期清理过期会话
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.mutex.Lock()
		for connID, session := range m.sessions {
			if session.IsExpired() {
				delete(m.sessions, connID)
			}
		}
		m.mutex.Unlock()
	}
}

func (m *SessionManager) GetSessionCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.sessions)
}
