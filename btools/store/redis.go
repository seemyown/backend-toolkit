package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Store interface {
	Set(ctx context.Context, key string, value any, expIn time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Scan(ctx context.Context, pattern string, count int64) ([]string, error)
	Delete(ctx context.Context, keys ...string) error
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (s *RedisStore) Set(ctx context.Context, key string, value any, expIn time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal err: %w", err)
	}
	return s.client.Set(ctx, key, data, expIn).Err()
}

func (s *RedisStore) Get(ctx context.Context, key string, dest any) error {
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return err
		}
		return fmt.Errorf("redis get err: %w", err)
	}
	return json.Unmarshal(data, dest)
}

func (s *RedisStore) Keys(ctx context.Context, pattern string) ([]string, error) {
	return s.client.Keys(ctx, pattern).Result()
}

func (s *RedisStore) Scan(ctx context.Context, pattern string, count int64) ([]string, error) {
	var (
		cursor uint64
		keys   []string
	)
	for {
		var newKeys []string
		var err error
		newKeys, cursor, err = s.client.Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, newKeys...)
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

func (s *RedisStore) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return s.client.Del(ctx, keys...).Err()
}
