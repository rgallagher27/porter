package store

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type Store struct {
	client *redis.Client
}

type Config struct {
	Address  string
	Password string
	DB       int
}

// New initialises a new Store
func New(cfg Config) (*Store, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Store{
		client: client,
	}, nil
}

// Insert inserts a given v in the store by key.
func (s *Store) Insert(key string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	err = s.client.Set(key, b, 0).Err()
	if err != nil {
		return fmt.Errorf("set: %w", err)
	}

	return nil
}
