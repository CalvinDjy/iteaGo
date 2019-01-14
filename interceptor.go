package itea

import (
	"net/http"
	"reflect"
)

type IInterceptor interface {
	Before(*http.Request, *reflect.Value) Error
	After(*http.Request, *reflect.Value) func()
}

type InterceptorManager struct {
	ioc *Ioc
}

func NewInterceptorManager(ioc *Ioc) *InterceptorManager {
	return &InterceptorManager{
		ioc: ioc,
	}
}

func (im *InterceptorManager) GetInterceptor() []map[string]reflect.Value {
	interceptorList := make([]map[string]reflect.Value, len(im.ioc.interceptorType))
	i := 0
	for _, t := range im.ioc.interceptorType {
		interceptInsVal := reflect.ValueOf(im.ioc.GetInstanceByType(t))
		interceptorList[i] = map[string]reflect.Value{
			"before": interceptInsVal.MethodByName("Before"),
			"after": interceptInsVal.MethodByName("After"),
		}
		i++
	}
	return interceptorList
}