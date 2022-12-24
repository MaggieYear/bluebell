package redis

import (
	"fmt"

	"bluebell/settings"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

// 初始化连接
func Init() (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			settings.RedisSettings.Host,
			settings.RedisSettings.Port),
		Password: settings.RedisSettings.Password,
		DB:       settings.RedisSettings.DbName,
		PoolSize: settings.RedisSettings.PoolSize, // 连接池大小
	})
	_, err = client.Ping().Result()
	return err
}

func Close() {
	_ = client.Close()
}
