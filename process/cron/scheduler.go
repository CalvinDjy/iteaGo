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

		name := p[TASK_KEY].(string)

		task := reflect.ValueOf(s.Ioc.InsByName(name))
		if !task.IsValid() {
			ilog.Error(fmt.Sprintf("task [%s] is nil, please check out if [%s] is registed", name, name))
			continue
		}
		
		method := task.MethodByName("Execute")
		if !method.IsValid() {
			ilog.Error(fmt.Sprintf("task [%s] need the method of `Execute`", name, name))
			continue
		}
		
		s.cron.AddFunc(p[CRON_KEY].(string), func() {
			method.Call([]reflect.Value{})
		})
	}
	
	s.cron.Start()

	ilog.Info("=== 【Scheduler】 Start ===")

	s.stop()
}

//Scheduler stop
func (s *Scheduler) stop() {
	for {
		select {
		case <-	s.Ctx.Done():
			ilog.Info("scheduler stop ...")
			s.cron.Stop()
			ilog.Info("scheduler stop success")
			return
		}
	}
}