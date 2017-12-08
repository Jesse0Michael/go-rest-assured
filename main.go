package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jesse0michael/go-rest-assured/assured"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	rootCtx := context.Background()

	errc := make(chan error)
	go func() {
		errc <- interrupt()
	}()

	port := flag.Int("port", 0, "a port to listen on. default automatically assigns a port.")
	logFile := flag.String("logger", "", "a file to send logs. default logs to STDOUT.")

	flag.Parse()

	if *logFile != "" {
		file, err := os.Create(*logFile)
		if err != nil {
			logger.Log("fatal", err.Error())
			os.Exit(1)
		}
		logger = kitlog.NewLogfmtLogger(file)
	}

	assured.NewClient(rootCtx, *port, &logger)

	logger.Log("fatal", <-errc)
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
