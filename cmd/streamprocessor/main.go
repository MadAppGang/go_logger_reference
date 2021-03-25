package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go_logger_reference/streamprocessor"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	service := streamprocessor.BuildService("some config")

	go func() {
		ch := make(chan os.Signal, 0)
		signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGSTOP)
		sig := <-ch
		log.Printf("got signal %s", sig)
		cancel()
	}()

	service.Run(ctx)
}
