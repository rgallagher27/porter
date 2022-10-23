package store

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/go-redis/redis"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rgallagher27/porter/internal/services/port"
)

func TestStore_InsertPort(t *testing.T) {
	addr, destroyFunc := startRedis(t)
	t.Cleanup(destroyFunc)

	str, err := New(Config{
		Address:  addr,
		Password: "",
		DB:       0,
	})
	require.NoError(t, err)

	testKey := "test_key"
	testPort := &port.Port{
		Name:     "a test port",
		City:     "Edinburgh",
		Province: "Scotland",
		Country:  "United Kingdom",
	}

	err = str.Insert(testKey, testPort)
	require.NoError(t, err)

	gotPort, err := str.client.Get(testKey).Result()
	require.NoError(t, err)

	assert.Equal(t, marshalPort(t, testPort), gotPort)
}

// marshalPort is a helper func to marshal a port struct to a string for test comaprison
func marshalPort(t *testing.T, p *port.Port) string {
	t.Helper()

	b, err := json.Marshal(p)
	require.NoError(t, err)

	return string(b)
}

// startRedis will spin up a redis docker image for testing
func startRedis(t *testing.T) (string, func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("starting dockertest: %+v", err)
	}

	resource, err := pool.Run("redis", "7-alpine", nil)
	if err != nil {
		t.Fatalf("starting redis: %+v", err)
	}

	addr := net.JoinHostPort("localhost", resource.GetPort("6379/tcp"))

	// wait for the container to be ready
	err = pool.Retry(func() error {
		var e error
		client := redis.NewClient(&redis.Options{Addr: addr})
		defer client.Close()

		_, e = client.Ping().Result()
		return e
	})

	if err != nil {
		t.Fatalf("ping Redis: %+v", err)
	}

	destroyFunc := func() {
		pool.Purge(resource)
	}

	return addr, destroyFunc
}
