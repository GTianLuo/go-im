package dao

import (
	"fmt"
	"github.com/go-redis/redis"
	"testing"
)

func redisConn() *redis.Client {
	rc := redis.NewClient(&redis.Options{
		Addr:     "101.42.38.110:6379",
		Password: "2g0t0l374yyds",
		DB:       0,
	})
	_, err := rc.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("redis failed connected:%s", err.Error()))
	}
	return rc
}

func Test1(t *testing.T) {
	conn := redisConn()
	conn.HSet("test", "t1", "hello world")
}
