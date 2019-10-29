package itea

import (
	"context"
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/CalvinDjy/iteaGo/ioc"
	"github.com/CalvinDjy/iteaGo/ioc/bean"
	"github.com/CalvinDjy/iteaGo/system"
	"github.com/CalvinDjy/iteaGo/process"
	"github.com/CalvinDjy/iteaGo/signal"
	"github.com/CalvinDjy/iteaGo/constant"
	"os"
	"sync"
)

var (
	sigs 			chan os.Signal
	s				chan bool
	ctx				context.Context
	config			*system.Config
)

type Itea struct {
	process			[]interface{}
	ioc 			*ioc.Ioc
}

//Create Itea
func New(appConfig string, debug bool) *Itea {
	system.InitConf(appConfig)
	ctx = context.Background()
	if debug {
		ctx = context.WithValue(ctx, constant.DEBUG, true)
	}
	system.InitLog()
	if process := system.Conf.GetStructArray("application.process", process.Process{}); process != nil {
		return &Itea{
			process: process,
			ioc: ioc.NewIoc(ctx),
		}
	} else {
		panic("Can not find config of process or process is nil")
	}
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
func (i *Itea) RegisterBean(beans ...[]*bean.Bean) *Itea {
	if i == nil {
		return nil
	}
	var beanList []*bean.Bean
	for _, bean := range beans{
		beanList = append(beanList, bean...)
	}
	i.ioc.RegisterBeans(beanList)
	return i
}

//Run Itea
func (i *Itea) Run() {
	switch true {
	case system.Stop:
		i.stop()
		break
	case system.Help:
		break
	case system.Start:
		i.start()
		break
	}
}

//Start Itea
func (i *Itea) start() {
	go signal.LogProcessInfo()

	s = make(chan bool)
	defer close(s)

	sigs = make(chan os.Signal)
	go signal.ProcessSignal(sigs, s)

	ctx, stop := context.WithCancel(ctx)

	go func() {
		if <-s {
			ilog.Info("Itea stop ...")
			stop()
		}
	}()

	var wg sync.WaitGroup
	for _, p := range i.process {
		var process = p.(*process.Process)
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
		signal.RemovePid()
		os.Exit(0)
	}

}

//Stop itea
func (i *Itea) stop() {
	signal.StopProcess()
}

//type IteaTest struct {
////	Ioc 	*ioc.Ioc
////}
////
//////Create IteaTest
////func NewIteaTest(appConfig string) *IteaTest {
////	config = system.InitConf(appConfig)
////	ctx = context.WithValue(context.Background(), DEBUG, true)
////	system.InitLog()
////	return &IteaTest{
////		Ioc: ioc.NewIoc(),
////	}
////}
