package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[VALUES] = _values
	generators[ORDERBY] = _orderBy
	generators[WHERE] = _where
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// INSERT INTO $tableName ($fields)
//tableName string, fields []string{}
func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

// SELECT ($fields) FROM $tableName
//string, []string{}
func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

// VALUES ($v1), ($v2), ...
// bug: 每一个value的参数个数不同， 如下
//_values([]interface{}{"dev", 18, "男"}, []interface{}{"dev", 18})
//VALUES (?, ?, ?), (?, ?, ?)
func _values(values ...interface{}) (string, []interface{}) {
	var bindString string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")

	for i, value := range values {
		v := value.([]interface{})
		if bindString == "" {
			bindString = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindString))

		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		// 所有的? 对应的值都被存到了vars
		vars = append(vars, v...)
	}

	return sql.String(), vars
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// return ?, ?, ?, ?
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	data := values[1].(map[string]interface{})

	var fields []string
	var vars []interface{}

	for key, val := range data {
		fields = append(fields, key+"=?")
		vars = append(vars, val)
	}

	return fmt.Sprintf("UPDATE %s SET %s", tableName,
		strings.Join(fields, ", ")), vars
}

func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
