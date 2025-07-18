package session

import "sync"

type Session struct {
	HandlerID string
	State     interface{}
}

type Manager struct {
	sessions map[int]*Session // userID -> session
	mu       sync.Mutex
}

func New() *Manager {
	return &Manager{
		sessions: make(map[int]*Session),
	}
}

func (sm *Manager) Get(userID int) (*Session, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session, exists := sm.sessions[userID]
	return session, exists
}

func (sm *Manager) Set(userID int, handlerID string, state interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[userID] = &Session{
		HandlerID: handlerID,
		State:     state,
	}
}

func (sm *Manager) Delete(userID int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, userID)
}
