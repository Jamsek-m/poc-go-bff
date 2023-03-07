package session

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"net/http"
)

type (
	Store interface {
		GetSessionValue(req *http.Request, paramName string) interface{}
		SetSessionValue(req *http.Request, paramName string, value string)
		CommitSessionChanges(req *http.Request)
		CreateSession(req *http.Request)
		DestroySession(req *http.Request) error
		Instance() *scs.SessionManager
	}

	storeImpl struct {
		sessionManager scs.SessionManager
	}
)

var store Store

func (s *storeImpl) GetSessionValue(req *http.Request, paramName string) interface{} {
	ctx := req.Context()
	if s.sessionManager.Exists(ctx, paramName) {
		return s.sessionManager.Get(ctx, paramName)
	}
	return nil
}

func (s *storeImpl) SetSessionValue(req *http.Request, paramName string, value string) {
	s.sessionManager.Put(req.Context(), paramName, value)
}

func (s *storeImpl) CommitSessionChanges(req *http.Request) {
	_ = s.sessionManager.RenewToken(req.Context())
}

func (s *storeImpl) CreateSession(req *http.Request) {

}

func (s *storeImpl) DestroySession(req *http.Request) error {
	if err := s.sessionManager.Destroy(req.Context()); err != nil {
		return fmt.Errorf("error destroying session: %w", err)
	}
	return nil
}

func (s *storeImpl) Instance() *scs.SessionManager {
	return &s.sessionManager
}

func InitStore() {
	sessionManager := initSessions()
	store = &storeImpl{
		sessionManager: *sessionManager,
	}
}

func Current() Store {
	return store
}
