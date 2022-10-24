package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPortReader_Read(t *testing.T) {
	type testType struct {
		Name string `json:"name"`
	}

	type expected struct {
		key  string
		data testType
		err  error
	}

	tests := []struct {
		name     string
		input    string
		expected []expected
	}{
		{
			name:  "valid rows",
			input: `{"ABC": {"name": "abc"}, "DEF": {"name": "def"}, "GHI": {"name": "ghi"}}`,
			expected: []expected{{
				key:  "ABC",
				data: testType{Name: "abc"},
				err:  nil,
			}, {
				key:  "DEF",
				data: testType{Name: "def"},
				err:  nil,
			}, {
				key:  "GHI",
				data: testType{Name: "ghi"},
				err:  nil,
			}},
		},
		{
			name:     "empty file",
			input:    `{}`,
			expected: []expected{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			prs, err := New[testType](strings.NewReader(tc.input))
			require.NoError(t, err)

			var i int
			for {
				key, data, err := prs.Read()
				if err == io.EOF {
					break
				}

				require.Equal(t, tc.expected[i].key, key)
				require.Equal(t, tc.expected[i].data, *data)
				require.Equal(t, tc.expected[i].err, err)

				i++
			}
		})
	}
}
