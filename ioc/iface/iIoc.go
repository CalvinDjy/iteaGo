package iface

import (
	"github.com/CalvinDjy/iteaGo/ioc/bean"
	"reflect"
)

type IIoc interface {
	InsByName(name string) interface{}
	InsByType(t reflect.Type) interface{}
	BeansByName(name string) *bean.Bean
}