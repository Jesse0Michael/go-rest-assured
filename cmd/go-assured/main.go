package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jesse0michael/go-rest-assured/v3/pkg/assured"
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
	preload := flag.String("preload", "", "a file to parse preloaded calls from.")
	trackMade := flag.Bool("track", true, "a flag to enable the storing of calls made to the service.")
	host := flag.String("host", "localhost", "a host to use in the client's url.")
	tlsCert := flag.String("tlsCert", "", "location of tls cert for serving https traffic. tlsKey also required, if specified.")
	tlsKey := flag.String("tlsKey", "", "location of tls key for serving https traffic. tlsCert also required, if specified")

	flag.Parse()

	client := assured.NewClient(assured.WithContext(rootCtx), assured.WithPort(*port),
		assured.WithCallTracking(*trackMade), assured.WithLogger(logger), assured.WithHost(*host), assured.WithTLS(*tlsCert, *tlsKey))

	// If preload file specified, parse the file and load all calls into the assured client
	if *preload != "" {
		b, err := os.ReadFile(*preload)
		if err != nil {
			_ = logger.Log("fatal", err.Error())
			os.Exit(1)
		}
		var preload Preload
		// TODO response won't unmarshal string to []byte
		if err := json.Unmarshal(b, &preload); err != nil {
			_ = logger.Log("fatal", err.Error())
			os.Exit(1)
		}
		if err = client.Given(preload.Calls...); err != nil {
			_ = logger.Log("fatal", err.Error())
			os.Exit(1)
		}
	}

	_ = logger.Log("fatal", <-errc)
}

func interrupt() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return fmt.Errorf("%s", <-c)
}
