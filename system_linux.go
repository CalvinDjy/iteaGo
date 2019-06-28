// +build !windows

package itea

import (
	"github.com/CalvinDjy/iteaGo/ilog"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"strings"
	"io/ioutil"
)

func logProcessInfo() {
	pid := strconv.Itoa(os.Getpid())
	ilog.Info("linux pid : ", pid)
	file, err := os.OpenFile("pid", os.O_CREATE|os.O_WRONLY,0)
	if err != nil {
		panic("open pid file error !")
	}
	file.WriteString(pid)
	file.Close()
}

func getPid() string {
	r, err := ioutil.ReadFile("pid")
	if err != nil {
		return ""
	}
	return string(r)
}

func removePid() {
	os.Remove("pid")
}

func stopProcess() {
	pid := getPid()
	if strings.EqualFold(pid, "") {
		return
	}
	iPid, err := strconv.Atoi(pid)
	if err != nil {
		ilog.Error("get pid error : ", err)
		return
	}
	process, err := os.FindProcess(iPid)
	if err != nil {
		ilog.Error("find process error : ", err)
		return
	}
	err = process.Signal(syscall.SIGINT)
	if err != nil {
		ilog.Error("process kill : ", err)
	}
}

func processSignal() {
	sigs = make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL,syscall.SIGUSR1, syscall.SIGUSR2, os.Interrupt)
	for{
		msg := <-sigs
		switch msg {
		default:
			ilog.Info("[linux] default: ", msg)
			break
		case syscall.SIGUSR1:
			//reload
			ilog.Info("[linux] SIGUSR1: ", msg)
			break
		case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
			//logger.Info("application stoping, signal[%v]", msg)
			//b.App.Stop()
			ilog.Info("[linux]", msg)
			signal.Stop(sigs)
			s <- true
			return
		}
	}
}
