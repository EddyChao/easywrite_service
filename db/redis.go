package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var (
	// Redis 声明一个全局的rdb变量
	Redis   *redis.Client
	Context = context.Background()
)

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

func (c *RedisConfig) copyToOption(option *redis.Options) {
	option.Addr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	option.Password = c.Password
	option.DB = c.Database
}

// InitRedis 初始化连接
func InitRedis(redisConfig RedisConfig) {

	option := &redis.Options{}
	redisConfig.copyToOption(option)
	Redis = redis.NewClient(option)

	_, err := Redis.Ping(Context).Result()
	if err != nil {
		panic(err)
	}
}

func RedisIsExist(key string) bool {
	if Redis.Exists(Context, key).Val() > 0 {
		return true
	}
	return false
}
