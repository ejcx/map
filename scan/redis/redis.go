package redis

import (
	"fmt"

	redis "gopkg.in/redis.v3"
)

const (
	identifier = "redis"
)

var (
	Doer = &RedisDoer{}
)

type RedisDoer struct{}

func (p *RedisDoer) Identifier() string {
	return identifier
}
func (p *RedisDoer) Do(addr string) (bool, []string) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0, // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return false, nil
	}
	fmt.Printf("%s : Redis")
	return true, nil
}
