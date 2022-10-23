package reader

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPortReader_Read(t *testing.T) {
	rdr := openFile(t, path.Join("testdata", "valid_ports.json"))

	portReader, err := NewPortReader(rdr)
	require.NoError(t, err)

	for {
		key, port, err := portReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
		}

		t.Log(key, port)
	}

}

func openFile(t *testing.T, name string) io.Reader {
	t.Helper()

	f, err := os.Open(name)
	require.NoError(t, err)

	return f
}
