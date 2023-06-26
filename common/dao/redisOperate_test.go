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
	account := "2985496686"
	conn := redisConn()
	result, err := conn.Exists(GateWayConnsStatus + account).Result()
	if err != nil {
		panic(err)
	}
	if result == 0 {
		// 用户不在线
		fmt.Println("用户不在线")
	}

	cmd := conn.Get(GateWayConnsStatus + account)
	if err := cmd.Err(); err != nil {
		panic(err)
	}
	rs := cmd.String()
	fmt.Printf("用户在线: %s", rs)
}
