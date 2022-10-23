package store

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"

	"github.com/rgallagher27/porter/internal/types"
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

// InsertPort inserts a given port in the store by its key.
func (s *Store) InsertPort(key string, p *types.Port) error {
	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("json marshal port: %w", err)
	}

	err = s.client.Set(key, b, 0).Err()
	if err != nil {
		return fmt.Errorf("set port: %w", err)
	}

	return nil
}
