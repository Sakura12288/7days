package main

import (
	"7days/ORM"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	engine, _ := ORM.NewEngine("mysql", "root:123456@/orm?charset=utf8")
	defer engine.Close()
	s := engine.NewSess()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name varchar(20));").Exec()
	_, _ = s.Raw("CREATE TABLE  IF NOT EXISTS User(Name varchar(20));").Exec()
	result, _ := s.Raw("INSERT into User values(?),(?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
