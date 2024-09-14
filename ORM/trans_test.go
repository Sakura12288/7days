package ORM

import (
	"7days/ORM/session"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("mysql", "root:123456@/orm?charset=utf8")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	//t.Run("rollback", func(t *testing.T) {
	//	transactionRollback(t)
	//})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}
func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	ss := engine.NewSess()
	_ = ss.Model(&User{}).DropTable()
	_ = ss.Model(&User{}).CreateTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{})
		_, err = s.Insert(&User{"Tom", 18})
		fmt.Println(s.Count())
		return nil, errors.New("Error")
	})
	//|| ss.HasTable()
	if err == nil {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	ss := engine.NewSess()
	_ = ss.Model(&User{}).DropTable()
	_ = ss.Model(&User{}).CreateTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{})
		_, err = s.Insert(&User{"Tom", 18})
		return
	})
	u := &User{}
	_ = ss.First(u)
	fmt.Println(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}
