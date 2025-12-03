package redis

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", errors.Wrap(err, "could not get value")
	}

	return val, nil
}

func (r *Redis) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	_, err := r.client.Set(ctx, key, value, expiration).Result()
	if err != nil {
		return errors.Wrap(err, "could not set value")
	}

	return nil
}

// TODO: delete не понадобится, так как в Set мы указываем время жизни записи в кеше, после которой она сама удаляется
//func (r *Redis) Delete(ctx context.Context, key string) error {
//	_, err := r.client.Del(ctx, key).Result()
//	if err != nil {
//		return errors.Wrap(err, "could not delete value")
//	}
//	return nil
//}
