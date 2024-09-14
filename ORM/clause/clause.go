package clause

import "strings"

//用来构建常用的sql语句，进行封装，方便直接调用
//对generator生成的子句进行拼接

type Type int

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

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

//Set 方法根据 Type 调用对应的 generator，生成该子句对应的 SQL 语句,然后保存在clause
//Build 方法根据传入的 Type 的顺序，构造出最终的 SQL 语句。

func (c *Clause) Set(name Type, values ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
	}
	if c.sqlVars == nil {
		c.sqlVars = make(map[Type][]interface{})
	}
	gen := generators[name]
	c.sql[name], c.sqlVars[name] = gen(values...)
}

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
