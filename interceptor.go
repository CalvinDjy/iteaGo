package itea

import (
	"net/http"
	"reflect"
	"sort"
)

type IInterceptor interface {
	Enter(*http.Request) error
	Exit(*http.Request, http.ResponseWriter, *interface{}) func()
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
	list := make(map[int][]reflect.Value)
	var key []int
	for _, t := range im.ioc.typeN {
		if t.Implements(reflect.TypeOf(new(IInterceptor)).Elem()) {
			v := reflect.ValueOf(im.ioc.GetInstanceByType(t))
			list[im.ioc.idT[t]] = []reflect.Value{
				v.MethodByName("Enter"),
				v.MethodByName("Exit"),
			}
			key = append(key, im.ioc.idT[t])
		}
	}

	sort.Ints(key)
	
	var ilist [][]reflect.Value
	for _, k := range key {
		ilist = append(ilist, list[k])
	}

	return ilist
}