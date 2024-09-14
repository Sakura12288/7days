package RPC

import (
	"go/ast"
	"log"
	"reflect"
	"sync/atomic"
)

//将服务类型与方法体分开

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	numCalls  uint64
}

func (m *methodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}

//构造参数类型的反射零值

func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value
	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem())
	} else {
		argv = reflect.New(m.ArgType).Elem()
	}
	return argv
}

func (m *methodType) newReplyv() reflect.Value {
	replyv := reflect.New(m.ReplyType.Elem())
	//针对map和切片类型，需要初始化
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return replyv
}

//service 的定义也是非常简洁的，name 即映射的结构体的名称，比如 T，比如 WaitGroup；
//typ 是结构体的类型；rcvr 即结构体的实例本身，保留 rcvr 是因为在调用时需要 rcvr 作为第 0 个参数；
//method 是 map 类型，存储映射的结构体的所有符合条件的方法

//每一个service对应一个具体的类型及调用，后面再用一个把所有的service合在一起

type service struct {
	name   string
	typ    reflect.Type
	rcvr   reflect.Value //receiver
	method map[string]*methodType
}

func newService(rcvr interface{}) *service {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	//注意name与typ取法区别，名字需要把指针去掉，即&m与m对应一个名字，但类型不行
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	s.typ = reflect.TypeOf(rcvr) //注意这里只能用rcvr, 不要用s.rcvr,s的那个已经是reflect.Value类型了
	if !ast.IsExported(s.name) {
		log.Fatalf("%s 不可导出", s.name)
	}
	s.registerMethods()
	return s
}

//只把符合rpc格式的方法存储

func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		//注意第一个入参为结构体
		mtype := method.Type
		if mtype.NumIn() != 3 || mtype.NumOut() != 1 {
			continue
		}
		if mtype.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mtype.In(1), mtype.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Printf("rpc server: register %s.%s\n", s.name, method.Name)
	}
}

// 判断是否是导出类型或者内置类型
func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

//调用

func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}
