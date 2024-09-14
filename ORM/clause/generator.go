package clause

import (
	"fmt"
	"strings"
)

//根据值生成相应的语句段，values有些时候形式是表名+值
//统一格式，后一个值不一定有实际值

type generator func(values ...interface{}) (string, []interface{})

//根据key映射到相应的构造器

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[WHERE] = _where
	generators[SELECT] = _select
	generators[VALUES] = _values
	generators[ORDERBY] = _orderBy
	generators[LIMIT] = _limit
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}
func genBindVars(num int) string {
	vars := make([]string, num)
	for i := 0; i < num; i++ {
		vars[i] = "?"
	}
	return strings.Join(vars, ",")
}

// values 表名加字段，可以 s吗？
func _insert(values ...interface{}) (string, []interface{}) {
	table := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("insert into %s (%s)", table, fields), []interface{}{}
}

func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("where %s", desc), vars
}
func _limit(values ...interface{}) (string, []interface{}) {
	return "Limit ?", values
}
func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("order by %s", values[0]), []interface{}{}
}
func _select(values ...interface{}) (string, []interface{}) {
	tablename := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("select %v from %s", fields, tablename), []interface{}{}
}

//生成插入时用的值，可能是多条插入值

func _values(values ...interface{}) (string, []interface{}) {
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars

}

//传入表名和键值对

func _update(values ...interface{}) (string, []interface{}) {
	table := values[0]
	vars := values[1].(map[string]interface{})
	keys := []string{}
	va := []interface{}{}
	for key, val := range vars {
		keys = append(keys, key+" = ?")
		va = append(va, val)
	}
	return fmt.Sprintf("update %s set %v", table, strings.Join(keys, ",")), va
}

//delete只需要表名即可，其他的在where语句

func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("delete from %v", values[0]), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
