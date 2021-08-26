package session

import (
	"reflect"

	"github.com/devhg/gofly-orm/log"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

func (s *Session) CallMethod(name string, value interface{}) {
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(name)

	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(name)
	}

	params := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		if v := fm.Call(params); len(v) > 0 {
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}
