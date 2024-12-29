package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/babaYaga451/social/internal/store"
	"github.com/redis/go-redis/v9"
)

type UsersStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (s *UsersStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)
	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil

}

func (s *UsersStore) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.SetEx(ctx, cacheKey, json, UserExpTime).Err()
}

func (s *UsersStore) Delete(ctx context.Context, userID int64) {
	cacheKey := fmt.Sprintf("user-%v", userID)

	s.rdb.Del(ctx, cacheKey)
}
