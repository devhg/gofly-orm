package gform

import (
	"database/sql"

	"github.com/devhg/gofly-orm/dialect"
	"github.com/devhg/gofly-orm/log"
	"github.com/devhg/gofly-orm/session"
)

type Engine struct {
	// 官方提供的 可以包含一个或者多个连接的 并发安全的 数据库连接池，
	db     *sql.DB
	dbName string

	// 用于获取 golang类型 和 数据库类型 的 对应关系
	dialect dialect.Dialect
}

// 创建一个orm映射引擎
// 通过此引擎可以创建 session 会话
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
	// 获取当前数据的 dialect
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

// create a session  传入 sql.DB 对象和 dialect 创建一个会话 session.Session
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) (interface{}, error)

// 用于支持事务的sql操作
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
