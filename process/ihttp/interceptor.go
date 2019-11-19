package ihttp

import (
	"fmt"
	"github.com/CalvinDjy/iteaGo/ilog"
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
			ilog.Error(fmt.Sprintf("can not find beans of [%s]", name))
			continue
		}
		t = b.GetConcreteType()
		if !t.Implements(IType) {
			ilog.Error(fmt.Sprintf("interceptor [%s] is not impliment of ihttp.IInterceptor", name))
			continue
		}
		ins := ioc.InsByType(t)
		if ins == nil {
			ilog.Error(fmt.Sprintf("interceptor [%s] is nil, please check out if [%s] is registed", name, name))
			continue
		}
		list = append(list, ins.(IInterceptor))
	}
	return list
}