package session

import (
	"fmt"
	"github.com/QXQZX/gofly-orm/gform/log"
	"testing"
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
	s := NewSession_().Model(&Account{})
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
