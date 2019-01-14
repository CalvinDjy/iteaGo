package itea

import (
	"reflect"
)

type Register struct {

}

func process() []interface{} {
	return [] interface{}{
		HttpServer{},
	}
}

func module() []interface{} {
	return [] interface{}{

	}
}

//Create register
func NewRegister() (c *Register) {
	return &Register{}
}

//Init system beans
func (r *Register)Init() (map[string] reflect.Type, map[string] reflect.Type){
	m, i := make(map[string] reflect.Type), make(map[string] reflect.Type)
	list := append(process(), module()...)
	var t reflect.Type
	for _, v := range list {
		t = reflect.TypeOf(v)
		m[t.Name()] = t
		if t.Implements(reflect.TypeOf(new(IInterceptor)).Elem()) {
			i[t.Name()] = t
		}
	}
	return m, i
}

//Register beans
func (r *Register)Register(beans [] interface{}) (map[string] reflect.Type, map[string] reflect.Type){
	m, i := make(map[string] reflect.Type), make(map[string] reflect.Type)
	var t reflect.Type
	for _, v := range beans {
		t = reflect.TypeOf(v)
		m[t.Name()] = t
		if t.Implements(reflect.TypeOf(new(IInterceptor)).Elem()) {
			i[t.Name()] = t
		}
	}
	return m, i
}
