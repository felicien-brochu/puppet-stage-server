package db

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	addr = ":6379"
)

var (
	pool = newPool()
)

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

// func connect() redis.Conn {
// 	connection, err := redis.Dial("tcp", ":6379")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return connection
// }
