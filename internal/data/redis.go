package data

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewKVStore(addr, pwd string) (*Store, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0,
	})

	if _, err := c.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &Store{
		client: c,
	}, nil
}

func (s *Store) Set(ctx context.Context, k string, v []byte) error {
	c := s.client.Set(ctx, k, v, 60*time.Second)
	return c.Err()
}

func (s *Store) Get(ctx context.Context, k string) ([]byte, error) {
	c := s.client.Get(ctx, k)
	return c.Bytes()
}
