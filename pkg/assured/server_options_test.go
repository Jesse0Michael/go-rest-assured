package assured

import (
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestServerOptions_applyOptions(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	tests := []struct {
		name   string
		option ServerOption
		want   ServerOptions
	}{
		{
			name:   "with http client",
			option: WithHTTPClient(*http.DefaultClient),
			want: ServerOptions{
				httpClient: http.DefaultClient,
			},
		},
		{
			name:   "with host",
			option: WithHost("rest-assured"),
			want: ServerOptions{
				host: "rest-assured",
			},
		},
		{
			name:   "with port",
			option: WithPort(8889),
			want: ServerOptions{
				Port: 8889,
			},
		},
		{
			name:   "with tls",
			option: WithTLS("cert", "key"),
			want: ServerOptions{
				tlsCertFile: "cert",
				tlsKeyFile:  "key",
			},
		},
		{
			name:   "with track made calls",
			option: WithCallTracking(true),
			want: ServerOptions{
				trackMadeCalls: true,
			},
		},
		{
			name:   "with logger",
			option: WithLogger(logger),
			want: ServerOptions{
				logger: logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := ServerOptions{}
			o.applyOptions(tt.option)
			if !reflect.DeepEqual(o, tt.want) {
				t.Errorf("applyOptions() = %#v, want %#v", o, tt.want)
			}
		})
	}
}
