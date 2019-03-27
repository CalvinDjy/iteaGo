// +build !windows

package itea

import (
	"github.com/CalvinDjy/iteaGo/ilog"
	"os"
	"os/signal"
	"syscall"
)

func logProcessInfo() {
	ilog.Info("linux pid : ", os.Getpid())
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
			log.Println("[linux] SIGUSR1: ", msg)
			break
		case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
			//logger.Info("application stoping, signal[%v]", msg)
			//b.App.Stop()
			log.Println("[linux]", msg)
			signal.Stop(sigs)
			s <- true
			return
		}
	}
}
