package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jesse0michael/go-rest-assured/v5/pkg/assured"
)

// Preload is the expected format for preloading assured endpoints through the go rest assured application
type Preload struct {
	Calls []assured.Call `json:"calls"`
}

func main() {
	ctx, cancel := context.WithCancelCause(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		cancel(fmt.Errorf("%s", <-sig))
	}()

	port := flag.Int("port", 0, "a port to listen on. default automatically assigns a port.")
	preload := flag.String("preload", "", "a file to parse preloaded calls from.")
	trackMade := flag.Bool("track", true, "a flag to enable the storing of calls made to the service.")
	host := flag.String("host", "localhost", "a host to use in the client's url.")
	tlsCert := flag.String("tlsCert", "", "location of tls cert for serving https traffic. tlsKey also required, if specified.")
	tlsKey := flag.String("tlsKey", "", "location of tls key for serving https traffic. tlsCert also required, if specified")

	flag.Parse()

	a := assured.NewAssured(
		assured.WithPort(*port),
		assured.WithCallTracking(*trackMade),
		assured.WithHost(*host),
		assured.WithTLS(*tlsCert, *tlsKey))

	go func() {
		slog.InfoContext(ctx, "starting assured server", "port", a.Port)
		if err := a.Serve(ctx); err != nil {
			slog.InfoContext(ctx, "assured server stopped serving", "error", err)
		}
	}()

	// If preload file specified, parse the file and load all calls into the assured client
	if *preload != "" {
		b, err := os.ReadFile(*preload)
		if err != nil {
			slog.InfoContext(ctx, "failed to read preload file", "error", err)
			cancel(err)
		}
		var preload Preload
		// TODO response won't unmarshal string to []byte
		if err := json.Unmarshal(b, &preload); err != nil {
			slog.InfoContext(ctx, "failed to unmarshal preload file", "error", err)
			cancel(err)
		}
		if err = a.Given(ctx, preload.Calls...); err != nil {
			slog.InfoContext(ctx, "failed to set given preload file calls", "error", err)
			cancel(err)
		}
	}

	<-ctx.Done()
	if err := a.Close(); err != nil {
		slog.Info("failed to close assured", "error", err)
	}
	slog.Info("exiting assured")
}
