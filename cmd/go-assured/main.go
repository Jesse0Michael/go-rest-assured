package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jesse0michael/go-rest-assured/pkg/assured"
)

// Preload is the expected format for preloading assured endpoints through the go rest assured application
type Preload struct {
	Calls []assured.Call `json:"calls"`
}

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	rootCtx := context.Background()

	errc := make(chan error)
	go func() {
		errc <- interrupt()
	}()

	port := flag.Int("port", 0, "a port to listen on. default automatically assigns a port.")
	logFile := flag.String("logger", "", "a file to send logs. default logs to STDOUT.")
	preload := flag.String("preload", "", "a file to parse preloaded calls from.")
	trackMade := flag.Bool("track", true, "a flag to enable the storing of calls made to the service.")

	flag.Parse()

	// If logger specified, set the assured client to write to the file
	if *logFile != "" {
		file, err := os.Create(*logFile)
		if err != nil {
			logger.Log("fatal", err.Error())
			os.Exit(1)
		}
		logger = kitlog.NewLogfmtLogger(file)
	}

	settings := assured.Settings{
		Logger:         logger,
		Port:           *port,
		TrackMadeCalls: *trackMade,
		HTTPClient:     *http.DefaultClient,
	}
	client := assured.NewClient(rootCtx, settings)

	// If preload file specified, parse the file and load all calls into the assured client
	if *preload != "" {
		file, err := os.Open(*preload)
		if err != nil {
			logger.Log("fatal", err.Error())
			os.Exit(1)
		}
		b, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Log("fatal", err.Error())
			os.Exit(1)
		}
		var preload Preload
		// TODO response won't unmarshal string to []byte
		if err := json.Unmarshal(b, &preload); err != nil {
			logger.Log("fatal", err.Error())
			os.Exit(1)
		}
		client.Given(preload.Calls...)
		if err != nil {
			logger.Log("fatal", err.Error())
			os.Exit(1)
		}
	}

	logger.Log("fatal", <-errc)
}

func interrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
