package redis

import (
	redis "gopkg.in/redis.v3"
)

const (
	identifier = "redis"
)

type RedisDoer struct {
	Password string
}

func (p *RedisDoer) Identifier() string {
	return identifier
}

func (r *RedisDoer) Do(addr string) (bool, []string) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: r.Password,
		DB:       0, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return false, nil
	}
	return true, nil
}
