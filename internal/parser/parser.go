package parser

import (
	"encoding/json"
	"fmt"
	"io"
)

// Parser is a generic json parser for processing documents of the format:
//
//	{
//	  "key": {
//	     ... object [T]
//	   },
//	...
//	}
type Parser[T any] struct {
	dec *json.Decoder
}

// New initialises a new parser of the type T
func New[T any](rdc io.Reader) (*Parser[T], error) {
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

	return &Parser[T]{
		dec: dec,
	}, nil
}

// Read reads the next key and data T from the underlying parser.
// When no more records remain, an EOF error will be returned
func (p *Parser[T]) Read() (string, *T, error) {
	if ok := p.dec.More(); !ok {
		return "", nil, io.EOF
	}

	t, err := p.dec.Token()
	if err != nil {
		return "", nil, err
	}
	key := t.(string)

	var data T
	if err := p.dec.Decode(&data); err != nil {
		return key, nil, err
	}

	return key, &data, err
}
