package parser

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPortReader_Read(t *testing.T) {
	rdr := openFile(t, path.Join("testdata", "valid_ports.json"))

	type testType struct {
		Name string `json:"name"`
	}

	portParser, err := New[testType](rdr)
	require.NoError(t, err)

	key, port, err := portParser.Read()
	require.NoError(t, err)

	assert.Equal(t, "AEAJM", key)
	assert.NotNil(t, port)
	assert.Equal(t, "Ajman", port.Name)

	key, port, err = portParser.Read()
	require.NoError(t, err)

	assert.Equal(t, "AEAUH", key)
	assert.NotNil(t, port)
	assert.Equal(t, "Abu Dhabi", port.Name)

	key, port, err = portParser.Read()
	require.NoError(t, err)

	assert.Equal(t, "AEDXB", key)
	assert.NotNil(t, port)
	assert.Equal(t, "Dubai", port.Name)

	_, _, err = portParser.Read()
	require.Error(t, err)
	require.ErrorIs(t, err, io.EOF)

}

func openFile(t *testing.T, name string) io.Reader {
	t.Helper()

	f, err := os.Open(name)
	require.NoError(t, err)

	return f
}
