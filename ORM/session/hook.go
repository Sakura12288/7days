package session

import (
	"7days/ORM/logg"
	"reflect"
)

// 以*session作为参数，在变化前后进行一些行为
// CallMethod calls the registered hooks
func (s *Session) CallMethod(method string, value interface{}) (err error) {
	val := reflect.ValueOf(s.RefTable().Model).Interface()
	if value != nil {
		val = reflect.ValueOf(value).Interface()
	}
	switch method {
	case "BeforeQuery":
		v, ok := val.(BeforeQuery)
		if !ok {
			break
		}
		err = v.BeforeQuery(s)
	case "AfterQuery":
		v, ok := val.(AfterQuery)
		if !ok {
			break
		}
		err = v.AfterQuery(s)
	case "BeforeInsert":
		v, ok := val.(BeforeInsert)
		if !ok {
			break
		}
		err = v.BeforeInsert(s)
	case "AfterInsert":
		v, ok := val.(AfterInsert)
		if !ok {
			break
		}
		err = v.AfterInsert(s)
	case "BeforeUpdate":
		v, ok := val.(BeforeUpdate)
		if !ok {
			break
		}
		err = v.BeforeUpdate(s)
	case "AfterUpdate":
		v, ok := val.(AfterUpdate)
		if !ok {
			break
		}
		err = v.AfterUpdate(s)
	case "BeforeDelete":
		v, ok := val.(BeforeDelete)
		if !ok {
			break
		}
		err = v.BeforeDelete(s)
	case "AfterDelete":
		v, ok := val.(AfterDelete)
		if !ok {
			break
		}
		err = v.AfterDelete(s)
	default:
		logg.Error("钩子未注册")
	}
	if err != nil {
		logg.Error(err)
	}
	return
}
