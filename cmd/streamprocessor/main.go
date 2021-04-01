package main

import (
	"context"
	"go.uber.org/zap/zapcore"
	"go_logger_reference/log"
	"go_logger_reference/sample"
	"os"
	"os/signal"
	"syscall"

	"go_logger_reference/streamprocessor"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	loggerCfg := log.ZapConfig{
		Format:      log.JSON,
		MinLogLevel: zapcore.DebugLevel,
	}
	defaultLogger, err := log.NewZap(loggerCfg)
	if err != nil {
		panic(err)
	}

	service := streamprocessor.BuildService(loggerCfg)
	handler := sample.NewSampleHandler(defaultLogger)
	handler.Handle(context.WithValue(context.Background(), sample.RequestIDKey, "ID123"))

	go func() {
		ch := make(chan os.Signal, 0)
		signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGSTOP)
		sig := <-ch
		defaultLogger.Infow("got signal %s", sig)
		cancel()
	}()

	service.Run(ctx)
}
