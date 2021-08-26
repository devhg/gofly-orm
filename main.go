package main

import (
	"fmt"
	"log"

	gform "github.com/devhg/gofly-orm"
	"github.com/devhg/gofly-orm/session"
)

type User struct {
	Name string `gform:"PRIMARY KEY"`
	Age  int
}

var (
	user1 = &User{"Tom", 18}
	user2 = &User{"Sam", 25}
	user3 = &User{"Jack", 25}
)

func main() {
	engine, _ := gform.NewEngine("sqlite3", "gofly.db")
	defer engine.Close()
	s := engine.NewSession().Model(&User{})

	testInsert(s)
	// test_Limit(s)
	// testUpdateFirstOrder(s)
	testDeleteAndCount(s)
}

func testInsert(s *session.Session) {
	err1 := s.DropTable()
	err2 := s.CreateTable()

	rows, err3 := s.Insert(user1, user2, user3)

	if err1 != nil || err2 != nil || err3 != nil {
		log.Fatal("failed init test records")
	}
	fmt.Println("插入成功:", rows)

	var users []User
	if err := s.Find(&users); err != nil {
		log.Fatal("failed to query all")
	}
	fmt.Println(users)
}

func testLimit(s *session.Session) {
	var users []User
	err := s.Limit(2).Find(&users)
	if err != nil || len(users) != 2 {
		log.Fatal("failed to query with limit condition")
	}
	fmt.Println(users)
}

func testUpdateFirstOrder(s *session.Session) {
	affected, _ := s.Where("Name = ?", "Tom").Update("Age", 310)
	u := &User{}
	_ = s.OrderBy("Age DESC").First(u)

	if affected != 1 || u.Age != 310 {
		log.Fatal("failed to update")
	}
	fmt.Println(u)
}

func testDeleteAndCount(s *session.Session) {
	affected, _ := s.Where("Name = ?", "Tom").Delete()
	fmt.Println(affected)
	count, _ := s.Count()

	if affected != 1 || count != 2 {
		log.Fatal("failed to delete or count")
	}
	fmt.Println(count)
}
