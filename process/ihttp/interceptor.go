package ihttp

import (
	"github.com/CalvinDjy/iteaGo/ioc/iface"
	"net/http"
	"reflect"
)

type IInterceptor interface {
	Handle(func(*http.Request, *Response) error) func(*http.Request, *Response) error
}

func ActionInterceptor(interceptors []string, ioc iface.IIoc) []IInterceptor {
	var list []IInterceptor
	IType := reflect.TypeOf(new(IInterceptor)).Elem()
	l := len(interceptors)
	for i := l-1; i >= 0; i-- {
		name := interceptors[i]
		var t reflect.Type
		b := ioc.BeansByName(name)
		if b == nil {
			continue
		}
		t = b.GetConcreteType()
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