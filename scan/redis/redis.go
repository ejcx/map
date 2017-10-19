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
	fmt.Println("Redis scanning")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0, // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		return false, nil
	}
	fmt.Println(pong)
	return true, nil
}
