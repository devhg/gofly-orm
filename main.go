package main

import (
	"fmt"
	"github.com/QXQZX/gofly-orm/gform"
	"log"

	// 导入时会注册 sqlite3 的驱动
	_ "github.com/mattn/go-sqlite3"
)

type Users struct {
	Name string `gform:"PRIMARY KEY"`
	Age  int
}

func main() {

	engine, _ := gform.NewEngine("sqlite3", "gofly.db")
	defer engine.Close()
	s := engine.NewSession()
	//_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	//_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	//_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	//result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	//count, _ := result.RowsAffected()
	//fmt.Printf("Exec success, %d affected\n", count)

	m := s.Model(&Users{})
	_ = m.DropTable()
	_ = m.CreateTable()
	if !m.HasTable() {
		log.Fatal("Failed to create table User")
	}
	fmt.Println(m.HasTable())
}
