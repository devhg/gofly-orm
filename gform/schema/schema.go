package schema

import (
	"go/ast"
	"reflect"

	"github.com/devhg/gofly-orm/dialect"
)

// Field represents a column of database
type Field struct {
	Name string
	Type string
	Tag  string
}

// 设计 Schema，利用反射(reflect)完成结构体和数据库表结构的映射，
// 包括表名、字段名、字段类型、字段 tag 等
// Schema represents a table of database
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	FieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.FieldMap[name]
}

// RecordValues 根据scheme对象结构字段信息，获取dest对象里对应字段的值
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

// Parse 利用反射(reflect)完成结构体和数据库表结构的映射，
// 包括表名、字段名、字段类型、 字段 tag 等
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// 为设计的入参是一个对象的指针，因此需要 reflect.Indirect() 获取指针指向的实例
	// dest结构体的信息
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	// fmt.Println(modelType.Name()) // User
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		FieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("gform"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.FieldMap[p.Name] = field
		}
	}
	return schema
}
