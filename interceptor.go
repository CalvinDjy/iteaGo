package itea

import (
	"net/http"
	"reflect"
)

type IInterceptor interface {
	Handle(func(*http.Request, *Response) error) func(*http.Request, *Response) error
}

func ActionInterceptor(interceptors []string, ioc *Ioc) []IInterceptor {
	var list []IInterceptor
	IType := reflect.TypeOf(new(IInterceptor)).Elem()
	l := len(interceptors)
	for i := l-1; i >= 0; i-- {
		name := interceptors[i]
		var t reflect.Type
		if _, ok := ioc.beansN[name]; !ok {
			continue
		}
		t = ioc.beansN[name].getConcreteType()
		if !t.Implements(IType) {
			continue
		}
		ins := ioc.InsByType(t)
		if ins == nil {
			continue
		}
		list = append(list, ins.(IInterceptor))
	}
	return list
}