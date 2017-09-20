package assured

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallID(t *testing.T) {
	call := Call{
		Path:   "/given/test/assured",
		Method: "GET",
	}

	require.Equal(t, "GET:/given/test/assured", call.ID())
}

func TestCallIDNil(t *testing.T) {
	call := Call{}

	require.Equal(t, ":", call.ID())
}

func TestCallString(t *testing.T) {
	call := Call{
		Response: []byte("GO assured is one way to GO"),
	}

	require.Equal(t, "GO assured is one way to GO", call.String())
}

func TestCallStringNil(t *testing.T) {
	call := Call{}

	require.Equal(t, "", call.String())
}
