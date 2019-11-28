// +build !windows

package signal

import (
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/CalvinDjy/iteaGo/util/integer"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"strings"
	"io/ioutil"
)

func LogProcessInfo() {
	pid := integer.Itos(os.Getpid())
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

func RemovePid() {
	os.Remove("pid")
}

func StopProcess() {
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

func ProcessSignal(sigs chan os.Signal, s chan bool) {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL,syscall.SIGUSR1, syscall.SIGUSR2, os.Interrupt)
	for{
		msg := <-sigs
		switch msg {
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
		case nil:
			break
		default:
			ilog.Info("[linux] default: ", msg)
			break
		}
	}
}
