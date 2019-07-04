package itea

import "reflect"

type Bean struct {
	Name string
	Scope string
	Abstract interface{}
	Concrete interface{}
	abstractType reflect.Type
	concreteType reflect.Type
}

func (b *Bean)setAbstractType(t reflect.Type) {
	b.abstractType = t
}

func (b *Bean)setConcreteType(t reflect.Type) {
	b.concreteType = t
}

func (b *Bean)getAbstractType() reflect.Type {
	return b.abstractType
}

func (b *Bean)getConcreteType() reflect.Type{
	return b.concreteType
}
