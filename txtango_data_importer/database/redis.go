package database

import (
	"os"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

//InitRedis initialize redis
func InitRedis() error {
	redisAddr := os.Getenv("REDIS_ADDR")

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default db
	})

	_, err := client.Ping().Result()
	if err != nil {
		return err
	}

	rdb = client
	return nil
}

//RDB returns a handle to the Redis object
func RDB() *redis.Client {
	return rdb
}
