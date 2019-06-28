package itea

import (
	"net/http"
	"reflect"
)

type IInterceptor interface {
	Enter(*http.Request) error
	Exit(*http.Request, *Response)
}

func ActionInterceptor(interceptors []string, ioc *Ioc) [][]reflect.Value {
	var list [][]reflect.Value
	for _, name := range interceptors {
		var t reflect.Type
		if _, ok := ioc.typeN[name]; !ok {
			continue
		}
		t = ioc.typeN[name]
		if !t.Implements(reflect.TypeOf(new(IInterceptor)).Elem()) {
			continue
		}
		ins := ioc.InsByType(t)
		if ins == nil {
			continue
		}
		v := reflect.ValueOf(ins)
		list = append(list, []reflect.Value{
			v.MethodByName("Enter"),
			v.MethodByName("Exit"),
		})
	}
	return list
}