package session

import (
	"7days/ORM/dialect"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testing"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("mysql")
)

func TestMain(m *testing.M) {
	TestDB, _ = sql.Open("mysql", "root:123456@/orm?charset=utf8")
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSessi() *Session {
	return NewSession(TestDB, TestDial)
}

func TestSession_Exec(t *testing.T) {
	s := NewSessi()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(Name) values (?), (?)", "Tom", "Sam").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRows(t *testing.T) {
	s := NewSessi()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	row := s.Raw("SELECT count(*) FROM User").QueryRow()
	var count int
	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}
