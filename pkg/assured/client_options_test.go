package assured

import (
	"net/http"
	"reflect"
	"testing"
)

func TestClientOptions_applyOptions(t *testing.T) {
	tests := []struct {
		name    string
		options []ClientOption
		want    ClientOptions
	}{
		{
			name:    "default",
			options: nil,
			want: ClientOptions{
				httpClient: http.DefaultClient,
				baseURL:    "http://localhost",
			},
		},
		{
			name: "with http client",
			options: []ClientOption{
				WithClientHTTPClient(*http.DefaultClient),
			},
			want: ClientOptions{
				httpClient: http.DefaultClient,
				baseURL:    "http://localhost",
			},
		},
		{
			name: "with url",
			options: []ClientOption{
				WithClientBaseURL("https://example.com"),
			},
			want: ClientOptions{
				httpClient: http.DefaultClient,
				baseURL:    "https://example.com",
			},
		},
		{
			name: "combined options",
			options: []ClientOption{
				WithClientHTTPClient(*http.DefaultClient),
				WithClientBaseURL("http://localhost:1234"),
			},
			want: ClientOptions{
				httpClient: http.DefaultClient,
				baseURL:    "http://localhost:1234",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := DefaultClientOptions
			o.applyOptions(tt.options...)
			if !reflect.DeepEqual(o, tt.want) {
				t.Errorf("applyOptions() = %#v, want %#v", o, tt.want)
			}
		})
	}
}
