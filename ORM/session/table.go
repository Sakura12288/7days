package session

import (
	"7days/ORM/logg"
	"7days/ORM/schema"
	"fmt"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session { //设置或更新表结构
	// nil or different model, update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		logg.Error("Model is not set")
	}
	return s.refTable
}
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE if not exists %s (%s);", table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return judge(tmp, s.RefTable().Name) //注意mysql不区分大小写
}
func judge(s1, s2 string) bool { //不区分大小写的判断是否相等
	t1 := strings.ToLower(s1)
	t2 := strings.ToLower(s2)
	return t1 == t2
}
