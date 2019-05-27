package ilog

import (
	"bytes"
	"log"
	"os"
	"strings"
	"sync"
	"time"
	"fmt"
)

var (
	loggor ILog
	logfile string
	rotate bool
)

const (
	DEFAULT_LOG_FILE = "itea.log"
)

type ILog interface {
	Init()
	Done() bool
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}

/*****************************************************  日志文件  ****************************************************/

type Log struct {
	log *log.Logger
	wg sync.WaitGroup
	logfile *os.File
}

func (l *Log) Init() {
	if strings.EqualFold(logfile, "") {
		logfile = DEFAULT_LOG_FILE
	}
	if rotate {
		go l.rotateFile(logfile)
		logfile = l.rotateName(logfile)
	}
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0)
	if err != nil {
		panic("open log file error !")
	}
	l.logfile = file
	l.log = log.New(file, "", log.LstdFlags)
}

func (l *Log) rotateName(logfile string) string {
	f := strings.Split(logfile, ".")
	var s bytes.Buffer
	s.WriteString(strings.Join(f[0:len(f) - 1], "."))
	s.WriteString("-")
	s.WriteString(time.Now().Format("2006-01-02"))
	s.WriteString(".")
	s.WriteString(f[len(f) - 1])
	return s.String()
}

func (l *Log) rotateFile(logfile string) {
	filename := logfile
	for {
		now := time.Now()
		// 计算下一个零点
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
		name := l.rotateName(filename)
		for {
			file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0)
			if err == nil {
				l.logfile.Close()
				l.logfile = file
				l.log = log.New(file, "", log.LstdFlags)
				break
			}
		}
	}
}

func (l *Log) Done() bool {
	l.wg.Wait()
	return true
}

func (l *Log) Debug(v ...interface{}) {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.log.SetPrefix("[Debug] ")
		l.log.Println(v)
	}()
}

func (l *Log) Info(v ...interface{}) {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.log.SetPrefix("[INFO] ")
		l.log.Println(v)
	}()
}

func (l *Log) Error(v ...interface{}) {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.log.SetPrefix("[ERROR] ")
		l.log.Println(v)
	}()
}

func (l *Log) Fatal(v ...interface{}) {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		l.log.SetPrefix("[Fatal] ")
		l.log.Println(v)
	}()
}
/*****************************************************  日志文件  ****************************************************/

/******************************************************  控制台  *****************************************************/

type Console struct {

}

func (c *Console) Init() {

}

func (c *Console) Done() bool{
	return true
}

func (c *Console) Debug(v ...interface{}){
	fmt.Println("[DEBUG]", time.Now().Format("2006/01/02 15:04:05"),  v)
}

func (c *Console) Info(v ...interface{}){
	fmt.Println("[INFO]", time.Now().Format("2006/01/02 15:04:05"), v)
}

func (c *Console) Error(v ...interface{}){
	fmt.Println("[ERROR]", time.Now().Format("2006/01/02 15:04:05"), v)
}

func (c *Console) Fatal(v ...interface{}){
	fmt.Println("[FATAL]", time.Now().Format("2006/01/02 15:04:05"), v)
}

/******************************************************  控制台  *****************************************************/


func Init(logtype string, file string, ro bool) {
	logfile = file
	rotate = ro
	if strings.EqualFold(logtype, "file") {
		loggor = new(Log)
	}else {
		loggor = new(Console)
	}
	loggor.Init()
}

func Done() bool {
	return loggor.Done()
}

func Debug(v ...interface{}) {
	loggor.Debug(v...)
}

func Info(v ...interface{}) {
	loggor.Info(v...)
}

func Error(v ...interface{}) {
	loggor.Error(v...)
}

func Fatal(v ...interface{}) {
	loggor.Fatal(v...)
}