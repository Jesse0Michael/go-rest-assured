package assured

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	kitlog "github.com/go-kit/kit/log"
)

func Test_applyOptions(t *testing.T) {
	log := kitlog.NewNopLogger()
	ctx := context.TODO()
	tests := []struct {
		name   string
		option Option
		want   Options
	}{
		{
			name:   "with context",
			option: WithContext(ctx),
			want: Options{
				ctx: ctx,
			},
		},
		{
			name:   "with logger",
			option: WithLogger(log),
			want: Options{
				logger: log,
			},
		},
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
