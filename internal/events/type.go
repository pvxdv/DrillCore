package events

import "sync"

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
	Callback
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}

type Session struct {
	HandlerID string
	State     interface{}
}

type SessionManager struct {
	sessions map[int]*Session // userID -> session
	mu       sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[int]*Session),
	}
}

func (sm *SessionManager) Get(userID int) (*Session, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	session, exists := sm.sessions[userID]
	return session, exists
}

func (sm *SessionManager) Set(userID int, handlerID string, state interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.sessions[userID] = &Session{
		HandlerID: handlerID,
		State:     state,
	}
}

func (sm *SessionManager) Delete(userID int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.sessions, userID)
}
