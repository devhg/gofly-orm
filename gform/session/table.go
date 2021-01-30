package session

import (
	"fmt"
	"github.com/cddgo/gofly-orm/log"
	"github.com/cddgo/gofly-orm/schema"
	"reflect"
	"strings"
)

//Model() 方法用于给 refTable 赋值。解析操作是比较耗时的，
//因此将解析的结果保存在成员变量 refTable 中，即使 Model()
//被调用多次，如果传入的结构体名称不发生变化，则不会更新
//refTable 的值。
func (s *Session) Model(value interface{}) *Session {
	// nil or different model, update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// 创建
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s",
			field.Name, field.Type, field.Tag))
	}

	join := strings.Join(columns, ", ")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, join)).Exec()
	return err
}

//删除
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s;", s.RefTable().Name)).Exec()
	return err
}

//存在性判断
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSql(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
