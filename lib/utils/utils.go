package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func TrapOSInterrupt(done chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	<-c
	close(done)
}
