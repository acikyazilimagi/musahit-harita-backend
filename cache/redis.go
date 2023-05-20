package cache

import (
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore() *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &RedisStore{client: client}
}

func (r *RedisStore) SetKey(key string, value []byte, ttl time.Duration) {
	err := r.client.Set(r.client.Context(), key, value, ttl).Err()
	if err != nil {
		log.Logger().Info("Unable to set key in redis" + key + err.Error())
	}
}

func (r *RedisStore) Get(key string) []byte {
	get := r.client.Get(r.client.Context(), key)

	resp, err := get.Bytes()
	if err != nil {
		log.Logger().Info("Unable to get key in redis" + key + err.Error())
	}

	return resp
}

func (r *RedisStore) Delete(key string) error {
	return r.client.Del(r.client.Context(), key).Err()
}

func (r *RedisStore) DeleteAll() error {
	return r.client.FlushDB(r.client.Context()).Err()
}

func (r *RedisStore) Close() {
	r.client.Close()
}
