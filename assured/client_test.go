package assured

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	httpClient := &http.Client{}
	client := NewClient(kitlog.NewLogfmtLogger(ioutil.Discard), nil).Run()
	url := client.URL()

	require.NoError(t, client.Given(call1))
	require.NoError(t, client.Given(call2))

	req, err := http.NewRequest(http.MethodGet, url+"test/assured", bytes.NewReader([]byte(`{"calling":"you"}`)))
	require.NoError(t, err)

	hit1, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, hit1.StatusCode)
	body, err := ioutil.ReadAll(hit1.Body)
	require.NoError(t, err)
	require.Equal(t, []byte(`{"assured": true}`), body)

	req, err = http.NewRequest(http.MethodGet, url+"test/assured", bytes.NewReader([]byte(`{"calling":"again"}`)))
	require.NoError(t, err)

	hit2, err := httpClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusConflict, hit2.StatusCode)
	body, err = ioutil.ReadAll(hit2.Body)
	require.NoError(t, err)
	require.Equal(t, []byte("error"), body)
}
