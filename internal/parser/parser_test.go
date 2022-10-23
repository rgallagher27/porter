package parser

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rgallagher27/porter/internal/services/port"
)

func TestPortReader_Read(t *testing.T) {
	rdr := openFile(t, path.Join("testdata", "valid_ports.json"))

	portParser, err := New[port.Port](rdr)
	require.NoError(t, err)

	for {
		key, port, err := portParser.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
		}

		t.Log(key, port)
		portParser.Return(port)
	}

}

func openFile(t *testing.T, name string) io.Reader {
	t.Helper()

	f, err := os.Open(name)
	require.NoError(t, err)

	return f
}
