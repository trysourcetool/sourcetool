package session

import (
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
)

const (
	// TTL for disconnected sessions
	disconnectedSessionTTL = 2 * time.Minute
	// Maximum number of disconnected sessions to keep
	maxDisconnectedSessions = 128
)

type disconnectedSession struct {
	session        *Session
	disconnectedAt time.Time
}

type Session struct {
	ID     uuid.UUID
	PageID uuid.UUID
	State  *State
}

func New(id uuid.UUID, pageID uuid.UUID) *Session {
	return &Session{
		ID:     id,
		PageID: pageID,
		State:  newState(),
	}
}

type SessionManager struct {
	activeSessions       map[uuid.UUID]*Session
	disconnectedSessions map[uuid.UUID]*disconnectedSession
	mu                   sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		activeSessions:       make(map[uuid.UUID]*Session),
		disconnectedSessions: make(map[uuid.UUID]*disconnectedSession),
	}
}

func (s *SessionManager) removeOldestDisconnectedSession() {
	if len(s.disconnectedSessions) == 0 {
		return
	}

	var oldestID uuid.UUID
	var oldestTime time.Time

	for id, ds := range s.disconnectedSessions {
		if oldestTime.IsZero() || ds.disconnectedAt.Before(oldestTime) {
			oldestID = id
			oldestTime = ds.disconnectedAt
		}
	}

	delete(s.disconnectedSessions, oldestID)
}

func (s *SessionManager) GetSession(id uuid.UUID) *Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.activeSessions[id]
}

func (s *SessionManager) SetSession(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ds, ok := s.disconnectedSessions[session.ID]; ok {
		session.State = ds.session.State
		delete(s.disconnectedSessions, session.ID)
	}

	s.activeSessions[session.ID] = session
}

func (s *SessionManager) DisconnectSession(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if session, ok := s.activeSessions[id]; ok {
		if len(s.disconnectedSessions) >= maxDisconnectedSessions {
			s.removeOldestDisconnectedSession()
		}

		s.disconnectedSessions[id] = &disconnectedSession{
			session:        session,
			disconnectedAt: time.Now(),
		}
		delete(s.activeSessions, id)

		go func(sessionID uuid.UUID) {
			time.Sleep(disconnectedSessionTTL)
			s.mu.Lock()
			delete(s.disconnectedSessions, sessionID)
			s.mu.Unlock()
		}(id)
	}
}
