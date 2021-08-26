package session

import (
	"fmt"
	"testing"

	"github.com/devhg/gofly-orm/log"
)

type Account struct {
	ID  int `gform:"primary key"`
	Pwd string
}

func (a *Account) BeforeInsert(s *Session) error {
	log.Info("before insert")
	a.ID += 1000
	return nil
}

func (a *Account) AfterQuery(s *Session) {
	log.Info("after insert")
	a.Pwd = "******"
}

func TestSession_CallMethod(t *testing.T) {
	s := NewSession2().Model(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&Account{1, "asa"}, &Account{2, "qwe"})

	account := &Account{}

	err := s.First(account)
	if err != nil {
		panic(err)
	}
	fmt.Println(account)
}
