package dbConf

import (
	"fmt"
	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
)

func cache(addr, password string, db int) {
	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := rc.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("redis failed connected:%s", err.Error()))
	}
	fmt.Println()
	redisClient = rc
}

func NewRedisClient() *redis.Client {
	return redisClient
}
