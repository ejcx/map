package redis

import (
	"fmt"

	redis "gopkg.in/redis.v3"
)

var (
	Doer = &RedisDoer{}
)

type RedisDoer struct{}

func (p *RedisDoer) Do(addr string) bool {
	fmt.Println("Redis scanning")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0, // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(pong)
	return true
}
