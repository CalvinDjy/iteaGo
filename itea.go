package itea

import (
	"context"
	"github.com/CalvinDjy/iteaGo/ilog"
	"os"
	"sync"
)

var (
	wg 		sync.WaitGroup
	sigs 	chan os.Signal
	s		chan bool
	mutex 	*sync.Mutex
	ctx		context.Context
	conf	*Config
)

type Itea struct {
	beans 	[]Bean
	ioc 	*Ioc
}

//Create Itea
func New(appConfig string) *Itea {
	mutex = new(sync.Mutex)
	conf = InitConf(appConfig)
	ctx = context.WithValue(context.Background(), DEBUG, false)

	if process := conf.Beans(PROCESS_CONFIG); process != nil {
		return &Itea{
			beans: process,
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

//Register beans
func (i *Itea) Register(beans ...[] interface{}) *Itea {
	var beanList [] interface{}
	for _, bean := range beans{
		beanList = append(beanList, bean...)
	}
	i.ioc.RegisterBeans(beanList)
	return i
}

//Start Itea
func (i *Itea) Start() {
	InitLog()
	go logProcessInfo()

	s = make(chan bool)

	wg.Add(1)
	go func() {
		defer wg.Done()
		processSignal()
	}()

	ctx, stop := context.WithCancel(ctx)

	go func() {
		if <-s {
			ilog.Info("Itea stop ...")
			stop()
		}
	}()
	for _, p := range i.beans {
		var process = p
		wg.Add(1)
		go func() {
			defer wg.Done()
			i.ioc.InitProcess(ctx, process)
		}()
	}
	wg.Wait()

	ilog.Info("Itea stop success. Good bye ")

	if ilog.Done() {
		os.Exit(0)
	}

}

type IteaTest struct {
	Ioc 	*Ioc
}

//Create IteaTest
func NewIteaTest(appConfig string) *IteaTest {
	mutex = new(sync.Mutex)
	conf = InitConf(appConfig)
	ctx = context.WithValue(context.Background(), DEBUG, true)
	InitLog()
	return &IteaTest{
		Ioc: NewIoc(),
	}
}
