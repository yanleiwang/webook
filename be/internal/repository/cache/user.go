package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/be/internal/domain"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewRedisUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (r *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := r.key(id)
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}

	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

func (r *RedisUserCache) Set(ctx context.Context, user domain.User) error {
	key := r.key(user.Id)
	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, bytes, r.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
