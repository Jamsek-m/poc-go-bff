package session

import (
	"database/sql"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"log"
	"poc-go-bff/config"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func initSessions() *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Name = config.GetConfig().Sessions.CookieName
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Path = "/"

	if config.GetConfig().Sessions.Store == "sqlite3" {
		sessionManager.Store = configureSqlite3Store()
	} else if config.GetConfig().Sessions.Store == "redis" {
		panic("not implemented")
	} else {
		sessionManager.Store = configureMemoryStore()
	}
	return sessionManager
}

func configureMemoryStore() scs.Store {
	return memstore.New()
}

func configureSqlite3Store() scs.Store {
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sessions (" +
		"token TEXT PRIMARY KEY " +
		", data BLOB NOT NULL " +
		", expiry REAL NOT NULL" +
		");")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry);")
	if err != nil {
		panic(err)
	}

	return sqlite3store.New(db)
}
