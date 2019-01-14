package itea

import (
	"encoding/json"
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
)

type Itea struct {
	config 	map[string]*json.RawMessage
	beans 	[]Bean
	ioc 	*Ioc
}

//Create Itea
func New(appConfig string) *Itea {

	conf := InitConf(appConfig)
	ctx := context.Background()
	mutex = new(sync.Mutex)

	wg.Add(2)
	go func() {
		if d := conf.Path(DB_CONFIG); d != "" {
			mutex.Lock()
			ctx = context.WithValue(ctx, DB_CONFIG, d)
			mutex.Unlock()
		}
		wg.Done()
	}()
	go func() {
		if i := conf.PathList(IMPORT_CONFIG); i != nil {
			mutex.Lock()
			ctx = context.WithValue(ctx, IMPORT_CONFIG, i)
			mutex.Unlock()
		}
		wg.Done()
	}()
	wg.Wait()

	if process := conf.Beans(PROCESS_CONFIG); process != nil {
		return &Itea{
			beans: process,
			ioc: NewIoc(ctx),
		}
	} else {
		panic("Can not find config of process")
	}

}

//Register beans
func (i *Itea) Register(beans ...[] interface{}) {
	var beanList [] interface{}
	for _, bean := range beans{
		beanList = append(beanList, bean...)
	}
	i.ioc.RegisterBeans(beanList)
}

//Start Itea
func (i *Itea) Start() {
	go logProcessInfo()
	sigs = make(chan os.Signal)
	s = make(chan bool)

	wg.Add(1)
	go func() {
		processSignal()
		wg.Done()
	}()

	ctx, stop := context.WithCancel(context.Background())
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
