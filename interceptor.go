package itea

import (
	"net/http"
	"reflect"
)

type IInterceptor interface {
	Enter(*http.Request) error
	Exit(*http.Request, *interface{}) func()
}

type InterceptorManager struct {
	ioc *Ioc
}

func NewInterceptorManager(ioc *Ioc) *InterceptorManager {
	return &InterceptorManager{
		ioc: ioc,
	}
}

func (im *InterceptorManager) GetInterceptor() [][]reflect.Value {
	var ilist [][]reflect.Value
	for _, t := range im.ioc.typeN {
		if t.Implements(reflect.TypeOf(new(IInterceptor)).Elem()) {
			v := reflect.ValueOf(im.ioc.GetInstanceByType(t))
			ilist = append(ilist, []reflect.Value{
				v.MethodByName("Enter"),
				v.MethodByName("Exit"),
			})
		}
	}
	return ilist
}