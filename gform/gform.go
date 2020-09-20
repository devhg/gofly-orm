package gform

import (
	"database/sql"
	"github.com/QXQZX/gofly-orm/gform/log"
	"github.com/QXQZX/gofly-orm/gform/session"
)

type Engine struct {
	db *sql.DB
}

// create a new engine
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}

	e = &Engine{db: db}
	log.Info(driver, "Database connect success")
	return
}

// close the engine
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close the Database")
	}
	log.Info("Closed Database success")
}

// create a session
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}
