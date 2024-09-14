package dialect

import "reflect"

//用来尽量实现对不同数据库类型的解耦，例如mysql，redis等

var dialectsMap = make(map[string]Dialect)

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []interface{})
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

func GetDialect(name string) (Dialect, bool) {
	dia, ok := dialectsMap[name]
	return dia, ok
}
