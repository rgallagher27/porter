package store

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rgallagher27/porter/internal/testhelper"
)

type testType struct {
	Name string `json:"name"`
	City string `json:"city"`
}

func TestStore_InsertPort(t *testing.T) {
	addr, destroyFunc := testhelper.StartRedis(t)
	t.Cleanup(destroyFunc)

	str, err := New(Config{
		Address:  addr,
		Password: "",
		DB:       0,
	})
	require.NoError(t, err)

	testKey := "test_key"
	testValue := &testType{
		Name: "a test port",
		City: "Edinburgh",
	}

	err = str.Insert(testKey, testValue)
	require.NoError(t, err)

	gotPort, err := str.client.Get(testKey).Result()
	require.NoError(t, err)

	assert.Equal(t, marshalValue(t, testValue), gotPort)
}

// marshalValue is a helper func to marshal a testType struct to a string for test comparison
func marshalValue(t *testing.T, p *testType) string {
	t.Helper()

	b, err := json.Marshal(p)
	require.NoError(t, err)

	return string(b)
}
