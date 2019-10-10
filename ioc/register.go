package ioc

import (
	"github.com/CalvinDjy/iteaGo/ioc/bean"
	"github.com/CalvinDjy/iteaGo/process/cron"
	"github.com/CalvinDjy/iteaGo/process/ihttp"
	"github.com/CalvinDjy/iteaGo/process/kafka"
	"github.com/CalvinDjy/iteaGo/process/thrift"
	"reflect"
	"strings"
)

const SINGLETON = "singleton"

type Register struct {

}

//Create register
func NewRegister() (c *Register) {
	return &Register{}
}

func (r *Register) process() []interface{} {
	return [] interface{}{
		ihttp.HttpServer{},
		thrift.ThriftServer{},
		cron.Scheduler{},
		kafka.KafkaConsumer{},
	}
}

func (r *Register) module() []interface{} {
	return [] interface{}{
		ihttp.Route{},
	}
}

//Init system beans
func (r *Register) Init() []*bean.Bean {
	return r.Register(append(r.process(), r.module()...))
}

//Register beans
func (r *Register) Register(class []interface{}) []*bean.Bean {
	var beans []*bean.Bean
	for _, b := range class {
		t := reflect.TypeOf(b)
		bean := &bean.Bean{
			Name: t.Name(),
			Scope: SINGLETON,
			Abstract: b,
			Concrete: b,
		}
		bean.SetAbstractType(t)
		bean.SetConcreteType(t)
		beans = append(beans, bean)
	}
	return beans
}

//Register beans
func (r *Register) RegisterBeans(beans []*bean.Bean) []*bean.Bean {
	for _, bean := range beans {
		if bean.Concrete == nil {
			panic("concrete of bean should not be nil")
		}

		tc := reflect.TypeOf(bean.Concrete)
		bean.SetConcreteType(tc)
		if bean.Abstract == nil {
			bean.Abstract = bean.Concrete
		}
		ta := reflect.TypeOf(bean.Abstract)
		bean.SetAbstractType(ta)

		if strings.EqualFold(bean.Name, "") {
			bean.Name = ta.Name()
		}

		if strings.EqualFold(bean.Scope, "") {
			bean.Scope = "singleton"
		}
	}
	return beans
}
