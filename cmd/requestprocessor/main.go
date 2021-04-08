package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_logger_reference/log"
	"go_logger_reference/requestprocessor"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	loggerCfg := log.ZapConfig{
		Format:      log.JSON,
		MinLogLevel: zapcore.DebugLevel,
	}
	defaultLogger, err := log.NewZap(loggerCfg)

	service := requestprocessor.BuildService("some config", defaultLogger)

	if err = service.StartListening(); err != nil {
		defaultLogger.Fatal(err)
	}

	c1 := createClientFor("http://localhost:8080/some", "", "", defaultLogger)
	c2 := createClientFor("http://localhost:8080/report/months", "Bob", "reporter", defaultLogger)
	c3 := createClientFor("http://localhost:8080/report/days", "Linus", "developer", defaultLogger)
	c4 := createClientFor("http://localhost:8080/admin/format-c", "Nick", "admin", defaultLogger)
	c5 := createClientFor("http://localhost:8080/admin/list", "Alice", "hr", defaultLogger)

	ticker := time.NewTicker(500 * time.Millisecond)
	timer := time.NewTimer(5 * time.Second)

	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			return
		case <-ticker.C:
			go c1.CallAndLog()
			go c2.CallAndLog()
			go c3.CallAndLog()
			go c4.CallAndLog()
			go c5.CallAndLog()
		}
	}

	go func() {
		ch := make(chan os.Signal, 0)
		signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGSTOP)
		sig := <-ch
		defaultLogger.Infow("got signal", "sig", sig)
		cancel()
	}()

}

func createClientFor(url, user, role string, logger *zap.SugaredLogger) client {
	customLogger := logger.With("who", "test_client").With("url", url)
	return client{url, user, role, customLogger}
}

type client struct {
	url, user, role string
	logger          *zap.SugaredLogger
}

func (c client) CallAndLog() {
	request, _ := http.NewRequest(http.MethodGet, c.url, nil)
	request.Header.Set("token", fmt.Sprintf(`{"username":"%s", "role":"%s"}`, c.user, c.role))

	response, err2 := http.DefaultClient.Do(request)
	if err2 != nil {
		c.logger.Infow("Get url failed", "err", err2)
		return
	}
	answer, err3 := io.ReadAll(response.Body)
	if err3 != nil {
		c.logger.Infow("Read response failed", "err", err2)
		return
	}

	c.logger.Infow("received response", "statusCode", response.StatusCode, "answer", answer)
}
