package bean

import "reflect"

type Bean struct {
	Name string
	Scope string
	Abstract interface{}
	Concrete interface{}
	abstractType reflect.Type
	concreteType reflect.Type
}

func (b *Bean)SetAbstractType(t reflect.Type) {
	b.abstractType = t
}

func (b *Bean)SetConcreteType(t reflect.Type) {
	b.concreteType = t
}

func (b *Bean)GetAbstractType() reflect.Type {
	return b.abstractType
}

func (b *Bean)GetConcreteType() reflect.Type{
	return b.concreteType
}
