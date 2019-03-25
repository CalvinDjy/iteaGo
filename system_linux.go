// +build !windows

package itea

import (
	"log"
	"os"
	"syscall"
	"os/signal"
)

func logProcessInfo() {
	log.Println("linux pid : ", os.Getpid())
}

func processSignal() {
	sigs = make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL,syscall.SIGUSR1, syscall.SIGUSR2, os.Interrupt)
	for{
		msg := <-sigs
		switch msg {
		default:
			log.Println("[linux] default: ", msg)
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
