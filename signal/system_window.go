// +build windows

package signal

import (
	"github.com/CalvinDjy/iteaGo/ilog"
	"github.com/CalvinDjy/iteaGo/util/integer"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func LogProcessInfo() {
	pid := integer.Itos(os.Getpid())
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
	err = process.Kill()
	if err != nil {
		ilog.Error("process kill error : ", err)
	}
}

func ProcessSignal(sigs chan os.Signal, s chan bool) {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGKILL)
	for{
		msg := <-sigs
		switch msg {
		case syscall.SIGINT, syscall.SIGKILL:
			ilog.Info("[windows]", msg)
			signal.Stop(sigs)
			s <- true
			return
		case nil:
			break
		default:
			ilog.Info("[windows] default: ", msg)
			//case syscall.SIG:
			//reload
			//b.App.Reload(b.Conf)
			break
		}
	}
}
