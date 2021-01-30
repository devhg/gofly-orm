package dialect

import "reflect"

//为适配不同的数据库，映射数据类型和特定的 SQL 语句，
//创建 Dialect 层结局数据库之间差异。
//day2

var dialectMap = map[string]Dialect{}

type Dialect interface {
	//用于将 Go 语言的类型转换为该数据库的数据类型。
	DataTypeOf(typ reflect.Value) string
	//返回某个表是否存在的sql语句，参数为表名
	TableExistSql(table string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
	dialect, ok := dialectMap[name]
	return dialect, ok
}
