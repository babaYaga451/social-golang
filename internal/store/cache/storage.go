package cache

import (
	"context"

	"github.com/babaYaga451/social/internal/store"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	User interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
		Delete(context.Context, int64)
	}
}

func NewRedisStore(rdb *redis.Client) Storage {
	return Storage{
		User: &UsersStore{rdb: rdb},
	}
}
