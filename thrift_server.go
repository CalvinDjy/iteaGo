package itea

import (
	"github.com/CalvinDjy/iteaGo/ilog"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type ThriftServer struct {
	Name   			string
	Ip				string
	Port 			string
	Multiplexed		bool
	Processor 		[]interface{}
	Ctx             context.Context
	Ioc 			*Ioc
	ser 			*thrift.TSimpleServer
}

//Thrift server start
func (ts *ThriftServer) Execute() {

	addr := fmt.Sprintf("%s:%s", ts.Ip, ts.Port)

	serverTransport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		ilog.Info(err)
		panic(err)
	}
	
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	ts.ser = thrift.NewTSimpleServer4(ts.processor(), serverTransport, transportFactory, protocolFactory)
	
	go ts.stop()

	ilog.Info("=== Thrift server [", ts.Name, "] start [", addr, "] ===")
	err = ts.ser.Serve()
	if err != nil {
		ilog.Error(err)
		return
	}
	
}

//Thrift processor
func (ts *ThriftServer) processor() thrift.TProcessor {
	if ts.Multiplexed {
		processor := thrift.NewTMultiplexedProcessor()
		for _, v := range ts.Processor {
			if p, ok := ts.Ioc.InsByName(v.(string)).(IProcessor); ok {
				processor.RegisterProcessor(p.Name(), p.Processor())
				ilog.Info("--- Register multiplexed thrift processor [", p.Name(), "] ---")
			}
		}
		return processor
	} else {
		if ts.Processor != nil && len(ts.Processor) > 0 {
			if p, ok := ts.Ioc.InsByName(ts.Processor[0].(string)).(IProcessor); ok {
				processor := p.Processor()
				ilog.Info("--- Register thrift processor [", p.Name(), "] ---")
				return processor
			} 
		}
		ilog.Info("Thrift processor config error")
		panic("Thrift processor config error")
	}
}

//Thrift server stop
func (ts *ThriftServer) stop() {
	for {
		select {
		case <-	ts.Ctx.Done():
			ilog.Info("Thrift server stop ...")
			ts.ser.Stop()
			ilog.Info("Thrift server stop success")
			return
		default:
			break
		}
	}
}