package utils

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

//RedisPool  .
var RedisPool *redis.Pool

func init() {
	RedisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "127.0.0.1:6379") },
	}
}
