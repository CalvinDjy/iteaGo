package thrift

import "github.com/apache/thrift/lib/go/thrift"

type IProcessor interface {
	Name() string
	Processor() thrift.TProcessor
}