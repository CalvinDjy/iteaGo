package itea

import (
	"log"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type ThriftServer struct {
	Name   			string
	Ip				string
	Port 			string
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
		log.Println(err)
		panic(err)
	}
	
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	
	multiProcess := thrift.NewTMultiplexedProcessor()

	for _, v := range ts.Processor {
		processor := ts.Ioc.GetInstanceByName(v.(string))
		if _, ok := processor.(IProcessor); ok {
			p := processor.(IProcessor)
			multiProcess.RegisterProcessor(p.Name(), p.Processor())
			log.Println("--- Register thrift processor [", p.Name(), "] ---")
		}
	}

	ts.ser = thrift.NewTSimpleServer4(multiProcess, serverTransport, transportFactory, protocolFactory)
	
	go ts.stop()
	
	log.Print("=== Thrift server [", ts.Name, "] start [", addr, "] ===")
	err = ts.ser.Serve()
	if err != nil {
		log.Println(err)
		return
	}
	
}

//Thrift server stop
func (ts *ThriftServer)stop() {
	for {
		select {
		case <-	ts.Ctx.Done():
			log.Println("Thrift server stop ...")
			ts.ser.Stop()
			log.Println("Thrift server stop success")
			return
		default:
			break
		}
	}
}