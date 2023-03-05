package memory

import (
	"github.com/google/uuid"
	"poc-go-bff/session/store"

	"sync"
)

type inMemorySessionStore struct {
	store map[string]*store.UserSession
	lock  sync.Mutex
}

func NewSessionStore() store.Store {
	return &inMemorySessionStore{
		store: make(map[string]*store.UserSession),
	}
}

func (s *inMemorySessionStore) getFromMap(key string) *store.UserSession {
	s.lock.Lock()
	defer s.lock.Unlock()
	return (s.store)[key]
}

func (s *inMemorySessionStore) saveToMap(key string, value *store.UserSession) {
	s.lock.Lock()
	defer s.lock.Unlock()
	(s.store)[key] = value
}

func (s *inMemorySessionStore) removeFromMap(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, exists := (s.store)[key]; exists {
		delete(s.store, key)
	}
}

func (s *inMemorySessionStore) Get(sessionID string) *store.UserSession {
	return s.getFromMap(sessionID)
}

func (s *inMemorySessionStore) Start() (string, *store.UserSession) {
	sessionId := uuid.New().String()
	userSession := &store.UserSession{
		ID: &sessionId,
	}
	s.saveToMap(sessionId, userSession)
	return sessionId, userSession
}

func (s *inMemorySessionStore) Destroy(sessionID string) {
	s.removeFromMap(sessionID)
}

func (s *inMemorySessionStore) Update(sessionId string, session *store.UserSession) *store.UserSession {
	existingSession := s.Get(sessionId)
	if existingSession == nil {
		id, newSession := s.Start()
		existingSession = newSession
		sessionId = id
	}
	s.saveToMap(*existingSession.ID, session)
	return s.Get(*existingSession.ID)
}
