package session

import (
	"errors"
	"reflect"

	"github.com/devhg/gofly-orm/clause"
)

func (s *Session) runSQL(sql string, vars []interface{}) (int64, error) {
	exec, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)

	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)

	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

// 1) destSlice.Type().Elem() 获取切片的单个元素的类型 destType，
// 使用 reflect.New() 方法创建一个 destType 的实例，作为 Model() 的入参，
// 映射出表结构 RefTable()。
// 2）根据表结构，使用 clause 构造出 SELECT 语句，查询到所有符合条件的记录rows。
// 3）遍历每一行记录，利用反射创建 destType 的实例 dest，将 dest 的所有字段平铺开，
// 构造切片values。【重点】
// 4）调用 rows.Scan() 将该行记录每一列的值依次赋值给 values 中的每一个字段。【重点】
// 5）将 dest 添加到切片 destSlice 中。循环直到所有的记录都添加到切片 destSlice 中
func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()

	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	// fmt.Println(destSlice)                                // []
	// fmt.Println(destSlice.Type())                         // []session.User
	// fmt.Println(destSlice.Type().Elem())                  // session.User

	// fmt.Println(reflect.New(destType))                    // 创建空对象指针
	// fmt.Println(reflect.New(destType).Elem())             // 获取空对象value封装
	// fmt.Println(reflect.New(destType).Elem().Interface()) // 空对象转interface

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}

		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// fmt.Println(values) 反射得到的values切片是创建的dest各个字段的地址 [0xc0000a65c0 0xc0000a65d0]
		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// 支持Update(map[string]interface{}{"Name": "zgh", "Age": 18})
// 同时支持Update("Name", "zgh", "Age", 18, ...)
func (s *Session) Update(values ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	m, ok := values[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		length := len(values)
		for i := 0; i < length; i += 2 {
			m[values[i].(string)] = values[i+1]
		}
	}

	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	s.CallMethod(AfterUpdate, nil)
	return s.runSQL(sql, vars)
}

func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	s.CallMethod(AfterDelete, nil)
	return s.runSQL(sql, vars)
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)

	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// support chain
// "name=? and age>?", "devhui", 20
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// LIMIT ?
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// ORDERBY %s
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()

	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}

	if destSlice.Len() == 0 {
		return errors.New("not found")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
