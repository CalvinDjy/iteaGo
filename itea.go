package itea

import (
	"context"
	"github.com/CalvinDjy/iteaGo/ilog"
	"os"
	"sync"
	"fmt"
)

var (
	wg 				sync.WaitGroup
	sigs 			chan os.Signal
	s				chan bool
	mutex 			*sync.Mutex
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
	mutex = new(sync.Mutex)
	config = InitConf(appConfig)
	ctx = context.WithValue(context.Background(), DEBUG, false)
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

//Get environment
func Env() string {
	if config == nil {
		panic("Please init itea")
	}
	return config.Env
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
	num := len(os.Args)
	if num <= 1 {
		i.start()
	}
	switch os.Args[1] {
	case "start":
		i.start()
	case "stop":
		i.stop()
	default:
		fmt.Println("error cmd")
	}
}

//Start Itea
func (i *Itea) start() {
	go logProcessInfo()

	s = make(chan bool)
	go processSignal()

	ctx, stop := context.WithCancel(ctx)

	go func() {
		if <-s {
			ilog.Info("Itea stop ...")
			stop()
		}
	}()
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
	mutex = new(sync.Mutex)
	config = InitConf(appConfig)
	ctx = context.WithValue(context.Background(), DEBUG, true)
	InitLog()
	return &IteaTest{
		Ioc: NewIoc(),
	}
}
