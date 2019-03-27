// +build windows

package itea

import (
	"github.com/CalvinDjy/iteaGo/ilog"
	"os"
	"os/signal"
	"syscall"
)

func logProcessInfo() {
	ilog.Info("windows pid : ", os.Getpid())
}

func processSignal() {
	sigs = make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	for{
		msg := <-sigs
		switch msg {
		default:
			ilog.Info("[windows] default: ", msg)
			//case syscall.SIG:
			//reload
			//b.App.Reload(b.Conf)
			break
		case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
			ilog.Info("[windows]", msg)
			signal.Stop(sigs)
			s <- true
			return
		}
	}
}
