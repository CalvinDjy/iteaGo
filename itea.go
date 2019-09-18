package itea

import (
	"context"
	"github.com/CalvinDjy/iteaGo/ilog"
	"os"
	"sync"
)

const ITEAGO_VERSION = "v0.4.5"

var (
	sigs 			chan os.Signal
	s				chan bool
	ctx				context.Context
	config			*Config
)

type Process struct {
	Name 			string
	Class 			string
	ExecuteMethod 	string
	Params 			map[string]interface{}
}

type Itea struct {
	process			[]interface{}
	ioc 			*Ioc
}

//Create Itea
func New(appConfig string) *Itea {
	config = InitConf(appConfig)
	ctx = context.Background()
	InitLog()
	if process := config.GetStructArray("application.process", Process{}); process != nil {
		return &Itea{
			process: process,
			ioc: NewIoc(),
		}
	} else {
		panic("Can not find config of process")
	}
}

//Debug
func (i *Itea) Debug() *Itea {
	ctx = context.WithValue(ctx, DEBUG, true)
	return i
}

//Register simple beans
func (i *Itea) Register(beans ...[]interface{}) *Itea {
	if i == nil {
		return nil
	}
	var beanList [] interface{}
	for _, bean := range beans{
		beanList = append(beanList, bean...)
	}
	i.ioc.Register(beanList)
	return i
}

//Register beans
func (i *Itea) RegisterBean(beans ...[]*Bean) *Itea {
	if i == nil {
		return nil
	}
	var beanList []*Bean
	for _, bean := range beans{
		beanList = append(beanList, bean...)
	}
	i.ioc.RegisterBeans(beanList)
	return i
}

//Run Itea
func (i *Itea) Run() {
	switch true {
	case Stop:
		i.stop()
		break
	case Help:
		break
	case Start:
		i.start()
		break
	}
}

//Start Itea
func (i *Itea) start() {
	go logProcessInfo()

	s = make(chan bool)
	defer close(s)

	go processSignal()

	ctx, stop := context.WithCancel(ctx)

	go func() {
		if <-s {
			ilog.Info("Itea stop ...")
			stop()
		}
	}()

	var wg sync.WaitGroup
	for _, p := range i.process {
		var process = p.(*Process)
		wg.Add(1)
		go func() {
			defer wg.Done()
			i.ioc.ExecProcess(ctx, process)
		}()
	}
	wg.Wait()

	ilog.Info("Itea stop success. Good bye ")
	
	if ilog.Done() {
		close(sigs)
		removePid()
		os.Exit(0)
	}

}

//Stop itea
func (i *Itea) stop() {
	stopProcess()
}

type IteaTest struct {
	Ioc 	*Ioc
}

//Create IteaTest
func NewIteaTest(appConfig string) *IteaTest {
	config = InitConf(appConfig)
	ctx = context.WithValue(context.Background(), DEBUG, true)
	InitLog()
	return &IteaTest{
		Ioc: NewIoc(),
	}
}
