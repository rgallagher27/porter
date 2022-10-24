package port

import (
	"context"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Run(t *testing.T) {
	tests := []struct {
		name           string
		ignoreReadErrs bool
		testFileName   string
		insertErrors   []error
		expectErr      require.ErrorAssertionFunc
		expectPorts    []string
	}{
		{
			name:           "valid file",
			ignoreReadErrs: true,
			testFileName:   "valid.json",
			insertErrors:   []error{nil, nil, nil, nil},
			expectErr:      require.NoError,
			expectPorts:    []string{"prt:AEAJM", "prt:AEAUH", "prt:AEDXB", "prt:AEFJR"},
		},
		{
			name:           "invalid file format",
			ignoreReadErrs: true,
			testFileName:   "invalid_format.json",
			insertErrors:   []error{},
			expectErr: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "new parser")
			},
			expectPorts: []string{},
		},
		{
			name:           "continues on invalid row error",
			ignoreReadErrs: true,
			testFileName:   "invalid_row.json",
			insertErrors:   []error{nil, nil, nil, nil},
			expectErr:      require.NoError,
			expectPorts:    []string{"prt:AEAJM", "prt:AEDXB", "prt:AEFJR"},
		},
		{
			name:           "returns on invalid row error",
			ignoreReadErrs: false,
			testFileName:   "invalid_row.json",
			insertErrors:   []error{nil},
			expectErr:      require.Error,
			expectPorts:    []string{"prt:AEAJM"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var (
				insertCount int
				written     []string
			)
			mockStore := storeMock{
				InsertFunc: func(key string, _ any) error {
					defer func() { insertCount++ }()
					written = append(written, key)

					return tc.insertErrors[insertCount]
				},
			}

			service := New(&mockStore, tc.ignoreReadErrs)

			testFile := readFile(t, path.Join("testdata", tc.testFileName))
			inputRdr := strings.NewReader(testFile)

			err := service.Run(context.Background(), inputRdr)
			tc.expectErr(t, err)

			assert.Len(t, written, len(tc.expectPorts))
			for _, v := range tc.expectPorts {
				assert.Contains(t, written, v)
			}

		})
	}
}

func readFile(t *testing.T, name string) string {
	t.Helper()

	b, err := os.ReadFile(name)
	require.NoError(t, err)

	return string(b)
}
