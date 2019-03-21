package itea

import (
	"log"
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
)

type ThriftServer struct {
	Ip				string
	Port 			string
	Processor 		[]interface{}
	Ctx             context.Context
	Ioc 			*Ioc
}

func (ts *ThriftServer) Execute() {

	addr := fmt.Sprintf("%s:%s", ts.Ip, ts.Port)

	serverTransport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		log.Println(err)
		return
	}
	
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	
	multiProcess := thrift.NewTMultiplexedProcessor()

	for _, v := range ts.Processor {
		processor := ts.Ioc.GetInstanceByName(v.(string))
		if _, ok := processor.(IProcessor); ok {
			p := processor.(IProcessor)
			multiProcess.RegisterProcessor(p.Name(), p.Processor())
			log.Println("---- Register thrift processor [", p.Name(), "] ----")
		}
	}

	server := thrift.NewTSimpleServer4(multiProcess, serverTransport, transportFactory, protocolFactory)
	log.Print("Starting server... on ", addr)
	err = server.Serve()
	if err != nil {
		log.Println(err)
		return
	}
	
}