package session

import (
	"7days/ORM/clause"
	"7days/ORM/dialect"
	"7days/ORM/logg"
	"7days/ORM/schema"
	"database/sql"
	"fmt"
	"strings"
)

type Session struct {
	db       *sql.DB         //连接的数据库
	sql      strings.Builder //存储sql语句，后面进行操作
	sqlVars  []interface{}   //sql语句涉及的变量
	dialect  dialect.Dialect //不同软件
	refTable *schema.Schema  //正在操作的表
	clause   clause.Clause   //生成sql语句
	tx       *sql.Tx         //用于事务
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func NewSession(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

//用于验证是否开启了事务

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

//清空上一次的查询语句，以致于后面可以使用

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

//写入查询语句

func (s *Session) Raw(sql string, sqlVars ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, sqlVars...)
	return s
}

//封装 Exec, Query(), QueryRow() 方法

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	logg.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		fmt.Println("nn")
		logg.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	logg.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	logg.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		logg.Error(err)
	}
	return
}
