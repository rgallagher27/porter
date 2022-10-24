package testhelper

import (
	"net"
	"testing"

	"github.com/go-redis/redis"
	"github.com/ory/dockertest"
)

// StartRedis will spin up a redis docker image for testing
func StartRedis(t *testing.T) (string, func()) {
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
