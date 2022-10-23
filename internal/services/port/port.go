package port

import (
	"context"
	"fmt"
	"io"

	"github.com/rgallagher27/porter/internal/parser"
)

const (
	portKeyPrefix = "prt"
)

type store interface {
	Insert(key string, v any) error
}

type Service struct {
	store        store
	ignoreErrors bool
}

func New(s store, ignoreErrors bool) *Service {
	return &Service{
		store:        s,
		ignoreErrors: ignoreErrors,
	}
}

// Run parses the data from the reader row by row, validating each port record
// and inserting into the store
func (s *Service) Run(ctx context.Context, reader io.Reader) error {
	prs, err := parser.New[Port](reader)
	if err != nil {
		return fmt.Errorf("new parser: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// continue with read loop
		}

		key, port, err := prs.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("parse read: %w", err)
		}

		// Validate the port type
		if err := port.Validate(); err != nil {
			if s.ignoreErrors {
				continue
			}

			return fmt.Errorf("port validation: %w", err)
		}

		// Prefix key with a port specific identifier
		key = fmt.Sprintf("%s:%s", portKeyPrefix, key)
		if err := s.store.Insert(key, port); err != nil {
			return fmt.Errorf("store insert: %w", err)
		}

		// Return port object back to the parser pool
		prs.Return(port)
	}

	return nil
}
