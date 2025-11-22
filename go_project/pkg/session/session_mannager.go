package session

import "sync"

type SessionManager struct {
	sessions sync.Map
}

func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

func (m *SessionManager) GetSession(connID string) (*Session, bool) {
	session, ok := m.sessions.Load(connID)
	if !ok {
		return nil, false
	}
	return session.(*Session), true
}

func (m *SessionManager) AddSession(session *Session) {
	m.sessions.Store(session.ConnID, session)
}

func (m *SessionManager) RemoveSession(connID string) {
	m.sessions.Delete(connID)
}

func (sm *SessionManager) Range(f func(connID string, s *Session) bool) {
	sm.sessions.Range(func(key, value any) bool {
		return f(key.(string), value.(*Session))
	})
}
