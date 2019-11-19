package thrift

import (
	"context"
	"fmt"
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/CalvinDjy/iteaGo/ioc/iface"
	"github.com/apache/thrift/lib/go/thrift"
)

type ThriftServer struct {
	Ctx             context.Context
	Ioc 			iface.IIoc
	Name   			string
	Ip				string
	Port 			int
	Multiplexed		bool
	Processor 		[]interface{}
	ser 			*thrift.TSimpleServer
}

//Thrift server start
func (ts *ThriftServer) Execute() {

	addr := fmt.Sprintf("%s:%d", ts.Ip, ts.Port)

	serverTransport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		ilog.Info(err)
		panic(err)
	}
	
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	ts.ser = thrift.NewTSimpleServer4(ts.processor(), serverTransport, transportFactory, protocolFactory)
	
	go ts.stop()

	ilog.Info(fmt.Sprintf("=== 【Thrift】Server [%s] start [%s] ===", ts.Name, addr))
	if err = ts.ser.Serve(); err != nil {
		ilog.Error(err)
		return
	}
}

//Thrift processor
func (ts *ThriftServer) processor() thrift.TProcessor {
	if ts.Multiplexed {
		processor := thrift.NewTMultiplexedProcessor()
		for _, v := range ts.Processor {
			if p := ts.check(v.(string)); p != nil {
				processor.RegisterProcessor(p.Name(), p.Processor())
				ilog.Info(fmt.Sprintf("... 【Thrift】Register processor [%s] multiplexed", p.Name()))
			}
		}
		return processor
	} else {
		if ts.Processor != nil && len(ts.Processor) > 0 {
			if p := ts.check(ts.Processor[0].(string)); p != nil {
				processor := p.Processor()
				ilog.Info(fmt.Sprintf("... 【Thrift】Register processor [%s]", p.Name()))
				return processor
			}
		}
		panic("thrift processor config error")
	}
}

func (ts *ThriftServer) check(name string) IProcessor {
	i := ts.Ioc.InsByName(name)
	if i == nil {
		ilog.Error(fmt.Sprintf("processor [%s] is nil, please check out if [%s] is registed", name, name))
		return nil
	}
	
	if p, ok := i.(IProcessor); ok {
		return p
	} else {
		ilog.Error(fmt.Sprintf("processor [%s] is not impliment of thrift.IProcessor", i))
	}
	
	return nil
}

//Thrift server stop
func (ts *ThriftServer) stop() {
	for {
		select {
		case <-	ts.Ctx.Done():
			ilog.Info("thrift server stop ...")
			ts.ser.Stop()
			ilog.Info("thrift server stop success")
			return
		}
	}
}