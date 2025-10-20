package assured

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientInvalidPort(t *testing.T) {
	client := NewServer(WithPort(-1))

	require.Error(t, client.Serve())
}
