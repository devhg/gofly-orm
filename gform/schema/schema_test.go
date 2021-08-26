package schema

import (
	"fmt"
	"testing"

	"github.com/devhg/gofly-orm/dialect"
)

// schema_test.go
type User struct {
	Name string `gform:"PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("sqlite3")

func TestParse(t *testing.T) {
	schema := Parse(&User{}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}

func TestSchema_RecordValues(t *testing.T) {
	u := User{
		Name: "sss",
		Age:  0,
	}
	schema := Parse(&User{}, TestDial)
	values := schema.RecordValues(u)

	fmt.Println(values)
}
