package session

import (
	"context"
	"sync"
	"time"
)

type Session struct {
	CreatedAt *time.Time
	UpdatedAt *time.Time
	State     interface{}
}

type Manager struct {
	sessions map[int]*Session // userID -> session
	mu       sync.RWMutex
}

func New() *Manager {
	return &Manager{
		sessions: make(map[int]*Session),
	}
}

func (m *Manager) Get(ctx context.Context, userID int) (*Session, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	session, exists := m.sessions[userID]
	return session, exists
}

func (m *Manager) Set(ctx context.Context, userID int, s *Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[userID] = s
	return nil
}

func (m *Manager) Delete(ctx context.Context, userID int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, userID)

	return nil
}
