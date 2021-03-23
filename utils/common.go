package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func CreateShutdownSignalReceiver() <-chan os.Signal {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	return signalChannel
}
