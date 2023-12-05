package assured

import (
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func Test_applyOptions(t *testing.T) {
	tests := []struct {
		name   string
		option Option
		want   Options
	}{
		{
			name:   "with http client",
			option: WithHTTPClient(*http.DefaultClient),
			want: Options{
				httpClient: http.DefaultClient,
			},
		},
		{
			name:   "with host",
			option: WithHost("rest-assured"),
			want: Options{
				host: "rest-assured",
			},
		},
		{
			name:   "with port",
			option: WithPort(8889),
			want: Options{
				Port: 8889,
			},
		},
		{
			name:   "with track",
			option: WithCallTracking(true),
			want: Options{
				trackMadeCalls: true,
			},
		},
		{
			name:   "with logger",
			option: WithLogger(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))),
			want: Options{
				logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{})),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{}
			o.applyOptions(tt.option)
			if !reflect.DeepEqual(o, tt.want) {
				t.Errorf("applyOptions() = %v, want %v", o, tt.want)
			}
		})
	}
}
