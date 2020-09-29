package gform

import (
	"database/sql"
	"github.com/QXQZX/gofly-orm/gform/dialect"
	"github.com/QXQZX/gofly-orm/gform/log"
	"github.com/QXQZX/gofly-orm/gform/session"
)

type Engine struct {
	db      *sql.DB
	dbName  string
	dialect dialect.Dialect
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

	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	e = &Engine{
		db:      db,
		dbName:  driver,
		dialect: dial,
	}
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
	return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) (interface{}, error)

func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = s.RollBack()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.RollBack() // err is non-nil; don't change it
		} else {
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}
