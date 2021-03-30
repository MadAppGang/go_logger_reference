package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"go_logger_reference/requestprocessor"
	"go_logger_reference/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	service := requestprocessor.BuildService("some config")

	err := service.StartListening()
	if err != nil {
		log.Fatal(err)
	}

	c1 := createClientFor("http://localhost:8080/some", "", "")
	c2 := createClientFor("http://localhost:8080/report/months", "Bob", "reporter")
	c3 := createClientFor("http://localhost:8080/report/days", "Linus", "developer")
	c4 := createClientFor("http://localhost:8080/admin/format-c", "Nick", "admin")
	c5 := createClientFor("http://localhost:8080/admin/list", "Alice", "hr")

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
		log.Printf("got signal %s", sig)
		cancel()
	}()

}

func createClientFor(url, user, role string) client {
	return client{url, user, role}
}

type client struct {
	url, user, role string
}

func (c client) CallAndLog() {
	logger := logrus.New()
	logger.AddHook(utils.LogDefaultField("who", "test_client"))
	logger.AddHook(utils.LogDefaultField("url", c.url))
	request, _ := http.NewRequest(http.MethodGet, c.url, nil)
	request.Header.Set("token", fmt.Sprintf(`{"username":"%s", "role":"%s"}`, c.user, c.role))

	response, err2 := http.DefaultClient.Do(request)
	if err2 != nil {
		log.Printf("Get %s failed: %v", c.url, err2)
		return
	}
	answer, err3 := io.ReadAll(response.Body)
	if err3 != nil {
		log.Printf("Read response from %s failed: %v", c.url, err3)
		return
	}
	log.Printf("%s: %d - %s", c.url, response.StatusCode, string(answer))
}
