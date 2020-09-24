package clause

import (
	"fmt"
	"testing"
)

func Test__insert(t *testing.T) {
	_insert("2323", []string{"1212", "sdsd"})
}

func testVals(vals ...interface{}) {
	//fmt.Println(vals)
}

func TestMY(t *testing.T) {
	testVals()
}

func Test_genBindVars(t *testing.T) {
	vars := genBindVars(4)
	fmt.Println(vars)
}

func Test__values(t *testing.T) {
	values, i := _values([]interface{}{"dev", 18, "男"}, []interface{}{"dev", 18})
	fmt.Println(values)
	fmt.Println(i)
}
