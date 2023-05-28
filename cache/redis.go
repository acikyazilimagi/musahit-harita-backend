package cache

import (
	"fmt"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisStore struct {
	client *redis.Client
	conf   *Config
}

var RedisCache *RedisStore

func NewRedisStore() (*RedisStore, error) {
	if RedisCache != nil {
		return RedisCache, nil
	}

	conf := NewConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost,
		Password: conf.RedisPassword,
		DB:       0,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		log.Logger().Info(fmt.Sprintf("Unable to connect to redis: %s", err.Error()))
		return nil, err
	} else {
		log.Logger().Info("Connected to redis")
	}

	RedisCache = &RedisStore{client: client, conf: conf}

	return RedisCache, nil
}

func (r *RedisStore) Set(key string, value interface{}, ttl time.Duration) {
	err := r.client.Set(r.client.Context(), key, value, ttl).Err()
	if err != nil {
		log.Logger().Info("Unable to set key in redis" + key + err.Error())
	}
}

func (r *RedisStore) Get(key string) interface{} {
	get := r.client.Get(r.client.Context(), key)

	resp, err := get.Result()
	if err != nil {
		log.Logger().Info("Unable to get key in redis" + key + err.Error())
	}

	return resp
}

func (r *RedisStore) SetCacheResponse(key string, value []byte, ttl time.Duration) {
	err := r.client.Set(r.client.Context(), key, value, ttl).Err()
	if err != nil {
		log.Logger().Info("Unable to set key in redis" + key + err.Error())
		return
	}
}

func (r *RedisStore) GetCacheResponse(key string) []byte {
	get := r.client.Get(r.client.Context(), key)
	if get == nil {
		log.Logger().Info("Unable to get key in redis" + key)
		return nil
	}

	if get.Err() != nil {
		log.Logger().Info("Unable to get key in redis" + key + get.Err().Error())
		return nil
	}

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
