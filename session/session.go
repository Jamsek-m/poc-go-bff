package session

import (
	"fmt"
	"net/http"
	"poc-go-bff/session/store"
	"poc-go-bff/session/store/memory"
	"sync"
	"time"
)

var lock sync.Mutex

var sessionStorage store.Store

func GetStore() store.Store {
	if sessionStorage == nil {
		lock.Lock()
		defer lock.Unlock()
		if sessionStorage != nil {
			return sessionStorage
		}
		sessionStorage = memory.NewSessionStore()
	}
	return sessionStorage
}

func GetSessionFromCookie(cookieName string, req *http.Request) (*store.UserSession, error) {
	cookie, err := req.Cookie(cookieName)
	fmt.Println(cookie)
	if err == nil {
		sessionId := cookie.Value
		fmt.Println("Cookie session: ", sessionId)
		if existingSession := GetStore().Get(sessionId); existingSession != nil {
			fmt.Println("found session:", existingSession)
			return existingSession, nil
		} else {
			fmt.Println("not found session!")
		}
	}
	return nil, fmt.Errorf("error retrieving session from cookie")
}

func NewSessionCookie(sessionID string, res http.ResponseWriter) *http.Cookie {
	return &http.Cookie{
		Name:     "BFF_ID",
		Value:    sessionID,
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	}
}
