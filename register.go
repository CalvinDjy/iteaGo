package itea

import (
	"reflect"
	"strings"
)

const SINGLETON = "singleton"

type Register struct {

}

func process() []interface{} {
	return [] interface{}{
		HttpServer{},
		ThriftServer{},
		Scheduler{},
	}
}

func module() []interface{} {
	return [] interface{}{
		Route{},
	}
}

//Create register
func NewRegister() (c *Register) {
	return &Register{}
}

//Init system beans
func (r *Register) Init() []*Bean {
	return r.Register(append(process(), module()...))
}

//Register beans
func (r *Register) Register(class []interface{}) []*Bean {
	var beans []*Bean
	for _, b := range class {
		t := reflect.TypeOf(b)
		bean := &Bean{
			Name: t.Name(),
			Scope: SINGLETON,
			Abstract: b,
			Concrete: b,
		}
		bean.setAbstractType(t)
		bean.setConcreteType(t)
		beans = append(beans, bean)
	}
	return beans
}

//Register beans
func (r *Register) RegisterBeans(beans []*Bean) []*Bean {
	for _, bean := range beans {
		if bean.Concrete == nil {
			panic("concrete of bean should not be nil")
		}

		tc := reflect.TypeOf(bean.Concrete)
		bean.setConcreteType(tc)
		if bean.Abstract == nil {
			bean.Abstract = bean.Concrete
		}
		ta := reflect.TypeOf(bean.Abstract)
		bean.setAbstractType(ta)

		if strings.EqualFold(bean.Name, "") {
			bean.Name = ta.Name()
		}

		if strings.EqualFold(bean.Scope, "") {
			bean.Scope = "singleton"
		}
	}
	return beans
}
