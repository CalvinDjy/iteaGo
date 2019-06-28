// +build windows

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
	ilog.Info("windows pid : ", pid)
	file, err := os.OpenFile("pid", os.O_CREATE|os.O_TRUNC,0)
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
	err = process.Kill()
	if err != nil {
		ilog.Error("process kill error : ", err)
	}
}

func processSignal() {
	sigs = make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL)
	for{
		msg := <-sigs
		switch msg {
		default:
			ilog.Info("[windows] default: ", msg)
			//case syscall.SIG:
			//reload
			//b.App.Reload(b.Conf)
			break
		case syscall.SIGINT, syscall.SIGKILL:
			ilog.Info("[windows]", msg)
			signal.Stop(sigs)
			s <- true
			return
		}
	}
}
