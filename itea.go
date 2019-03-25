package itea

import (
	"sync"
	"os"
	"context"
	"log"
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
	ctx = context.WithValue(context.Background(), ENV, env())
	conf = InitConf(ctx, appConfig)

	wg.Add(2)
	go dbConfig()
	go importConfig()
	wg.Wait()

	if process := conf.Beans(PROCESS_CONFIG); process != nil {
		ctx = context.WithValue(ctx, DEBUG, false)
		return &Itea{
			beans: process,
			ioc: NewIoc(ctx),
		}
	} else {
		panic("Can not find config of process")
	}

}

//Env
func env() string {
	num := len(os.Args)
	if num > 1 {
		return os.Args[1]
	}
	return DEFAULT_ENV
}

//Extract db config
func dbConfig() {
	if d := conf.Path(DB_CONFIG); d != "" {
		mutex.Lock()
		ctx = context.WithValue(ctx, DB_CONFIG, d)
		mutex.Unlock()
	}
	wg.Done()
}

//Extract import config
func importConfig() {
	if i := conf.PathList(IMPORT_CONFIG); i != nil {
		mutex.Lock()
		ctx = context.WithValue(ctx, IMPORT_CONFIG, i)
		mutex.Unlock()
	}
	wg.Done()
}

//Debug
func (i *Itea) Debug() *Itea {
	i.ioc.ctx = context.WithValue(i.ioc.ctx, DEBUG, true)
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
	go logProcessInfo()

	s = make(chan bool)

	wg.Add(1)
	go func() {
		processSignal()
		wg.Done()
	}()

	ctx, stop := context.WithCancel(i.ioc.ctx)
	go func() {
		if <-s {
			log.Println("Itea stop ...")
			stop()
		}
	}()
	for _, p := range i.beans {
		var process = p
		wg.Add(1)
		go func() {
			i.ioc.InitProcess(ctx, process)
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("Itea stop success. Good bye ")
	os.Exit(0)
}
