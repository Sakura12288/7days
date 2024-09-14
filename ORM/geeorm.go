package ORM

import (
	"7days/ORM/dialect"
	"7days/ORM/logg"
	"7days/ORM/session"
	"database/sql"
	"fmt"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

//driver 是数据库类型，source 是地址

func NewEngine(driver string, source string) (engine *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		logg.Error(err)
		return
	}
	//验证是否连接成功
	if err = db.Ping(); err != nil {
		logg.Error(err, "连接失败")
	}
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		logg.Errorf("dialect %s not found", driver)
	}
	engine = &Engine{db: db, dialect: dial}
	logg.Info("连接成功")
	return
}

//记得关闭数据库的连接，不然会浪费资源

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		logg.Error(err, "关闭失败")
	}
	logg.Info("成功关闭")
}

func (e *Engine) NewSess() *session.Session {
	return session.NewSession(e.db, e.dialect)
}

//提供事务接口

type TxFunc func(*session.Session) (interface{}, error)

func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSess()
	if err = s.Begin(); err != nil {
		logg.Error(err)
		return
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.RollBack()
			fmt.Println("11")
			panic(p) //回滚后再panic
		} else if err != nil {
			fmt.Println("13")
			//er := s.RollBack()
			//if er != nil {
			//	logg.Error(er, "mi")
			//}
		} else {
			fmt.Println("12")
			defer func() {
				if err != nil {
					_ = s.RollBack()
				}
			}()
			err = s.Commit()
		}
		fmt.Println(s.Count())
	}()
	return f(s)
}
