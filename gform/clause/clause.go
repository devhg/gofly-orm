package clause

import "strings"

type Type int

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

func (c *Clause) Set(p Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}

	// 根据Type key从容器中获取生成器函数，调用不同的生成器函数，获得sql和vars
	sql, rVars := generators[p](vars...)
	c.sql[p] = sql
	c.sqlVars[p] = rVars
}

//根据传入序列生成 SELECT  WHERE ORDERBY...
//bug：复杂sql存在问题
func (c *Clause) Build(orders ...Type) (string, []interface{}) {

	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
