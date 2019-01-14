package itea

import (
	"log"
	"os"
)

type ILog interface {
	Init()
	Debug(s string, v ...interface{})
	Info(s string, v ...interface{})
	Error(s string, v ...interface{})
	Fatal(s string, v ...interface{})
}

type Log struct {
	l *log.Logger
}

func(l *Log) Init() {
	fileName := "ll.log"
	logFile, err  := os.Create(fileName)
	if err != nil {
		panic("open file error !")
	}
	l.l = log.New(logFile, "", log.LstdFlags)
}

func(l *Log) Debug(s string, v ...interface{}) {

}

func(l *Log) Info(s string, v ...interface{}) {
	l.l.SetPrefix("[INFO]")
	l.l.Println(v)
}

func(l *Log) Error(s string, v ...interface{}) {

}

func(l *Log) Fatal(s string, v ...interface{}) {

}