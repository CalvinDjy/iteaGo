// +build windows

package itea

import (
	"log"
	"os"
	"syscall"
	"os/signal"
)

func logProcessInfo() {
	log.Println("windows pid : ", os.Getpid())
}

func processSignal() {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	for{
		msg := <-sigs
		switch msg {
		default:
			log.Println("[windows] default: ", msg)
			//case syscall.SIG:
			//reload
			//b.App.Reload(b.Conf)
			break
		case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
			log.Println("[windows]", msg)
			signal.Stop(sigs)
			s <- true
			return
		}
	}
}
