package session

import (
	"database/sql"
	"github.com/cddgo/gofly-orm/clause"
	"github.com/cddgo/gofly-orm/dialect"
	"github.com/cddgo/gofly-orm/log"
	"github.com/cddgo/gofly-orm/schema"
	"strings"

	// 导入时会注册 sqlite3 的驱动
	_ "github.com/mattn/go-sqlite3"
)

//Session 负责与数据库的交互，那交互前的准备工作（比如连接/测试数据库）
//封装有两个目的，一是统一打印日志（包括 执行的SQL 语句和错误日志）。
//二是执行完成后，清空操作。这样 Session 可以复用，开启一次会话，
//可以执行多次 SQL。
type Session struct {
	db  *sql.DB
	tx  *sql.Tx
	sql strings.Builder
	// sql 中占位符对应的值
	sqlVars []interface{}

	// 不同数据类型集，支持不同的数据库
	dialect dialect.Dialect

	// 结构体和数据表的映射
	refTable *schema.Schema

	// sql构造器
	clause clause.Clause
}

// CommonDB is a minimal function set of db
// sql.DB sql.Tx 都包含下面的函数，进而他们都实现了 下面的 接口
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// 验证接口的有效性
var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

// 创建一个 session 对象
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

// DB returns tx if a tx begins. otherwise return *sql.DB
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// 存储sql   和  参数
func (s *Session) Raw(sql string, vars ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, vars...)
	return s
}

// 清理sql串 和 参数
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// 执行sql
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// package the QueryRow() method
// QueryRow gets a record from db
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// package the Query() method
// QueryRows gets a list of records from db
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// for test
func NewSession_() *Session {
	db, err := sql.Open("sqlite3", "gofly.db")
	if err != nil {
		panic(err)
	}
	getDialect, _ := dialect.GetDialect("sqlite3")
	return New(db, getDialect)
}
