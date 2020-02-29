package sessions

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisSessionStore struct {
	Pool *redis.Pool
}

func dial(network, address string) (redis.Conn, error) {
	return redis.Dial(network, address)
}

func (redisSessionStore *RedisSessionStore) New(network, address string) error {
	redisSessionStore.Pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return dial(network, address)
		},
	}

	conn := redisSessionStore.Pool.Get()
	defer conn.Close()

	ping, err := redis.String(conn.Do("PING"))
	if err != nil {
		return err
	}

	if ping != "PONG" {
		return errors.New("unable to initiate connection")
	}

	return nil
}

func (redisSessionStore *RedisSessionStore) Save(key string, value map[string]interface{}, expiry time.Duration) (err error) {
	conn := redisSessionStore.Pool.Get()
	defer conn.Close()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if expiry == 0 {
		expiry = 3600
	}

	_, err = conn.Do("SETEX", key, int(expiry), data)

	return err
}

func (redisSessionStore *RedisSessionStore) Get(key string) (data map[string]interface{}, err error) {
	conn := redisSessionStore.Pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(s), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (redisSessionStore *RedisSessionStore) Delete(key string) error {
	conn := redisSessionStore.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (redisSessionStore *RedisSessionStore) Reset() error {
	conn := redisSessionStore.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHDB")
	return err
}
