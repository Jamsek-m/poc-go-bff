package session

import (
	"github.com/gorilla/sessions"
	"github.com/michaeljs1990/sqlitestore"
	"poc-go-bff/config"
)

var sessionStore sessions.Store

func InitSessions() {
	if config.GetConfig().Sessions.Store == "sqlite3" {
		sqliteStore, err := sqlitestore.NewSqliteStore("./db.sqlite", "bff_sessions", "/", 3600, []byte(config.GetConfig().Sessions.Secret))
		if err == nil {
			sessionStore = sqliteStore
		} else {
			panic(err)
		}
	} else if config.GetConfig().Sessions.Store == "memory" {
		sessionStore = sessions.NewCookieStore([]byte(config.GetConfig().Sessions.Secret))
	} else {

	}
}

func GetStore() sessions.Store {
	return sessionStore
}
