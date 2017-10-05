package db

import "github.com/garyburd/redigo/redis"

var conn = connect()

func connect() redis.Conn {
	connection, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	return connection
}
