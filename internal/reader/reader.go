package reader

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/rgallagher27/porter/internal/types"
)

type PortReader struct {
	dec *json.Decoder
}

func NewPortReader(rdc io.Reader) (*PortReader, error) {
	dec := json.NewDecoder(rdc)

	// Get first token
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}

	// Expect to be at the start of a JSON object
	if t != json.Delim('{') {
		return nil, fmt.Errorf("expected { but got %v", t)
	}

	return &PortReader{
		dec: dec,
	}, nil
}

// Read reads the next key and port from the underlying reader.
// When no more records remain, an EOF error will be returned
func (p *PortReader) Read() (string, *types.Port, error) {
	if ok := p.dec.More(); !ok {
		return "", nil, io.EOF
	}

	t, err := p.dec.Token()
	if err != nil {
		return "", nil, err
	}
	key := t.(string)

	var port types.Port
	if err := p.dec.Decode(&port); err != nil {
		return key, nil, err
	}

	return key, &port, err
}
