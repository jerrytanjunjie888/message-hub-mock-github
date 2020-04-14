package config

import (
	"fmt"

	"github.com/go-redis/redis"
)

// LoadRedis states starts up the redis server
func LoadRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return client
}

// // SetID receives sequenceId and stores into Redis
// func SetID(client *redis.Client) error {
// 	err := clien
// 	return
// }

// // GetID retrieves last inserted sequenceId from Redis
// func GetID(k string) (v string, err error) {
// 	return
// }
