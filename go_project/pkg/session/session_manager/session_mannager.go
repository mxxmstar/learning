package session_manager

import (
	"sync"
	"time"

	"github.com/mxxmstar/learning/pkg/session"
)

type SessionManager struct {
	sessions map[session.SessionId]*session.BaseSession
	mutex    *sync.RWMutex
}

func NewSessionManager() *SessionManager {
	r := &SessionManager{
		sessions: make(map[session.SessionId]*session.BaseSession),
		mutex:    &sync.RWMutex{},
	}

	// 启动一个后台 goroutine 定期清理过期会话
	go r.cleanupExpiredSessions()
	return r
}

func (m *SessionManager) GetSession(connId string) (*session.BaseSession, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	session, ok := m.sessions[session.SessionId(connId)]
	if !ok || session.IsExpired() {
		return nil, false
	}
	return session, ok
}

func (m *SessionManager) AddSession(session *session.BaseSession) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sessions[session.Id] = session
}

func (m *SessionManager) RemoveSession(connId string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.sessions, session.SessionId(connId))
}

func (m *SessionManager) GetSessionsByType(sessionType session.SessionType) []*session.BaseSession {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	sessions := make([]*session.BaseSession, 0)
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
		for connId, session := range m.sessions {
			if session.IsExpired() {
				delete(m.sessions, connId)
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
