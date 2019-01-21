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
func (r *Register)Init() []reflect.Type {
	return r.Register(append(process(), module()...))
}

//Register beans
func (r *Register)Register(beans [] interface{}) [] reflect.Type {
	var l []reflect.Type
	for _, v := range beans {
		l = append(l, reflect.TypeOf(v))
	}
	return l
}
