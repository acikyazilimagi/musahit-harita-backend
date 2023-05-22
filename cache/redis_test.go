package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func TestRedisStore(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{
		Addr: server.Addr(),
	})

	store := RedisStore{
		Client: client,
	}
	
	key := "test_key"
	value := []byte("test_value")
	ttl := 10 * time.Second

	store.SetKey(key, value, ttl)

	resp := client.Get(client.Context(), key)
	assert.NoError(t, resp.Err())
	assert.Equal(t, string(value), resp.Val())

	respValue := store.Get(key)
	assert.Equal(t, value, respValue)

	err := store.Delete(key)
	assert.NoError(t, err)

	resp = client.Get(client.Context(), key)
	assert.EqualError(t, redis.Nil, resp.Err().Error())

	err = store.DeleteAll()
	assert.NoError(t, err)


	keys := client.Keys(client.Context(), "*")
	assert.NoError(t, keys.Err())
	assert.Empty(t, keys.Val())

	store.Close()
	

	_, err = client.Ping(client.Context()).Result()
	assert.EqualError(t, err, "redis: client is closed")
}
