package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jesse0michael/go-rest-assured/assured"
	"github.com/phayes/freeport"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	rootCtx := context.Background()

	errc := make(chan error)
	go func() {
		errc <- interrupt()
	}()

	port := freeport.GetPort()
	assured.StartApplicationHTTPListener(port, logger, rootCtx, errc)

	logger.Log("fatal", <-errc)
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
