package cron

import (
	"context"
	"fmt"
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/CalvinDjy/iteaGo/ioc/iface"
	"github.com/robfig/cron"
	"reflect"
)

const (
	TASK_KEY = "Task"
	CRON_KEY = "Cron"
)

type Scheduler struct {
	Ctx             context.Context
	Ioc 			iface.IIoc
	Name			string
	Processor 		[]interface{}
	cron			*cron.Cron
}

func (s *Scheduler) Execute() {
	if len(s.Processor) == 0 {
		return
	}

	s.cron = cron.New()
	
	for _, process := range s.Processor {
		if _, ok := process.(map[interface{}]interface{}); !ok {
			continue
		}

		p := process.(map[interface{}]interface{})

		if _, ok := p[TASK_KEY]; !ok {
			continue
		}

		if _, ok := p[CRON_KEY]; !ok {
			continue
		}

		task := reflect.ValueOf(s.Ioc.InsByName(p[TASK_KEY].(string)))
		if !task.IsValid() {
			panic(fmt.Sprint("Controller [", p[TASK_KEY].(string), "] need register"))
		}
		
		method := task.MethodByName("Execute")
		if !method.IsValid() {
			continue
		}
		
		s.cron.AddFunc(p[CRON_KEY].(string), func() {
			method.Call([]reflect.Value{})
		})
	}

	go s.stop()
	
	s.cron.Start()

	ilog.Info("=== 【Scheduler】 Start ===")
}

//Scheduler stop
func (s *Scheduler) stop() {
	for {
		select {
		case <-	s.Ctx.Done():
			ilog.Info("Scheduler stop ...")
			s.cron.Stop()
			ilog.Info("Scheduler stop success")
			return
		}
	}
}