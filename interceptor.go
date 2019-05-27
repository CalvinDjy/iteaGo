package itea

import (
	"net/http"
	"reflect"
	"sort"
)

type IInterceptor interface {
	Enter(*http.Request) error
	Exit(*http.Request, *Response)
}

func GetInterceptor(ioc *Ioc) [][]reflect.Value {
	list := make(map[int][]reflect.Value)
	var key []int
	for _, t := range ioc.typeN {
		if t.Implements(reflect.TypeOf(new(IInterceptor)).Elem()) {
			v := reflect.ValueOf(ioc.InsByType(t))
			list[ioc.idT[t]] = []reflect.Value{
				v.MethodByName("Enter"),
				v.MethodByName("Exit"),
			}
			key = append(key, ioc.idT[t])
		}
	}

	sort.Ints(key)
	
	var ilist [][]reflect.Value
	for _, k := range key {
		ilist = append(ilist, list[k])
	}

	return ilist
}